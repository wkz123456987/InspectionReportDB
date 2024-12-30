package detection

import (
	"GoBasic/utils/fileutils"
	"bytes"
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
	cfg, err := ini.Load("../config/database_config.ini")
	if cfg == nil || err != nil {
		logWriter.WriteLog("无法读取配置文件: " + err.Error())
	}
	section := cfg.Section("Linux")
	user := section.Key("User").String()
	// 加载配置文件 "database_config.ini"
	cfg_password, err := ini.LoadSources(ini.LoadOptions{
		AllowBooleanKeys:    true,
		IgnoreInlineComment: true, // 禁止#注释
	}, "../config/database_config.ini")

	if cfg_password == nil || err != nil {
		logWriter.WriteLog("无法读取配置文件: " + err.Error())
	}
	// 获取配置文件中 "Linux" 节
	section_password := cfg_password.Section("Linux")
	// 使用StringWithShadows来获取包括注释在内的整行内容
	password := section_password.Key("Password").String()

	port, err := section.Key("Port").Int()
	if err != nil {
		logWriter.WriteLog("无法转换端口号: " + err.Error())
	}
	host := section.Key("Host").String()

	sshConf := SSHConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
	return sshConf
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
