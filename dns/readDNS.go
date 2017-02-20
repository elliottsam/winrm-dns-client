package dns

import (
	"fmt"
)

func readRecord(c Config, r Record, dnszone string, record ...string) (*DnsRecordResponse, error) {
	// Powershell script template to read record from DNS
	const tmplpscript = `
Get-DnsServerResourceRecord -ZoneName {{.Dnszone}}  -Name {{.Name}} | select DistinguishedName, HostName, RecordData, RecordType, TimeToLive | ConvertTo-Json
`
	// Setup
	if err := c.ConfigureWinRMClient(); err != nil {
		panic(err)
	}
	pscript, err := tmplExec(r, tmplpscript)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	command := Powershell(pscript)

	output, err := c.ExecutePowerShellScript(command)
	if err != nil {
		return nil, fmt.Errorf("Error running PowerShell script: %v\n", err)
	}

	resp := DnsRecordResponse{}
	if err := unmarshalResponse(output.stdout, resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
