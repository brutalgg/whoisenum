package cmd

import (
	"bufio"
	"os"
	"time"

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
	inFile, _ := ctx.Flags().GetString("file")
	j, _ := ctx.Flags().GetBool("json")

	switch {
	case l == "" && inFile == "":
		cli.Debugln("Missing -l and -f. Trying Stdin...")
		info, err := os.Stdin.Stat()
		if err != nil {
			cli.Fatalln("No input found in -l, -f, or Stdin. Exiting...")
		}
		if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
			cli.Fatalln("No input found in -l, -f, or Stdin. Exiting...")
		}
		result = domainScannerLogic(ctx, bufio.NewScanner(os.Stdin))
	case l != "" && inFile == "":
		cli.Info("Searching Whois Records for Domain %v", l)
		cli.Infoln("This make take some time depending on the number of queries and your internet connection")
		if r, e := queryDomain(l); e != nil {
			cli.Errorln("Whois lookup error", e)
		} else {
			result = append(result, r)
		}
	case inFile != "":
		f, _ := os.Open(inFile)
		defer f.Close()
		cli.Info("Searching Whois Records for Domains identified in %v", inFile)
		cli.Infoln("This make take some time depending on the number of queries and your internet connection")
		result = domainScannerLogic(ctx, bufio.NewScanner(f))
	}

	if j {
		jsonResultsOut(result)
	} else {
		domainResultsOut(result)
	}
}

func domainScannerLogic(ctx *cobra.Command, scanner *bufio.Scanner)[]rdap.WhoisDomainRecord{
	r, _ := ctx.Flags().GetInt("rate")
	var result []rdap.WhoisDomainRecord
	for scanner.Scan() {
		if i := scanner.Text(); i != "" {
			cli.Info("Searching Whois Records for Domain %v", i)
			if r, e := queryDomain(i); e != nil {
				cli.Errorln("Whois lookup error", e)
			} else {
				result = append(result, r)
			}
			time.Sleep(time.Duration(r) * time.Second)
		}
	}
	return []rdap.WhoisDomainRecord
}

func queryDomain(s string) (rdap.WhoisDomainRecord, error) {
	r, e := rdap.GetWhoisDomainResults(s)
	if e != nil {
		return rdap.WhoisDomainRecord{}, e
	}
	return r, nil
}
