/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"cli-tool/cmd"
	globalvar "cli-tool/globalVar"
	"fmt"
	"os"
)

func main() {
	if data, err := os.ReadFile(globalvar.SessionFile); err != nil && len(data) == 0 {
		fmt.Println("No connection is registered. Register a connection with the subcommand connect")
	} else {

		output := fmt.Sprintf("A session with the endpoint %s will be opened.\nInfo: To change the session use the subcommand connect.", string(data))

		fmt.Println(output)
	}
	cmd.Execute()
}
