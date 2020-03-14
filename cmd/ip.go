package cmd

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/brutalgg/whoisenum/internal/ipMath"
	"github.com/brutalgg/whoisenum/internal/rdap"
	"github.com/brutalgg/whoisenum/internal/utils"
	"github.com/brutalgg/whoisenum/pkg/cli"
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
	inFile, _ := ctx.Flags().GetString("file")
	j, _ := ctx.Flags().GetBool("json")

	switch {
	case l != "" && inFile == "":
		cli.Info("Searching Whois Records for IP %v", l)
		cli.Infoln("This make take some time depending on the number of queries and your internet connection")
		if r, e := queryIP(l); e != nil {
			cli.Errorln("Whois lookup error", e)
		} else {
			result = append(result, r)
		}
	case inFile != "":
		f, _ := os.Open(inFile)
		defer f.Close()
		cli.Info("Searching Whois Records for IPs identified in %v", inFile)
		cli.Infoln("This make take some time depending on the number of queries and your internet connection")
		result = ipFileLogic(ctx, f)
	}

	if j {
		utils.JsonResultsOut(result)
	} else {
		for _, v := range result {
			v.PrintResult()
		}
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

func ipFileLogic(ctx *cobra.Command, readr io.ReadWriteSeeker) []rdap.WhoisIPRecord {
	var result []rdap.WhoisIPRecord
	r, _ := ctx.Flags().GetString("rate")
	rd, _ := time.ParseDuration(r)
	utils.SizeCheck(readr)
	scanner := bufio.NewScanner(readr)
	for scanner.Scan() {
		if i := scanner.Text(); i != "" {
			if !uniqueNetworkCheck(i, result) {
				continue
			}
			cli.Info("Searching Whois Records for IP %v", i)
			if r, e := queryIP(i); e != nil {
				cli.Errorln("Whois lookup error ", e)
			} else {
				result = append(result, r)
			}
			time.Sleep(rd)
		}
	}
	return result
}

func queryIP(s string) (rdap.WhoisIPRecord, error) {
	r, e := rdap.GetWhoisIPResults(s)
	if e != nil {
		return rdap.WhoisIPRecord{}, e
	}
	return r, nil
}
