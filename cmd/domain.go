package cmd

import (
	"bufio"
	"os"

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

	switch {
	case l == "" && inFile == "":
		cli.Fatalln("Lookup and File flag not detected. The domain command requires at least one of these flags.")
	case l != "" && inFile == "":
		cli.Info("Searching Whois Records for IP %v", l)
		cli.Infoln("This make take some time depending on the number of queries and your internet connection")
		if r, e := queryDomain(l); e != nil {
			cli.Errorln("Whois lookup error", e)
		} else {
			result = append(result, r)
		}
	case inFile != "":
		f, _ := os.Open(inFile)
		defer f.Close()
		cli.Info("Searching Whois Records for IPs identified in %v", inFile)
		cli.Infoln("This make take some time depending on the number of queries and your internet connection")
		reader := bufio.NewScanner(f)
		for reader.Scan() {
			if i := reader.Text(); i != "" {
				if r, e := queryDomain(i); e != nil {
					cli.Errorln("Whois lookup error", e)
				} else {
					result = append(result, r)
				}
			}
		}
	}

	if j {
		jsonResultsOut(result)
	} else {
		domainResultsOut(result)
	}
}

func queryDomain(l string) (rdap.WhoisDomainRecord, error) {
	r, e := rdap.GetWhoisDomainResults(l)
	if e != nil {
		return rdap.WhoisDomainRecord{}, e
	}
	return r, nil
}
