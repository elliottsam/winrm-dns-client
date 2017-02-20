package dns

import (
	"encoding/base64"
	"fmt"

	"github.com/masterzen/winrm"
)

type Config struct {
	ServerName string
	Username   string
	Password   string
	Client     *winrm.Client
}

type Output struct {
	stdout   string
	stderr   string
	exitcode int
}

func GenerateConfig(sn, un, pwd string) *Config {
	return &Config{
		ServerName: sn,
		Username:   un,
		Password:   pwd,
	}
}

func (c *Config) ConfigureWinRMClient() error {
	endpoint := winrm.NewEndpoint(c.ServerName, 5985, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, c.Username, c.Password)
	if err != nil {
		return fmt.Errorf("Error creating WinRM client: %v\n\n", err)
	}
	c.Client = client

	return nil
}

func (c *Config) ExecutePowerShellScript(pscript string) (*Output, error) {
	command := winrm.Powershell(pscript)
	out, outerr, exitcode, err := c.Client.RunWithString(command, "")
	if err != nil {
		return nil, fmt.Errorf("Error executing script: %v\n\n", err)
	}

	return &Output{stdout: out, stderr: outerr, exitcode: exitcode}, nil
}

func Powershell(psCmd string) string {
	wideCmd := ""
	for _, b := range []byte(psCmd) {
		wideCmd += string(b) + "\x00"
	}
	input := []uint8(wideCmd)
	encodedCmd := base64.StdEncoding.EncodeToString(input)
	return fmt.Sprintf("powershell.exe -EncodedCommand %s", encodedCmd)
}
