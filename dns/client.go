package dns

import (
	"encoding/base64"
	"fmt"

	"github.com/masterzen/winrm"
)

// Client struct for holding winrm.Client configuration
type Client struct {
	ServerName string
	Username   string
	Password   string
	Client     *winrm.Client
}

// Output returned from running PS scripts on WinRm server
type Output struct {
	stdout   string
	stderr   string
	exitcode int
}

// GenerateClient generates the winrm.client configuration
func GenerateClient(sn, un, pwd string) *Client {
	return &Client{
		ServerName: sn,
		Username:   un,
		Password:   pwd,
	}
}

// ConfigureWinRMClient creates the connection to the winrm server
func (c *Client) ConfigureWinRMClient() error {
	endpoint := winrm.NewEndpoint(c.ServerName, 5985, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, c.Username, c.Password)
	if err != nil {
		return fmt.Errorf("Error creating WinRM client: %v", err)
	}
	c.Client = client

	return nil
}

// ExecutePowerShellScript runs a PS script on the winrm server
func (c *Client) ExecutePowerShellScript(pscript string) (*Output, error) {
	command := winrm.Powershell(pscript)
	out, outerr, exitcode, err := c.Client.RunWithString(command, "")
	if err != nil {
		return nil, fmt.Errorf("Error executing script: %v", err)
	}

	return &Output{stdout: out, stderr: outerr, exitcode: exitcode}, nil
}

func powershell(psCmd string) string {
	wideCmd := ""
	for _, b := range []byte(psCmd) {
		wideCmd += string(b) + "\x00"
	}
	input := []uint8(wideCmd)
	encodedCmd := base64.StdEncoding.EncodeToString(input)
	return fmt.Sprintf("powershell.exe -EncodedCommand %s", encodedCmd)
}
