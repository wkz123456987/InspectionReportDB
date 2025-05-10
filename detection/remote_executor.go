package detection

import (
	"GoBasic/config"
	"GoBasic/utils/fileutils"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/ssh"
	"gopkg.in/ini.v1"
)

// SSHConfig 用于配置SSH连接相关参数
type SSHConfig struct {
	User     string
	Password string
	Host     string
	Port     int
}

func GetSSHConfig(logWriter *fileutils.LogWriter) SSHConfig {
	// 先尝试加载配置文件获取基本信息
	cfg, err := ini.Load(config.ConfigPath)
	if cfg == nil || err != nil {
		logWriter.WriteLog("无法读取配置文件: " + err.Error())
	}
	section := cfg.Section("Linux")
	user := section.Key("User").String()
	port, err := section.Key("Port").Int()
	if err != nil {
		logWriter.WriteLog("无法转换端口号: " + err.Error())
	}
	host := section.Key("Host").String()

	// 加载配置文件 "database_config.ini"
	cfg_password, err := ini.LoadSources(ini.LoadOptions{
		AllowBooleanKeys:    true,
		IgnoreInlineComment: true, // 禁止#注释
	}, config.ConfigPath)

	if cfg_password == nil || err != nil {
		logWriter.WriteLog("无法读取配置文件: " + err.Error())
	}
	// 获取配置文件中 "Linux" 节
	section_password := cfg_password.Section("Linux")
	Encryption_method := section.Key("Encryption_method").String()
	// 使用StringWithShadows来获取包括注释在内的整行内容
	password_source := section_password.Key("Password").String()

	if Encryption_method == "plaintext" {
		password := password_source
		sshConf := SSHConfig{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
		}
		return sshConf

	} else {
		// 解密密码
		decryptedPassword, err := decrypt(password_source)
		if err != nil {
			logWriter.WriteLog("解密密码时出错: " + err.Error())
			return SSHConfig{}
		}
		password := decryptedPassword
		sshConf := SSHConfig{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
		}
		return sshConf
	}

}

// ExecuteRemoteCommand 通过SSH远程执行命令并返回结果
func ExecuteRemoteCommand(sshConf SSHConfig, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: sshConf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshConf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 建立SSH连接
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshConf.Host, sshConf.Port), config)
	if err != nil {
		return "", fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	defer client.Close()

	// 创建一个会话用于执行命令
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// 执行命令并获取输出
	var result bytes.Buffer
	session.Stdout = &result
	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("failed to execute command on remote server: %w", err)
	}

	return result.String(), nil
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
