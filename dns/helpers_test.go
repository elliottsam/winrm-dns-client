package dns

import (
	"testing"
)

const test1 string = `[{"DistinguishedName":"DC=test123,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","HostName":"test123","RecordData":{"CimClass":{"CimSuperClassName":"DnsServerResourceRecordData","CimSuperClass":"ROOT/Microsoft/Windows/DNS:DnsServerResourceRecordData","CimClassProperties":"HostNameAlias","CimClassQualifiers":"dynamic = True provider = \"DnsServerPSProvider\" ClassVersion = \"1.0.0\" locale = 1033","CimClassMethods":"","CimSystemProperties":"Microsoft.Management.Infrastructure.CimSystemProperties"},"CimInstanceProperties":["HostNameAlias = \"blah.com.\""],"CimSystemProperties":{"Namespace":"root/Microsoft/Windows/DNS","ServerName":"W2K12R2-DC","ClassName":"DnsServerResourceRecordCName","Path":null}},"RecordType":"CNAME","TimeToLive":{"Ticks":36000000000,"Days":0,"Hours":1,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":0.041666666666666664,"TotalHours":1,"TotalMilliseconds":3600000,"TotalMinutes":60,"TotalSeconds":3600}}]`
const test2 string = `[{"DistinguishedName":"DC=test123,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","HostName":"test123","RecordData":{"CimClass":{"CimSuperClassName":"DnsServerResourceRecordData","CimSuperClass":"ROOT/Microsoft/Windows/DNS:DnsServerResourceRecordData","CimClassProperties":"HostNameAlias","CimClassQualifiers":"dynamic = True provider = \"DnsServerPSProvider\" ClassVersion = \"1.0.0\" locale = 1033","CimClassMethods":"","CimSystemProperties":"Microsoft.Management.Infrastructure.CimSystemProperties"},"CimInstanceProperties":["HostNameAlias = \"blah.com.\""],"CimSystemProperties":{"Namespace":"root/Microsoft/Windows/DNS","ServerName":"W2K12R2-DC","ClassName":"DnsServerResourceRecordCName","Path":null}},"RecordType":"CNAME","TimeToLive":{"Ticks":36000000000,"Days":0,"Hours":1,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":0.041666666666666664,"TotalHours":1,"TotalMilliseconds":3600000,"TotalMinutes":60,"TotalSeconds":3600}},{"DistinguishedName":"DC=test123,DC=test.local,cn=MicrosoftDNS,DC=DomainDnsZones,DC=test,DC=local","HostName":"test123","RecordData":{"CimClass":{"CimSuperClassName":"DnsServerResourceRecordData","CimSuperClass":"ROOT/Microsoft/Windows/DNS:DnsServerResourceRecordData","CimClassProperties":"HostNameAlias","CimClassQualifiers":"dynamic = True provider = \"DnsServerPSProvider\" ClassVersion = \"1.0.0\" locale = 1033","CimClassMethods":"","CimSystemProperties":"Microsoft.Management.Infrastructure.CimSystemProperties"},"CimInstanceProperties":["HostNameAlias = \"blah.com.\""],"CimSystemProperties":{"Namespace":"root/Microsoft/Windows/DNS","ServerName":"W2K12R2-DC","ClassName":"DnsServerResourceRecordCName","Path":null}},"RecordType":"CNAME","TimeToLive":{"Ticks":36000000000,"Days":0,"Hours":1,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":0.041666666666666664,"TotalHours":1,"TotalMilliseconds":3600000,"TotalMinutes":60,"TotalSeconds":3600}}]`

func TestMakeResponseArray(t *testing.T) {
	const test1 string = `{
	"test": "string",
	"testint": 5
	}`

	const test2 string = `[{
	"test": "string",
	"testint": 5
	}]`

	result1 := makeResponseArray(test1)
	if result1 != test2 {
		t.Error("Failed to transform json into array")
	}
	result2 := makeResponseArray(test2)
	if result2 != test2 {
		t.Error("Failed to keep array as array")
	}
}

func TestUnmarshalResponse(t *testing.T) {
	resp, err := unmarshalResponse(test1)
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 1 {
		t.Error("Length should be 1")
	}

	resp, err = unmarshalResponse(test2)
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 2 {
		t.Error("Length should be 2")
	}
}

func TestConvertResponse(t *testing.T) {
	resp, _ := unmarshalResponse(test1)
	rec := convertResponse(resp, Record{Dnszone: "test.local"})
	if len(*rec) != 1 {
		t.Errorf("Expecting 1 record to be returned, got %v: %v", len(*rec), rec)
	}

	resp, _ = unmarshalResponse(test2)
	rec = convertResponse(resp, Record{Dnszone: "test.local"})
	if len(*rec) != 2 {
		t.Errorf("Expecting 2 record to be returned, got %v: %v", len(*rec), rec)
	}
}
