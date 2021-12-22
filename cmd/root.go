package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crypto-client",
	Short: "Interact with different crypto APIs.",
	Long: `Crypto-Client is a cli client for interacting with different crypto currency service providers. 

To see more options run this command with the api provider of your choice followed by the -h flag. 
The example below will output the help message for interacting with the Coinbase API.

	$ cyrpto-client coinbase -h

Supported APIs are listed in the table below.

	╔══════════╤══════════════════╗
	║ Provider │ Supported        ║
	╠══════════╪══════════════════╣
	║ Coinbase │ partial          ║
	╟──────────┼──────────────────╢
	║ Celsius  │ TBD              ║
	╚══════════╧══════════════════╝

Please note that if the vendor makes breaking changes to their API it could break the cypto-client cli.
`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
