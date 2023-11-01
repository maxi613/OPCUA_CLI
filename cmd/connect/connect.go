/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package connect

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/spf13/cobra"
)

var (
	url     string
	timeout int16
)

func registerConnection(url string) error {
	url_b := []byte(url)
	err := os.WriteFile("session.dat", url_b, 0644)

	if err != nil {
		return err
	}

	return nil
}

func connect(url string) (error, bool) {

	durationTimeout := time.Duration(timeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), durationTimeout)

	defer cancel()

	output := fmt.Sprintf("Try to connect to %s\n", url)
	fmt.Println(output)
	c, err := opcua.NewClient(url, opcua.SecurityMode(ua.MessageSecurityModeNone))
	defer c.Close(ctx)
	if err != nil {
		return err, false
	}

	if err := c.Connect(ctx); err != nil {
		return err, false
	}

	fmt.Println(os.Getenv("URL"))
	os.Setenv("URL", url)
	return nil, true
}

// connectCmd represents the connect command
var ConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect with a OPC UA Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err, success := connect(url); err != nil && success != true {
			fmt.Println(err)
			output := fmt.Sprintf("The connection to the Endpoint %s failed", url)
			fmt.Println(output)

		} else {
			output := fmt.Sprintf("The connection to the Endpoint %s was succesfull", url)
			fmt.Println(output)
		}

		if err := registerConnection(url); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Connection is registered")
		}
	},
}

func init() {

	ConnectCmd.Flags().StringVarP(&url, "url", "u", "", "The url of the server")
	ConnectCmd.Flags().Int16VarP(&timeout, "timeout", "t", 2000, "Timeout in milliseconds")
	if err := ConnectCmd.MarkFlagRequired("url"); err != nil {
		fmt.Println(err)
	}
}
