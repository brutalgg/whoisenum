package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/brutalgg/whoisenum/pkg/cli"
)

func SizeCheck(readr io.ReadSeeker) {
	if i, _ := lineCounter(readr); i > 500 {
		cli.Warn("The provided input is over 500 lines. This may result in bans from the queries APIs.")
		if !confirm() {
			cli.Info("Quiting...")
			os.Exit(0)
		}
	}
}

func confirm() bool {
	var s string

	fmt.Printf("Would you like to continue? (y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		cli.Error("%v", err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s[:1] == "y" {
		return true
	}
	return false
}

func JsonResultsOut(x interface{}) error {
	o, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}
	cli.WriteResults(string(o))
	return nil
}

func lineCounter(r io.ReadSeeker) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			r.Seek(0, io.SeekStart)
			return count, nil

		case err != nil:
			r.Seek(0, io.SeekStart)
			return count, err
		}
	}
}
