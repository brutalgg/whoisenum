package cmd

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/brutalgg/whoisenum/internal/cli"
	"github.com/brutalgg/whoisenum/internal/ipMath"
	"github.com/brutalgg/whoisenum/internal/rdap"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Search internet registries by ip address",
	Run:   baseIPCmd,
}

func baseIPCmd(ctx *cobra.Command, args []string) {
	var result []rdap.WhoisIPRecord
	l, _ := ctx.Flags().GetString("lookup")
	j, _ := ctx.Flags().GetBool("json")
	inFile, _ := ctx.Flags().GetString("file")

	if inFile != "" {
		f, _ := os.Open(inFile)
		defer f.Close()
		reader := bufio.NewScanner(f)
		for reader.Scan() {
			if i := reader.Text(); i != "" {
				if !uniqueNetworkCheck(i, result) {
					continue
				}
				qr, err := rdap.GetWhoisIPResults(i)
				if err != nil {
					cli.Fatalln("Whois lookup error:", err)
				}
				result = append(result, qr)
			}
		}
	} else if l != "" {
		qr, err := rdap.GetWhoisIPResults(l)
		if err != nil {
			cli.Fatalln("Whois lookup error:", err)
		}
		result = append(result, qr)
	} else {
		cli.Fatalln("Lookup and File flag not detected. The IP command requires at least one of these flags.")
	}

	if j {
		jsonResult(result)
	} else {
		standardResult(result)
	}
}

func uniqueNetworkCheck(i string, r []rdap.WhoisIPRecord) bool {
	for indx := range r {
		if ipMath.NetworksContain(i, r[indx].CIDR...) {
			r[indx].IPSearched = append(r[indx].IPSearched, i)
			return false
		}
	}
	return true
}

func jsonResult(x interface{}) error {
	o, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}
	cli.WriteResults(string(o))
	return nil
}

func standardResult(x []rdap.WhoisIPRecord) {
	for _, v := range x {
		cli.WriteResults("Registrar:", v.Registrar)
		cli.WriteResults("Starting IP:", v.NetworkAddress)
		cli.WriteResults("Ending IP:", v.BroadcastAddress)
		cli.WriteResults("CIDR:\t", strings.Join(v.CIDR, "\n\t "))
		cli.WriteResults("IP Version:", v.IPVersion)
		cli.WriteResults("Registration Type:", v.Type)
		cli.WriteResults("Parent Registration:", v.Parent)
		cli.WriteResults("Organization:", v.Organization)
		cli.WriteResults("IPs Searched:\t", strings.Join(v.IPSearched, "\n\t\t "))
		cli.WriteResults("--------------------")
	}
}
