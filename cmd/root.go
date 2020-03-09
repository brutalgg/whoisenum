package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "whoisenum",
	Short: "A tool to pull whois information from internet registries",
	Long: `A tool for pulling whois information from internet registries.
Useful for open source inteligence gathering during penetration 
testing or other security assessments.`,
	Version: "0.0.1a",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	//continue
}
