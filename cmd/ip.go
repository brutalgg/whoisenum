package cmd

import (
	"bufio"
	"os"
	"time"

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
		result = ipScannerLogic(ctx, bufio.NewScanner(os.Stdin))
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
		result = ipScannerLogic(ctx, bufio.NewScanner(f))
	}

	if j {
		jsonResultsOut(result)
	} else {
		ipResultsOut(result)
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

func ipScannerLogic(ctx *cobra.Command, scanner *bufio.Scanner) []rdap.WhoisIPRecord {
	var result []rdap.WhoisIPRecord
	r, _ := ctx.Flags().GetInt("rate")
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
			time.Sleep(time.Duration(r) * time.Second)
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
