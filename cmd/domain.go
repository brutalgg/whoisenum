package cmd

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/brutalgg/whoisenum/internal/rdap"
	"github.com/brutalgg/whoisenum/internal/utils"
	"github.com/brutalgg/whoisenum/pkg/cli"
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
		result = domainScannerLogic(ctx, f)
	}

	if j {
		utils.JsonResultsOut(result)
	} else {
		for _, v := range result {
			v.PrintResult()
		}
	}
}

func domainScannerLogic(ctx *cobra.Command, readr io.ReadWriteSeeker) []rdap.WhoisDomainRecord {
	var result []rdap.WhoisDomainRecord
	r, _ := ctx.Flags().GetString("rate")
	rd, _ := time.ParseDuration(r)
	utils.SizeCheck(readr)
	scanner := bufio.NewScanner(readr)
	for scanner.Scan() {
		if i := scanner.Text(); i != "" {
			cli.Info("Searching Whois Records for Domain %v", i)
			if r, e := queryDomain(i); e != nil {
				cli.Errorln("Whois lookup error", e)
			} else {
				result = append(result, r)
			}
			time.Sleep(rd)
		}
	}
	return result
}

func queryDomain(s string) (rdap.WhoisDomainRecord, error) {
	r, e := rdap.GetWhoisDomainResults(s)
	if e != nil {
		return rdap.WhoisDomainRecord{}, e
	}
	return r, nil
}
