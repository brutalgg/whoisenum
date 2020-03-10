package rdap

import (
	"github.com/brutalgg/whoisenum/internal/cli"
	"github.com/brutalgg/whoisenum/internal/ipMath"

	"github.com/openrdap/rdap"
)

type WhoisIPRecord struct {
	NetworkAddress   string   `json:"network"`
	BroadcastAddress string   `json:"broadcast"`
	CIDR             []string `json:"cidr"`
	IPVersion        string   `json:"ip_version"`
	Type             string   `json:"type"`
	Parent           string   `json:"parent_network"`
	Organization     string   `json:"organization"`
	Registrar        string   `json:"registrar"`
	IPSearched       []string `json:"ip_searched"`
}

type WhoisDomainRecord struct {
	Handle      string   `json:"handle"`
	Name        string   `json:"name"`
	NameServers []string `json:"nameservers"`
	Status      []string `json:"status"`
	Reg         string   `json:"registration"`
	Exp         string   `json:"expiration"`
}

var tenDot WhoisIPRecord = WhoisIPRecord{
	NetworkAddress:   "10.0.0.0",
	BroadcastAddress: "10.255.255.255",
	CIDR:             []string{"10.0.0.0/8"},
	IPVersion:        "4",
	Type:             "PRIVATE-ADDRESS-ABLK-RFC1918-IANA-RESERVED",
	Parent:           "No Parent",
	Organization:     "PRIVATE-ADDRESS",
	Registrar:        "local",
}

var client = &rdap.Client{
	Verbose: func(text string) {
		cli.Debugln(text)
	},
}

func GetWhoisIPResults(ip string) (WhoisIPRecord, error) {
	whois := new(WhoisIPRecord)
	if ipMath.NetworksContain(ip, tenDot.CIDR...) {
		cli.Debug("10.0.0.0 - Special Case encounterd ")
		return tenDot, nil
	}
	queryResults, err := client.QueryIP(ip)
	if err != nil {
		return WhoisIPRecord{}, err
	}
	whois.NetworkAddress = queryResults.StartAddress
	whois.BroadcastAddress = queryResults.EndAddress
	whois.IPVersion = queryResults.IPVersion
	whois.Type = queryResults.Type
	whois.Parent = queryResults.ParentHandle
	whois.Organization = queryResults.Name
	whois.Registrar = queryResults.Port43
	whois.IPSearched = append(whois.IPSearched, ip)
	whois.CIDR = ipMath.Range2CIDRs(whois.NetworkAddress, whois.BroadcastAddress)

	return *whois, nil
}

func GetWhoisDomainResults(domain string) (WhoisDomainRecord, error) {
	whois := new(WhoisDomainRecord)
	queryResults, err := client.QueryDomain(domain)
	if err != nil {
		return WhoisDomainRecord{}, err
	}
	whois.Handle = queryResults.Handle
	whois.Name = queryResults.LDHName
	for _, v := range queryResults.Nameservers {
		whois.NameServers = append(whois.NameServers, v.LDHName)
	}
	whois.Status = queryResults.Status
	for _, v := range queryResults.Events {
		switch {
		case v.Action == "registration":
			whois.Reg = v.Date
		case v.Action == "expiration":
			whois.Exp = v.Date
		}
	}

	return *whois, nil
}
