package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "whoisenum",
	Version: "0.0.1a",
	PersistentPreRunE: func(ctx *cobra.Command, args []string) error {
		l, _ := ctx.Flags().GetString("lookup")
		f, _ := ctx.Flags().GetString("file")
		if l != "" && f != "" {
			fmt.Printf("Both --lookup and --file provided. --file will take priority and --lookup will be ignored.")
			//return errors.New("options --lookup and --file are mutually exclusive. Please only provide one of these options to continue")
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(ipCmd, domainCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON in place of the default format")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet extra program output and only print results")
	rootCmd.PersistentFlags().StringP("file", "f", "", "File with single IP/Domain per line")
	rootCmd.PersistentFlags().StringP("lookup", "l", "", "Single IP/Domain to lookup. Will not be used when --file is specified.")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Include verbose messages from program execution")
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.whoisenum.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
