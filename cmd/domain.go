package cmd

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"

	"github.com/brutalgg/cli"
	"github.com/brutalgg/whoisenum/internal/rdap"
	"github.com/brutalgg/whoisenum/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/net/publicsuffix"
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
			cli.Errorln("Whois lookup error: ", e)
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
ScannerLoop:
	for scanner.Scan() {
		if i := scanner.Text(); i != "" {
			cli.Info("Searching Whois Records for Domain %v", i)
			// Get root domain
			rootDomain, err := getRootDomain(i)
			if err != nil {
				cli.Errorln("Root Domain Error: ", err)
			}
			// Check if root romain has been queried
			if result != nil {

				for index, entry := range result {
					if entry.Name == strings.ToUpper(rootDomain) {
						//Root domain has already been queried, add to root domain entry "DomainsSearched" .
						result[index].DomainsSearched = append(result[index].DomainsSearched, strings.ToUpper(i))
						//Skip query and exit out of current iteration of ScannerLoop.
						continue ScannerLoop
					}
				}
				if r, e := queryDomain(rootDomain); e != nil {
					cli.Errorln("Whois lookup error: ", e)
				} else {
					result = append(result, r)
				}
				time.Sleep(rd)
			} else {
				if r, e := queryDomain(rootDomain); e != nil {
					cli.Errorln("Whois lookup error: ", e)
				} else {
					result = append(result, r)
					if i != rootDomain {
						result[0].DomainsSearched = append(result[0].DomainsSearched, strings.ToUpper(i))
					}
				}
			}
		}
	}
	return result
}

func getRootDomain(domain string) (string, error) {
	// Extract root domain from provided domain.
	rootDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return rootDomain, err
	}
	// Check to see if root domain is different than domain provided.
	if rootDomain != domain {
		cli.Debug("Using Root Domain '%v' instead of '%v'", rootDomain, domain)
	}
	return rootDomain, nil
}

func queryDomain(s string) (rdap.WhoisDomainRecord, error) {
	r, e := rdap.GetWhoisDomainResults(s)
	if e != nil {
		return rdap.WhoisDomainRecord{}, e
	}
	return r, nil
}
