package cmd

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/brutalgg/whoisenum/internal/cli"
	"github.com/brutalgg/whoisenum/internal/ipMath"
	"github.com/brutalgg/whoisenum/internal/rdap"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Search internet registries by ip address",
	Run:   basecmd,
}

func basecmd(ctx *cobra.Command, args []string) {
	inFile, _ := ctx.Flags().GetString("file")
	var result []rdap.WhoisIPEntry

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
	jsonOut(result)
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

func jsonOut(x interface{}) error {
	o, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}
	cli.NoFormatString(string(o))
	return nil
}
