package dns

import (
	"fmt"
	"log"
)

// Record containing information regarding DNS record
type Record struct {
	Dnszone string
	Name    string
	Type    string
	Value   string
	TTL     float64
	Id      string
}

// ReadRecord performs DNS Record lookup from server
func ReadRecord(c *Client, rec Record) []Record {
	// powershell script template to read record from DNS
	const tmplpscript = `
Get-DnsServerResourceRecord -ZoneName {{.Dnszone}}{{ if .Name }} -Name {{.Name}}{{end}} | ?{$_.RecordType -eq 'A' -or $_.RecordType -eq 'CNAME'} | select DistinguishedName, HostName, RecordData, RecordType, TimeToLive | ConvertTo-Json
`

	pscript, err := tmplExec(rec, tmplpscript)
	if err != nil {
		log.Fatal(fmt.Errorf("%v", err))
	}
	command := powershell(pscript)

	output, err := c.ExecutePowerShellScript(command)
	if err != nil {
		log.Fatalln(fmt.Errorf("Error running PowerShell script: %v", err))
	}
	output.stdout = makeResponseArray(output.stdout)
	resp, err := unmarshalResponse(output.stdout)
	if err != nil {
		log.Fatal(err)
	}
	return *convertResponse(resp, rec)
}

// CreateRecord creates new DNS records on server
func CreateRecord(c *Client, rec Record) (Record, error) {
	const tmplscriptA = `
Add-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} -A -IpAddress {{ .Value }} -TimeToLive (New-TimeSpan -Seconds {{ .TTL }}
`
	const tmplscriptCname = `
Add-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} -CName -HostNameAlias {{ .Value }} -TimeToLive (New-TimeSpan -Seconds {{ .TTL }})
`
	var (
		pscript string
		err     error
	)
	switch rec.Type {
	case "A":
		pscript, err = tmplExec(rec, tmplscriptA)
		if err != nil {
			return Record{}, fmt.Errorf("Error creating template: %v", err)
		}
	case "CNAME":
		pscript, err = tmplExec(rec, tmplscriptCname)
		if err != nil {
			return Record{}, fmt.Errorf("Error creating template: %v", err)
		}
	}
	command := powershell(pscript)
	_, err = c.ExecutePowerShellScript(command)
	if err != nil {
		return Record{}, fmt.Errorf("Error executing PowerShell script: %v", err)
	}
	return ReadRecord(c, rec)[0], nil
}
