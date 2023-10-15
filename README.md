# NetPuppy
A CLI tool for making and/ or binding TCP connections. Inspired by [Netcat](https://netcat.sourceforge.net/) & written in Python!
## Goals
NetPuppy will be able to:
- listen for & serve incoming TCP connections as well as initiate outgoing ones
- maintain a stable connections b/w both parties
- send and receive data from either endpoint
## Use:
The only dependencies at this time are:
- [argparse](https://docs.python.org/3/library/argparse.html)
- & [socket](https://docs.python.org/3/library/socket.html?highlight=socket#module-socket)

Once cloned, run `python -m netpuppy-pkg`, then you should be able to run NetPuppy using `./netpuppy`.
### Flags
As it's being planned & built right now, NetPuppy will have two modes: client & server. 
#### Client mode
Client is the default mode & will initiate a TCP connection to a remote host address (`-H`) and port number (`-p`):
```
./netpuppy -H 10.0.2.5 -p 69
```
#### Server/ Listener mode
Give the `-l` flag to put NetPuppy into listening mode. In this mode NetPuppy will default to your primary network interface and will listen to the posrt you provide via `-p`.
```
./netpuppy -l -p 69
```
## Why Python?
Python is the language for this tool for now, with plans to migrate to a lower-level, less-abstracted language in the future (such as Golang perhaps :) ?)
