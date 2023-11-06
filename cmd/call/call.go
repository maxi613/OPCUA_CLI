/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package call

import (
	globalvar "cli-tool/globalVar"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/spf13/cobra"
)

var (
	methodID string
	objectID string
	timeout  int16
)

func getUrl() (string, error) {
	data, err := os.ReadFile(globalvar.SessionFile)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func callFunction(url string, methodID string, objectID string) error {
	ctx := context.Background()

	c, err := opcua.NewClient(url, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err != nil {
		return err
	}
	if err := c.Connect(ctx); err != nil {
		return err
	}
	defer c.Close(ctx)

	req := &ua.CallMethodRequest{
		ObjectID:       ua.NewStringNodeID(2, "main"),
		MethodID:       ua.NewStringNodeID(2, "even"),
		InputArguments: []*ua.Variant{},
	}

	resp, err := c.Call(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	if got, want := resp.StatusCode, ua.StatusOK; got != want {
		log.Fatalf("got status %v want %v", got, want)
	}

	return nil
}

// callCmd represents the call command
var CallCmd = &cobra.Command{
	Use:   "call",
	Short: "Executes an opc ua function",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if url, err := getUrl(); err != nil {
			fmt.Println(err)
		} else {

			fmt.Println(fmt.Sprintf("Entered Node-ID: %s", methodID))

			if err := callFunction(url, methodID, objectID); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	CallCmd.Flags().StringVarP(&methodID, "methodID", "m", "", "ID of the function to call")
	CallCmd.Flags().StringVarP(&objectID, "objectID", "o", "", "Object-ID of the output argument")
	CallCmd.Flags().Int16VarP(&timeout, "timeout", "t", 2000, "Timeout in milliseconds")
	if err := CallCmd.MarkFlagRequired("methodID"); err != nil {
		fmt.Println(err)
	}

	if err := CallCmd.MarkFlagRequired("objectID"); err != nil {
		fmt.Println(err)
	}
}
