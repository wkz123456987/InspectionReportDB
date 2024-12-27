package detection

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// SSHConfig 用于配置SSH连接相关参数
type SSHConfig struct {
	User     string
	Password string
	Host     string
	Port     int
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
