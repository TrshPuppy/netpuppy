# NetPuppy
A CLI tool for making and/ or binding TCP connections. Inspired by [Netcat](https://netcat.sourceforge.net/) & written in Python!
## Goals
NetPuppy will be able to:
- listen for & serve incoming TCP connections as well as initiate outgoing ones
- maintain a stable connections b/w both parties
- send and receive data from either endpoint
## Use:

## Python Packages Used

This project utilizes several built-in Python packages to achieve its functionality. Here's a brief overview of these packages:

### 1. Argparse

[Argparse](https://docs.python.org/3/library/argparse.html) is a built-in Python library used for parsing command-line arguments. It provides a way to specify the inputs a program expects and generates helpful error messages for users.

**Usage in Project:** 
The `argparse` library is used to handle command-line arguments for the `netpuppy` tool, allowing users to specify options such as listening mode, host IP, and port.

### 2. Sockets

[Sockets](https://docs.python.org/3/library/socket.html) is a built-in Python library that provides a low-level networking interface. It's used for creating and working with both server and client sockets.

**Usage in Project:** 
The `socket` library is employed to establish network connections, send, and receive data.

## Development Dependencies

This project uses several development tools to ensure code quality, type safety, and to run tests. Here's a brief overview of these tools:

### 1. Black

[Black](https://black.readthedocs.io/en/stable/) is a code formatter that ensures consistent code formatting throughout the project.

**Usage:**
```
make format
```

### 2. Mypy

[Mypy](http://mypy-lang.org/) is a static type checker for Python. It helps catch type errors before runtime.

**Usage:**
```
make test
```

### 3. Pytest

[Pytest](https://docs.pytest.org/en/stable/) is a testing framework that makes it easy to write simple and scalable test cases.

**Usage:**
```
make test
```

## Dependency on Make

This project uses `make` to manage various tasks such as setting up the development environment, formatting code, running tests, and more. Ensure you have `make` installed on your system to use the provided `Makefile`.

**Examples of commands you can run:**

- Set up the virtual environment:
```
make venv
```

- Format the code:
```
make format
```

- Run tests & mypy:
```
make test
```

- Clean the project (remove build artifacts, virtual environment, etc.):
```
make clean
```

Once cloned, run `make env`, then afer activating the python venv with either `source .venv/bin/activate` or `source .venv\Scripts\activate` should be able to run NetPuppy using `netpuppy`.

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
