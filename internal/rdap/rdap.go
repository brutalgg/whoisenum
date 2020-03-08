package rdap

import (
	"github.com/brutalgg/whoisenum/internal/ipMath"

	"github.com/openrdap/rdap"
)

type WhoisIPEntry struct {
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

func GetWhoisIPResults(ip string) (WhoisIPEntry, error) {
	whois := new(WhoisIPEntry)
	client := &rdap.Client{}
	queryResults, err := client.QueryIP(ip)
	if err != nil {
		return WhoisIPEntry{}, err
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
