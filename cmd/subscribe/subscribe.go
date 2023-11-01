/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscribe called")
	},
}

func init() {

}
