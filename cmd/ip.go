package cmd

import (
	"bufio"
	"log"
	"os"

	"github.com/brutalgg/whoisenum/internal/ipMath"
	"github.com/brutalgg/whoisenum/internal/output"
	"github.com/brutalgg/whoisenum/internal/rdap"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Search internet registries by ip address",
	Run:   basecmd,
}

func init() {
	rootCmd.AddCommand(ipCmd)
}

func basecmd(cmd *cobra.Command, args []string) {
	var result []rdap.WhoisIPEntry

	// TODO Parse some CLI args to get this info
	inFile := "sample.txt"

	f, _ := os.Open(inFile)
	defer f.Close()

	reader := bufio.NewScanner(f)
	for reader.Scan() {
		if i := reader.Text(); i != "" {
			if !uniqueNetworkCheck(i, result) {
				continue
			}
			qr, err := rdap.GetWhoisIPResults(i)
			//qr, err := getWhoisIPResults(i)
			if err != nil {
				log.Fatalln("Error with Whois lookup:", err)
			}
			result = append(result, qr)
		}
	}
	output.JsonWriteOut(os.Stdout, result)
}

func uniqueNetworkCheck(i string, r []rdap.WhoisIPEntry) bool {
	for indx := range r {
		if ipMath.NetworksContain(i, r[indx].CIDR...) {
			r[indx].IPSearched = append(r[indx].IPSearched, i)
			return false
		}
	}
	return true
}
