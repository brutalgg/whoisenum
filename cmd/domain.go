package cmd

import (
	"github.com/brutalgg/whoisenum/internal/cli"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Search internet registries for a given domain",
	Args:  cobra.NoArgs,
	Run:   func(ctx *cobra.Command, args []string) { cli.Infoln("This is the domain command") },
}
