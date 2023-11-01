/*
Copyright © 2023 Maximilian Wiegand HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cli-tool/cmd/connect"
	"cli-tool/cmd/read"
	"cli-tool/cmd/subscribe"
	"cli-tool/cmd/write"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "An opc ua cli client",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func addSubcommands() {
	rootCmd.AddCommand(connect.ConnectCmd)
	rootCmd.AddCommand(read.ReadCmd)
	rootCmd.AddCommand(write.WriteCmd)
	rootCmd.AddCommand(subscribe.SubscribeCmd)
}

func init() {

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addSubcommands()
}
