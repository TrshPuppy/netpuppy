# NetPuppy
A CLI tool for making and/ or binding TCP connections. Inspired by [Netcat](https://netcat.sourceforge.net/) & written in Golang!
## Goals
NetPuppy will be able to:
- listen for & serve incoming TCP connections as well as initiate outgoing ones
- maintain a stable connections b/w both parties
- send and receive data from either endpoint
## Use:
NetPuppy has two primary modes: offense & connect-back.
### Offense
The 'offensive' peer is executed as a server and listens to the `0.0.0.0` address and a specified port. It will bind incoming TCP connections to that port.
### Connect-Back
The 'connect-back' peer starts w/ a client-like relationship to the offense peer. I.e. it will connect to the address and port you give it. However, once connected, it will start a subprocess on the local machine. For now, this subprocess is a bash shell. This shell will take input from the offensive peer and execute commands on the machine its running on. The output from the shell will then be echoed back to the offensive peer.
#### Flags:
- `-H` the host IP address you want to connect to (in connect-back mode)
- `-p` the port you want to start your peer on (both mode types)
- `-l` tell NetPuppy to listen, this will start NP in the offense mode. You can also give a port number.
## Examples:
### Offense Peer:
```
$ go run main.go -l -p 44444

... <banner>          

    *sneef sneef*
   .-.
  / (_          |Host:  0.0.0.0
 ( "  6\___o    |RPort: 44444
 /  (  ___/     |LPort: 44444
/     /  U      |Mode:  Offensive Server
```
### Connect-Back Peer:
```
$ go run main.go -H 0:0:0:0:0:0:0:1 -p 44444

...<banner>

        bork!
     __  /     |Host:  0:0:0:0:0:0:0:1
(___()'';      |RPort: 44444
/ )   /'       |LPort: 60804
/\'--/\        |Mode:  Client
```
## This project is still being built & improved!
### Contributing:
Just create a pull request and I (TrshPuppy) will review it.
### Python?
This project was originally written in Python. If you'd like to fork the Python branch (which isn't being updated), you can check it out [here](https://github.com/TrshPuppy/NetPuppy/tree/python-version-abandoned)!
