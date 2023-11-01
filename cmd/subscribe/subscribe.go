/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package subscribe

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
	nodeID   string
	timeout  int16
	interval int64
)

func getUrl() (string, error) {
	data, err := os.ReadFile(globalvar.SessionFile)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func sub(nodeId string, url string) (error, value any) {
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

	notifyCh := make(chan *opcua.PublishNotificationData)

	sub, err := c.Subscribe(ctx, &opcua.SubscriptionParameters{
		Interval: time.Duration(interval) * time.Millisecond,
	}, notifyCh)
	if err != nil {
		return err, nil
	}

	fmt.Printf("Created subscription with id %v", sub.SubscriptionID)

	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		return err, nil
	}

	var miCreateRequest *ua.MonitoredItemCreateRequest

	miCreateRequest = valueRequest(id)
	res, err := sub.Monitor(ctx, ua.TimestampsToReturnBoth, miCreateRequest)
	if err != nil || res.Results[0].StatusCode != ua.StatusOK {
		return err, nil
	}

	for {
		select {
		case <-ctx.Done():
			return
		case res := <-notifyCh:
			if res.Error != nil {
				return err, nil
			}

			return nil, res.Value
		}
	}
}

func valueRequest(nodeID *ua.NodeID) *ua.MonitoredItemCreateRequest {
	handle := uint32(42)
	return opcua.NewMonitoredItemCreateRequestWithDefaults(nodeID, ua.AttributeIDValue, handle)
}

// subscribeCmd represents the subscribe command
var SubscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if url, err := getUrl(); err != nil {
			fmt.Println(err)
		} else {

			fmt.Println(fmt.Sprintf("Entered Node-ID: %s", nodeID))

			if err, resp := sub(nodeID, url); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(resp)
			}
		}
	},
}

func init() {

	SubscribeCmd.Flags().StringVarP(&nodeID, "nodeID", "n", "", "ID of the node to read")
	SubscribeCmd.Flags().Int16VarP(&timeout, "timeout", "t", 2000, "Timeout in milliseconds")
	SubscribeCmd.Flags().Int64VarP(&interval, "interval", "i", 1000, "Subcribtion interval")

	if err := SubscribeCmd.MarkFlagRequired("nodeID"); err != nil {
		fmt.Println(err)
	}

}
