package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Search internet registries for a given domain",
	Args:  cobra.NoArgs,
	Run:   func(ctx *cobra.Command, args []string) { fmt.Println("This is the DOMAIN command") },
}
