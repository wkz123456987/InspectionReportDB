package inspection

import (
	"GoBasic/config"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/ini.v1"
)

func ConnectPostgreSQL(queryIdentifier string, dbname ...string) [][]string {
	// 获取当前项目工程的绝对路径
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取项目工程路径失败: %v", err)
	}

	// 拼接SQL文件的完整路径并读取文件内容
	sqlBytes, err := ioutil.ReadFile(filepath.Join(projectPath, "../SQL", "inspection.sql"))
	if err != nil {
		log.Fatalf("读取SQL文件失败: %v", err)
	}

	// 查找目标SQL语句，传递查询标识符参数
	targetSQL := extractTargetSQL(string(sqlBytes), queryIdentifier)
	if targetSQL == "" {
		log.Fatal("未找到匹配的SQL语句")
	}
	// 读取配置文件获取数据库配置信息
	cfg, err := ini.Load(config.ConfigPath)
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}
	section := cfg.Section("Database")
	dbName := section.Key("DBName").String()
	hostname := section.Key("Hostname").String()
	port := section.Key("Port").String()
	username := section.Key("Username").String()
	// 加载配置文件 "database_config.ini"
	cfg_password, err := ini.LoadSources(ini.LoadOptions{
		AllowBooleanKeys:    true,
		IgnoreInlineComment: true, // 禁止#注释
	}, config.ConfigPath)

	if cfg_password == nil || err != nil {
		log.Fatalf("无法读取配置文件: " + err.Error())
	}
	// 获取配置文件中 "Linux" 节
	section_password := cfg_password.Section("Database")
	// 使用StringWithShadows来获取包括注释在内的整行内容
	//password := section_password.Key("Password").String()

	Encryption_method := section.Key("Encryption_method").String()
	// 使用StringWithShadows来获取包括注释在内的整行内容
	password_source := section_password.Key("Password").String()
	var password *string
	if Encryption_method == "plaintext" {
		password = &password_source

	} else {
		// 解密密码
		decryptedPassword, err := decrypt(password_source)
		if err != nil {
			fmt.Println("解密密码时出错: " + err.Error())
		}
		password = &decryptedPassword
	}

	// 构建连接字符串并连接数据库，同时检查连接有效性
	var actualDBName string
	if len(dbname) > 0 {
		actualDBName = dbname[0]
	} else {
		actualDBName = dbName
	}

	// 构建连接字符串并连接数据库，同时检查连接有效性
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		username, *password, actualDBName, hostname, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// 查询数据库并获取表格形式的数据结果
	rows, err := db.Query(targetSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 用于存储表格数据，每一个内层切片代表一行数据
	tableData := [][]string{}
	// 读取列信息
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	// 用于存储每行数据的扫描结果，根据列数创建对应长度的切片
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	// 循环读取每一行数据并添加到tableData中，同时处理不同数据类型转换为字符串，并处理空值显示为空字符串
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}
		rowData := make([]string, len(columns))
		for j := range values {
			switch v := values[j].(type) {
			case []byte:
				rowData[j] = string(v)
			case time.Time:
				rowData[j] = v.Format("2006-01-02 15:04:05")
			default:
				if v == nil {
					rowData[j] = ""
				} else {
					rowData[j] = fmt.Sprintf("%v", v)
				}
			}
		}
		tableData = append(tableData, rowData)
	}

	return tableData
}

// extractTargetSQL从SQL文件内容中提取以指定标识符开头的目标SQL语句，接收SQL内容和标识符作为参数
func extractTargetSQL(sqlContent string, queryIdentifier string) string {
	lines := strings.Split(strings.TrimSpace(sqlContent), "\n")
	var targetIndex int
	for index, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, queryIdentifier) {
			targetIndex = index + 1
			break
		}
	}
	if targetIndex < len(lines) {
		var targetSQL string
		for i := targetIndex; i < len(lines); i++ {
			currentLine := strings.TrimSpace(lines[i])
			// 查找 -- 在当前行中的位置，若存在则截取到 -- 之前的内容（去除注释部分）
			commentIndex := strings.Index(currentLine, "--")
			if commentIndex >= 0 {
				currentLine = currentLine[:commentIndex]
			}
			// 去除当前行两边的空白字符
			currentLine = strings.TrimSpace(currentLine)
			if currentLine == "" {
				break
			}
			// 如果targetSQL为空，直接赋值当前行内容；否则用空格连接当前行内容
			if targetSQL == "" {
				targetSQL = currentLine
			} else {
				targetSQL += " " + currentLine
			}
		}
		// 去除最终目标SQL语句中多余的连续空格，将多个连续空格替换为单个空格
		targetSQL = strings.Join(strings.Fields(targetSQL), " ")
		return targetSQL
	}
	return ""
}

// 解密函数
func decrypt(encryptedString string) (string, error) {

	// 解密密钥，长度必须是 16、24 或 32 字节，分别对应 AES-128、AES-192 或 AES-256
	const decryptionKey = "1234567890abcdef1234567890abcdef"
	// 解码 Base64 字符串
	decodedData, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", fmt.Errorf("解码 Base64 时出错: %w", err)
	}

	// 提取 IV 和密文
	blockSize := aes.BlockSize
	iv := decodedData[:blockSize]
	ciphertext := decodedData[blockSize:]

	// 创建 AES 解密器
	block, err := aes.NewCipher([]byte(decryptionKey))
	if err != nil {
		return "", fmt.Errorf("创建 AES 解密器时出错: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去除 PKCS7 填充
	plaintext := pkcs7UnPadding(ciphertext)
	return string(plaintext), nil
}

// PKCS7 去填充
func pkcs7UnPadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
