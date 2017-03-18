# winrm-dns-client

#### Introduction

winrm-dns-client is a CLI and Go library for interacting with remote Microsoft DNS servers, it currently
utilises WinRM for remote connectivity, in the future when available, I will change this to use OpenSSH.
 
#### Requirements
In order to use this, the following must be met:
- DNS Server needs to be running on `Windows Server 2012`, or greater
- PowerShell must have the `DnsServer` module installed
 
#### Limitations
At present this only works with `A` and `CNAME` records, this is because of my own requirements for
this tool, I will look at adding more functionality as time progresses.

#### Usage
Configuration settings are found in $HOME/.winrm-dns-client.yaml
```
servername: <name of dns server>
username: <username to login>
password: <password to login>
```
To read all records within a zone
```
winrm-dns-client read -d <domain-name>
```

To read specific dns record within a zone
```
winrm-dns-client read -d >domain-name> -n <name>
```
To create a dns record within a zone
```
winrm-dns-client create -d <domain-name> -n <name> -t <record-type> -v <value> [-l <ttl>]
```