package cmd

import (
	"encoding/json"
	"strings"

	"github.com/brutalgg/whoisenum/internal/cli"
	"github.com/brutalgg/whoisenum/internal/rdap"
)

func domainResultsOut(x []rdap.WhoisDomainRecord) {
	for _, v := range x {
		cli.Infoln("Query Results")
		cli.WriteResults("Handle:", v.Handle)
		cli.WriteResults("Name:", v.Name)
		cli.WriteResults("Name Servers:\t", strings.Join(v.NameServers, "\n\t\t "))
		cli.WriteResults("Status:\t", strings.Join(v.Status, "\n\t "))
		cli.WriteResults("Registration:", v.Reg)
		cli.WriteResults("Expiration:", v.Exp)
		cli.WriteResults("")
	}
}

func ipResultsOut(x []rdap.WhoisIPRecord) {
	for _, v := range x {
		cli.Infoln("Query Results")
		cli.WriteResults("Registrar:", v.Registrar)
		cli.WriteResults("Starting IP:", v.NetworkAddress)
		cli.WriteResults("Ending IP:", v.BroadcastAddress)
		cli.WriteResults("CIDR:\t", strings.Join(v.CIDR, "\n\t "))
		cli.WriteResults("IP Version:", v.IPVersion)
		cli.WriteResults("Registration Type:", v.Type)
		cli.WriteResults("Parent Registration:", v.Parent)
		cli.WriteResults("Organization:", v.Organization)
		cli.WriteResults("IPs Searched:\t", strings.Join(v.IPSearched, "\n\t\t "))
		cli.WriteResults("")
	}
}

func jsonResultsOut(x interface{}) error {
	o, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}
	cli.WriteResults(string(o))
	return nil
}
