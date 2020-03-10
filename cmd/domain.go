package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/brutalgg/whoisenum/internal/cli"
	"github.com/brutalgg/whoisenum/internal/rdap"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Search internet registries for a given domain",
	Args:  cobra.NoArgs,
	Run:   baseDomainCmd,
}

func baseDomainCmd(ctx *cobra.Command, args []string) {
	var result []rdap.WhoisDomainRecord
	l, _ := ctx.Flags().GetString("lookup")
	j, _ := ctx.Flags().GetBool("json")
	inFile, _ := ctx.Flags().GetString("file")

	if inFile != "" {
		f, _ := os.Open(inFile)
		defer f.Close()
		reader := bufio.NewScanner(f)
		for reader.Scan() {
			if i := reader.Text(); i != "" {
				qr, err := rdap.GetWhoisDomainResults(i)
				if err != nil {
					cli.Fatalln("Whois lookup error:", err)
				}
				result = append(result, qr)
			}
		}
	} else if l != "" {
		qr, err := rdap.GetWhoisDomainResults(l)
		if err != nil {
			cli.Fatalln("Whois lookup error:", err)
		}
		result = append(result, qr)
	} else {
		cli.Fatalln("Lookup and File flag not detected. The domain command requires at least one of these flags.")
	}

	if j {
		jsonResult(result)
	} else {
		standardDResult(result)
	}
}

func standardDResult(x []rdap.WhoisDomainRecord) {
	for _, v := range x {
		cli.WriteResults("Handle:", v.Handle)
		cli.WriteResults("Name:", v.Name)
		cli.WriteResults("Name Servers:\t", strings.Join(v.NameServers, "\n\t\t "))
		cli.WriteResults("Status:\t", strings.Join(v.Status, "\n\t "))
		cli.WriteResults("Registration:", v.Reg)
		cli.WriteResults("Expiration:", v.Exp)
		cli.WriteResults("--------------------")
	}
}
