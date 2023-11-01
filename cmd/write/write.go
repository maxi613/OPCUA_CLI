/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package write

import (
	globalvar "cli-tool/globalVar"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/spf13/cobra"
)

var (
	nodeID  string
	timeout int16
	value   int16
)

func getUrl() (string, error) {
	data, err := os.ReadFile(globalvar.SessionFile)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writeNodeID(url string, nodeId string) (error, any) {
	durationTimeout := time.Duration(timeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), durationTimeout)
	defer cancel()
	c, err := opcua.NewClient(url, opcua.SecurityMode(ua.MessageSecurityModeNone))

	if err != nil {
		return err, nil
	}
	if err := c.Connect(ctx); err != nil {
		return err, nil
	}

	defer c.Close(ctx)

	id, err := ua.ParseNodeID(nodeId)
	if err != nil {
		return err, nil
	}

	v, err := ua.NewVariant(value)
	if err != nil {
		return err, nil
	}

	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					EncodingMask: ua.DataValueValue,
					Value:        v,
				},
			},
		},
	}

	resp, err := c.Write(ctx, req)
	if err != nil {
		return err, nil
	}
	return nil, resp.Results[0]
}

// writeCmd represents the write command
var WriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a Value to an Node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if url, err := getUrl(); err != nil {
			fmt.Println(err)
		} else {

			fmt.Println(fmt.Sprintf("Entered Node-ID: %s", nodeID))

			if err, resp := writeNodeID(url, nodeID); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(resp)
			}
		}
	},
}

func init() {
	WriteCmd.Flags().StringVarP(&nodeID, "nodeID", "n", "", "ID of the node to read")
	WriteCmd.Flags().Int16VarP(&timeout, "timeout", "t", 2000, "Timeout in milliseconds")
	WriteCmd.Flags().Int16VarP(&value, "value", "v", 0, "Value to write")

	if err := WriteCmd.MarkFlagRequired("nodeID"); err != nil {
		fmt.Println(err)
	}

	if err := WriteCmd.MarkFlagRequired("value"); err != nil {
		fmt.Println(err)
	}

}
