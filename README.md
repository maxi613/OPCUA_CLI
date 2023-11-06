# OPCUA CLI
A CLI-Tool to test OPC functions. 

Based on the native Go library: https://github.com/gopcua/opcua

Test the new fist version under: https://github.com/maxi613/OPCUA_CLI/tags

## Quick Start
### Connect

To start with testing some functions you have to register first a connection. 
You can do this by entering the following code in the command line:

`cli-tool connect -u opc.tcp://192.168.0.1`

If it was succesfull you will get a message back. 

### Write

So far it is only possible to write an variable with the type of integer. 
Use the subcommand `write` with the flag `-n` to tell the Node-ID and `-v` to enter the value. 
If the operation was successful you will get back the statuscode `0x0`.

`cli-tool write -n "ns=4;i=7" -v 2`

### Read

It is possible to read all type of varliables. 
You only have to specify the Node-ID with the flag `-n`

`cli-tool read -n "ns=4;i=7"`

### Subscribe 
With a subscribtion is possible to monitor all kind of variables. Enter `subscribe` as subcommand with the flag `-n` to specify again the Node-ID and with the flag `-i` you can enter an monitor interval. 
The standard value for interval is 1000. The unit is in milliseconds. 

`cli-tool subscribe -n "ns=4;i=7" -i 10`

## Status Features

| Feature  | Stauts |
| :-------------: | :-------------: |
| Read  | X|
| Write  | X |
| Subscribtion| X |
| Call Functions| Planend|
|Browse Nodes| Planned|
|Browse Endpoints|Planned|









