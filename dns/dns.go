package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

type Record struct {
	Dnszone string
	Name    string
	Type    string
	Value   string
	TTL     int
}

type DnsRecordResponse struct {
	DistinguishedName string `json:"DistinguishedName"`
	HostName          string `json:"HostName"`
	RecordData        struct {
		CimClass struct {
			CimClassMethods     string `json:"CimClassMethods"`
			CimClassProperties  string `json:"CimClassProperties"`
			CimClassQualifiers  string `json:"CimClassQualifiers"`
			CimSuperClass       string `json:"CimSuperClass"`
			CimSuperClassName   string `json:"CimSuperClassName"`
			CimSystemProperties string `json:"CimSystemProperties"`
		} `json:"CimClass"`
		CimInstanceProperties []string `json:"CimInstanceProperties"`
		CimSystemProperties   struct {
			ClassName  string      `json:"ClassName"`
			Namespace  string      `json:"Namespace"`
			Path       interface{} `json:"Path"`
			ServerName string      `json:"ServerName"`
		} `json:"CimSystemProperties"`
	} `json:"RecordData"`
	RecordType string `json:"RecordType"`
	TimeToLive struct {
		Days              int64   `json:"Days"`
		Hours             int64   `json:"Hours"`
		Milliseconds      int64   `json:"Milliseconds"`
		Minutes           int64   `json:"Minutes"`
		Seconds           int64   `json:"Seconds"`
		Ticks             int64   `json:"Ticks"`
		TotalDays         float64 `json:"TotalDays"`
		TotalHours        int64   `json:"TotalHours"`
		TotalMilliseconds int64   `json:"TotalMilliseconds"`
		TotalMinutes      int64   `json:"TotalMinutes"`
		TotalSeconds      int64   `json:"TotalSeconds"`
	} `json:"TimeToLive"`
}

func tmplExec(r Record, tp string) (string, error) {
	t := template.New("tmpl")
	t, err := t.Parse(tp)
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %v\n", err)
	}
	var result bytes.Buffer
	if err := t.Execute(result, r); err != nil {
		return "", fmt.Errorf("Error generating template: %v\n", err)
	}

	return result.String(), nil
}

func unmarshalResponse(resp string, data interface{}) error {
	if err := json.Unmarshal([]byte(resp), &data); err != nil {
		return fmt.Errorf("Error unmarshalling json: %v", err)
	}

	return nil
}
