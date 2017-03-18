package dns

import (
	"fmt"
	"strings"
)

// Record containing information regarding DNS record
type Record struct {
	Dnszone string
	Name    string
	Type    string
	Value   string
	TTL     float64
	ID      string
}

// ReadRecord performs DNS Record lookup from server
func ReadRecord(c *Client, rec Record) ([]Record, error) {
	// powershell script template to read record from DNS
	const tmplpscript = `
Get-DnsServerResourceRecord -ZoneName {{.Dnszone}}{{ if .Name }} -Name {{.Name}}{{end}} | ?{$_.RecordType -eq 'A' -or $_.RecordType -eq 'CNAME'} | select DistinguishedName, HostName, RecordData, RecordType, TimeToLive | ConvertTo-Json
`

	pscript, err := tmplExec(rec, tmplpscript)
	if err != nil {
		return []Record{}, fmt.Errorf("Error creating template: %v", err)
	}
	output, err := c.ExecutePowerShellScript(pscript)
	if err != nil {
		return []Record{}, fmt.Errorf("Error running PowerShell script: %v", err)
	}
	output.stdout = makeResponseArray(output.stdout)
	resp, err := unmarshalResponse(output.stdout)
	if err != nil {
		return []Record{}, fmt.Errorf("Error unmarshalling response: %v", err)
	}
	return *convertResponse(resp, rec), nil
}

// ReadRecordfromID retrieves specifc DNS record based on record ID
func ReadRecordfromID(c *Client, recID string) (Record, error) {
	id := strings.Split(recID, "|")
	if len(id) != 3 {
		return Record{}, fmt.Errorf("ID is incorrect")
	}
	rec := Record{
		Dnszone: id[0],
		Name:    id[1],
		Value:   id[2],
	}
	result, err := ReadRecord(c, rec)
	if err != nil {
		return Record{}, fmt.Errorf("Reading record: %v", err)
	}
	for i, v := range result {
		if v.ID == recID {
			return result[i], nil
		}
	}
	return Record{}, fmt.Errorf("Record not found: %v", recID)
}

// CreateRecord creates new DNS records on server
func CreateRecord(c *Client, rec Record) ([]Record, error) {
	const tmplscriptA = `
Add-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} -A -IPv4Address {{ .Value }} -TimeToLive (New-TimeSpan -Seconds {{ .TTL }})
`
	const tmplscriptCname = `
Add-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} -CName -HostNameAlias {{ .Value }} -TimeToLive (New-TimeSpan -Seconds {{ .TTL }})
`
	var (
		pscript string
		err     error
	)

	if RecordExist(c, rec) {
		return []Record{}, fmt.Errorf("Record already exists: %v", rec)
	}
	rec.ID = fmt.Sprintf("%s|%s|%s", rec.Dnszone, rec.Name, rec.Value)
	switch rec.Type {
	case "A":
		pscript, err = tmplExec(rec, tmplscriptA)
		if err != nil {
			return []Record{}, fmt.Errorf("Error creating template: %v", err)
		}
	case "CNAME":
		pscript, err = tmplExec(rec, tmplscriptCname)
		if err != nil {
			return []Record{}, fmt.Errorf("Error creating template: %v", err)
		}
	}
	fmt.Println(pscript)
	_, err = c.ExecutePowerShellScript(pscript)
	if err != nil {
		return []Record{}, fmt.Errorf("Error executing PowerShell script: %v", err)
	}
	record, err := ReadRecordfromID(c, rec.ID)
	if err != nil {
		return []Record{}, fmt.Errorf("Error reading record: %v", err)
	}

	var result []Record
	result = append(result, record)

	return result, nil
}

// RecordExist returns if record exists or not
func RecordExist(c *Client, rec Record) bool {
	records, _ := ReadRecord(c, rec)
	if len(records) > 0 {
		for _, v := range records {
			if v.Value == rec.Value {
				return true
			}
		}
	}
	return false
}
