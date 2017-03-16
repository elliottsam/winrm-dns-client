package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
)

type Record struct {
	Dnszone string
	Name    string
	Type    string
	Value   string
	TTL     float64
}

func ReadRecord(c Client, dnszone string, name string) *[]Record {
	// Powershell script template to read record from DNS
	const tmplpscript = `
Get-DnsServerResourceRecord -ZoneName {{.Dnszone}}{{ if .Name }} -Name {{.Name}}{{end}} | ?{$_.RecordType -eq 'A' -or $_.RecordType -eq 'CNAME'} | select DistinguishedName, HostName, RecordData, RecordType, TimeToLive | ConvertTo-Json
`
	r := Record{
		Dnszone: dnszone,
		Name:    name,
	}
	pscript, err := tmplExec(r, tmplpscript)
	if err != nil {
		log.Fatal(fmt.Errorf("%v", err))
	}
	command := Powershell(pscript)

	output, err := c.ExecutePowerShellScript(command)
	if err != nil {
		log.Fatal(fmt.Errorf("Error running PowerShell script: %v\n", err))
	}
	output.stdout = makeResponseArray(output.stdout)
	resp, err := unmarshalResponse(output.stdout)
	if err != nil {
		log.Fatal(err)
	}
	return convertResponse(resp, dnszone)
}

func tmplExec(r Record, tp string) (string, error) {
	t := template.New("tmpl")
	t, err := t.Parse(tp)
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %v\n", err)
	}
	var result bytes.Buffer
	if err := t.Execute(&result, r); err != nil {
		return "", fmt.Errorf("Error generating template: %v\n", err)
	}

	return result.String(), nil
}

func unmarshalResponse(resp string) ([]interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(resp), &data); err != nil {
		return nil, fmt.Errorf("Error unmarshalling json: %v", err)
	}
	return data.([]interface{}), nil
}

func convertResponse(r []interface{}, dnsZone string) *[]Record {
	records := []Record{}
	for i := range r {
		var rec Record
		switch r[i].(map[string]interface{})["RecordData"].(map[string]interface{})["CimInstanceProperties"].(type) {
		case []interface{}:
			rec = Record{
				Dnszone: dnsZone,
				Name:    r[i].(map[string]interface{})["HostName"].(string),
				Type:    r[i].(map[string]interface{})["RecordType"].(string),
				Value:   strings.Split(r[i].(map[string]interface{})["RecordData"].(map[string]interface{})["CimInstanceProperties"].([]interface{})[0].(string), "\"")[1],
				TTL:     r[i].(map[string]interface{})["TimeToLive"].(map[string]interface{})["TotalSeconds"].(float64),
			}
		case string:
			rec = Record{
				Dnszone: dnsZone,
				Name:    r[i].(map[string]interface{})["HostName"].(string),
				Type:    r[i].(map[string]interface{})["RecordType"].(string),
				Value:   strings.Split(r[i].(map[string]interface{})["RecordData"].(map[string]interface{})["CimInstanceProperties"].(string), "\"")[1],
				TTL:     r[i].(map[string]interface{})["TimeToLive"].(map[string]interface{})["TotalSeconds"].(float64),
			}
		}
		records = append(records, rec)

	}
	return &records
}

func makeResponseArray(r string) string {
	if rune(r[0]) != '[' && rune(r[(len(r)-1)]) != ']' {
		return fmt.Sprintf("[%s]", r)
	}
	return r
}

func OutputTable(rec *[]Record) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"DnsZone", "Name", "Type", "Value", "TTL"})
	for _, v := range *rec {
		table.Append([]string{v.Dnszone, v.Name, v.Type, v.Value, v.Type})
	}
	table.Render()
}
