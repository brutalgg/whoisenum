package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/openrdap/rdap"
)

type whoisIPEntry struct {
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

func main() {
	var result []*whoisIPEntry

	// TODO Parse some CLI args to get this info
	inFile := "sample.txt"
	outFile := "result.json"

	f, _ := os.Open(inFile)
	defer f.Close()

	reader := bufio.NewScanner(f)
	for reader.Scan() {
		if i := reader.Text(); i != "" {
			if !uniqueNetworkCheck(i, result) {
				continue
			}
			qr, err := getWhoisIPResults(i)
			if err != nil {
				log.Fatalln("Error with Whois lookup:", err)
			}
			result = append(result, qr)
		}
	}
	jsonWriteOut(outFile, result)
}

func jsonWriteOut(f string, x interface{}) {
	outFile, _ := os.Create(f)
	defer outFile.Close()

	o, _ := json.MarshalIndent(x, "", "  ")
	outFile.Write(o)
}

func uniqueNetworkCheck(i string, r []*whoisIPEntry) bool {
	for _, w := range r {
		if NetworksContain(i, w.CIDR...) {
			w.IPSearched = append(w.IPSearched, i)
			return false
		}
	}
	return true
}

func getWhoisIPResults(ip string) (*whoisIPEntry, error) {
	whois := new(whoisIPEntry)
	client := &rdap.Client{}
	queryResults, err := client.QueryIP(ip)
	if err != nil {
		return nil, err
	}
	whois.NetworkAddress = queryResults.StartAddress
	whois.BroadcastAddress = queryResults.EndAddress
	whois.IPVersion = queryResults.IPVersion
	whois.Type = queryResults.Type
	whois.Parent = queryResults.ParentHandle
	whois.Organization = queryResults.Name
	whois.Registrar = queryResults.Port43
	whois.IPSearched = append(whois.IPSearched, ip)
	whois.CIDR = range2CIDRs(whois.NetworkAddress, whois.BroadcastAddress)

	return whois, nil
}

// https://groups.google.com/forum/#!topic/golang-nuts/rJvVwk4jwjQ
func range2CIDRs(ip1, ip2 string) (r []string) {
	allFF := net.ParseIP("255.255.255.255").To4()
	maxLen := 32
	a1 := net.ParseIP(ip1).To4()
	a2 := net.ParseIP(ip2).To4()
	for cmp(a1, a2) <= 0 {
		l := 32
		for l > 0 {
			m := net.CIDRMask(l-1, maxLen)
			if cmp(a1, first(a1, m)) != 0 || cmp(last(a1, m), a2) > 0 {
				break
			}
			l--
		}
		r = append(r, fmt.Sprintf("%v/%v", a1, l))
		//r = append(r, &net.IPNet{IP: a1, Mask: net.CIDRMask(l, maxLen)})
		a1 = last(a1, net.CIDRMask(l, maxLen))
		if cmp(a1, allFF) == 0 {
			break
		}
		a1 = next(a1)
	}
	return r
}

func next(ip net.IP) net.IP {
	n := len(ip)
	out := make(net.IP, n)
	copy := false
	for n > 0 {
		n--
		if copy {
			out[n] = ip[n]
			continue
		}
		if ip[n] < 255 {
			out[n] = ip[n] + 1
			copy = true
			continue
		}
		out[n] = 0
	}
	return out
}

func cmp(ip1, ip2 net.IP) int {
	l := len(ip1)
	for i := 0; i < l; i++ {
		if ip1[i] == ip2[i] {
			continue
		}
		if ip1[i] < ip2[i] {
			return -1
		}
		return 1
	}
	return 0
}

func first(ip net.IP, mask net.IPMask) net.IP {
	return ip.Mask(mask)
}

func last(ip net.IP, mask net.IPMask) net.IP {
	n := len(ip)
	out := make(net.IP, n)
	for i := 0; i < n; i++ {
		out[i] = ip[i] | ^mask[i]
	}
	return out
}

func NetworksContain(i string, n ...string) bool {
	for _, ns := range n {
		_, c, _ := net.ParseCIDR(ns)
		if c.Contains(net.ParseIP(i)) {
			return true
		}
	}
	return false
}
