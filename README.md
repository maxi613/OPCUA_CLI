# OPCUA CLI
A CLI-Tool to test OPC functions. 

Based on the native Go library: https://github.com/gopcua/opcua

Test the new fist version under: https://github.com/maxi613/OPCUA_CLI/tags

## Connect

To start with testing some functions you have to register first a connection. 
You can do this by entering the following code in the command line:

`cli-tool connect -u opc.tcp://192.168.0.1`

If it was succesfull you will get a message back. 

## Write

So far it is only possible to write an variable with the type of integer 
Use the subcommand `write` with the flag `-n` to tell the Node-ID and `-v` to enter the value. 
If the operation was successful you will get back the statuscode `0x0`




