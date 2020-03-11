package cmd

import (
	"github.com/brutalgg/whoisenum/internal/banner"
	"github.com/brutalgg/whoisenum/internal/cli"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "whoisenum",
	Version:          "0.1beta",
	PersistentPreRun: setup,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func setup(ctx *cobra.Command, args []string) {
	// Check the relevant flags
	q, _ := ctx.Flags().GetBool("quiet")
	v, _ := ctx.Flags().GetBool("verbose")
	l, _ := ctx.Flags().GetString("lookup")
	f, _ := ctx.Flags().GetString("file")
	r, _ := ctx.Flags().GetInt("rate")

	// Set our output level
	switch {
	case q:
		cli.SetPrintLevel(cli.LevelFatal)
	case v:
		cli.SetPrintLevel(cli.LevelDebug)
	}

	// Print the Banner
	cli.WriteBanner(banner.Banner)

	// Warn if both --lookup and --file are used
	if l != "" && f != "" {
		cli.Warnln("File flag detected. Ignoring lookup flag...")
	}

	// Notify if rate limiting is active
	if r != 0 {
		cli.Debug("Rate limiting detected. Sending queries every %v second(s)", r)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cli.Fatalln(err)
	}
	cli.Infoln("Queries Complete...")
	cli.Infoln("Goodbye.")
}

func init() {
	// Add additional commands to our CLI interface
	rootCmd.AddCommand(ipCmd)
	rootCmd.AddCommand(domainCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().Bool("json", false, "Output results in JSON")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet extra program output and only print results")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Include verbose messages from program execution")
	rootCmd.PersistentFlags().StringP("file", "f", "", "File with single IP/Domain per line")
	rootCmd.PersistentFlags().StringP("lookup", "l", "", "Single IP/Domain to lookup. Has no effect when --file is also specified.")
	rootCmd.PersistentFlags().IntP("rate", "r", 0, "The number of seconds between queries")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
