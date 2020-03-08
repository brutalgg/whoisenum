package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/brutalgg/whoisenum/internal/ipMath"
	"github.com/brutalgg/whoisenum/internal/rdap"
)

func uniqueNetworkCheck(i string, r []rdap.WhoisIPEntry) bool {
	for indx := range r {
		if ipMath.NetworksContain(i, r[indx].CIDR...) {
			r[indx].IPSearched = append(r[indx].IPSearched, i)
			return false
		}
	}
	return true
}

func main() {
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
	jsonWriteOut(os.Stdout, result)
}

func jsonWriteOut(out io.Writer, x interface{}) error {
	o, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}
	_, err = out.Write(o)
	return err
}
