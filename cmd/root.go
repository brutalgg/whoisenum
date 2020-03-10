package cmd

import (
	"github.com/brutalgg/whoisenum/internal/banner"
	"github.com/brutalgg/whoisenum/internal/cli"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "whoisenum",
	Version: "0.0.1a",
	PersistentPreRun: func(ctx *cobra.Command, args []string) {
		l, _ := ctx.Flags().GetString("lookup")
		f, _ := ctx.Flags().GetString("file")
		if l != "" && f != "" {
			cli.Warnln("lookup and file flags both provided. Ignoring lookup flag...")
			//return errors.New("options --lookup and --file are mutually exclusive. Please only provide one of these options to continue")
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cli.NoFormatString(banner.Banner)
	if err := rootCmd.Execute(); err != nil {
		cli.Errorln(err)
	}
}

func init() {
	// Add additional commands to our CLI interface
	rootCmd.AddCommand(ipCmd)
	rootCmd.AddCommand(domainCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON in place of the default format")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet extra program output and only print results")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Include verbose messages from program execution")
	rootCmd.PersistentFlags().StringP("file", "f", "", "File with single IP/Domain per line")
	rootCmd.PersistentFlags().StringP("lookup", "l", "", "Single IP/Domain to lookup. Has no effect when --file is also specified.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
