/*
Copyright © 2023 wiegandmaximilian@gmail.com
*/
package subscribe

import (
	globalvar "cli-tool/globalVar"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
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

func sub(nodeId string, url string) error {
	ctx := context.Background()
	c, err := opcua.NewClient(url, opcua.SecurityMode(ua.MessageSecurityModeNone))

	if err != nil {
		return err
	}
	if err := c.Connect(ctx); err != nil {
		return err
	}

	defer c.Close(ctx)

	m, err := monitor.NewNodeMonitor(c)
	if err != nil {
		return err
	}

	m.SetErrorHandler(func(_ *opcua.Client, sub *monitor.Subscription, err error) {
		log.Printf("error: sub=%d err=%s", sub.SubscriptionID(), err.Error())
	})

	go startChanSub(ctx, m, time.Duration(interval), 0, nodeId)

	<-ctx.Done()
	return nil
}

func startChanSub(ctx context.Context, m *monitor.NodeMonitor, interval, lag time.Duration, nodes ...string) {
	ch := make(chan *monitor.DataChangeMessage, 16)
	sub, err := m.ChanSubscribe(ctx, &opcua.SubscriptionParameters{Interval: interval}, ch, nodes...)

	cleanChannel := make(chan os.Signal, 1)
	signal.Notify(cleanChannel, os.Interrupt)

	go func() {
		<-cleanChannel
		log.Printf("unsubscribe %d", sub.SubscriptionID())
		cleanup(ctx, sub)
		os.Exit(1)
	}()

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-ch:
			if msg.Error != nil {
				log.Printf(" sub=%d error=%s", sub.SubscriptionID(), msg.Error)
			} else {
				log.Printf(" sub=%d ts=%s node=%s value=%v", sub.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())
			}
			time.Sleep(lag)
		}
	}
}

func cleanup(ctx context.Context, sub *monitor.Subscription) {
	log.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	sub.Unsubscribe(ctx)
}

var SubscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if url, err := getUrl(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(fmt.Sprintf("Entered Node-ID: %s", nodeID))

			if err := sub(nodeID, url); err != nil {
				fmt.Println(err)
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
