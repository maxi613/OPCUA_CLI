/*
Copyright © 2023 wiegandmaximilian@gmail.com
*/
package read

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
)

func getUrl() (string, error) {
	data, err := os.ReadFile(globalvar.SessionFile)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func readNode(url string, nodeId string) (error, any) {
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

	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	resp, err := c.Read(ctx, req)

	if err != nil {
		return err, nil
	}

	if resp != nil && resp.Results[0].Status != ua.StatusOK {
		errorMessage := fmt.Sprintf("Status not OK: %s", resp.Results[0].Status)
		return fmt.Errorf(errorMessage), nil
	}

	return nil, resp.Results[0].Value.Value()
}

// readCmd represents the read command
var ReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Function to read a Node-ID",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if url, err := getUrl(); err != nil {
			fmt.Println(err)
		} else {

			fmt.Println(fmt.Sprintf("Entered Node-ID: %s", nodeID))

			if err, value := readNode(url, nodeID); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(value)
			}
		}
	},
}

func init() {
	ReadCmd.Flags().StringVarP(&nodeID, "nodeID", "n", "", "ID of the node to read")
	ReadCmd.Flags().Int16VarP(&timeout, "timeout", "t", 2000, "Timeout in milliseconds")
	if err := ReadCmd.MarkFlagRequired("nodeID"); err != nil {
		fmt.Println(err)
	}
}
