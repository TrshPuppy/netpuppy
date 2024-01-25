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

### 3. Threading
[Threading](https://docs.python.org/3/library/threading.html) is a built in Python lib which constructs higher-level threading interfaces.

**Usage in Project:**
The `threading` lib is used in this project to allow for data to be sent and received from either endpoint of the connection established by NetPuppy without blocking or collision.

### 4. Subprocess
[Subprocess](https://docs.python.org/3/library/subprocess.html) is a built in libraryu which allows you to spawn processes and connect to their input, output, and error pipes.

**Usage in Project:**
The `subprocess` module is being used to help us execute simple commands on the serving endpoint of the connection.

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

- Run tests:
```
make test
```

- Clean the project (remove build artifacts, virtual environment, etc.):
```
make clean
```

Once cloned, run `make env`, then afer activating the python venv with either `source .venv/bin/activate` should be able to run NetPuppy using `netpuppy`.

### Flags
As it's being planned & built right now, NetPuppy will have two modes: listen & connect. However, once the connection is made, both ends listen as well as serve.
#### Client mode
Client is the default mode & will initiate a TCP connection to a remote host address. Just tell NetPuppy to `listen`, followed by the host IP and port.
```
netpuppy connect 10.0.2.5 69
```
#### Server/ Listener mode
To have NetPuppy listen on `0.0.0.0` and on a specific port, just tell NetPuppy `listen` followed by the port number.
```
netpuppy listen 69
```
## Why Python?
Python is the language for this tool for now, with plans to migrate to a lower-level, less-abstracted language in the future (such as Golang perhaps :) ?)
# Contribution:
Fork this repo, make some changes, and submit a PR. Please be detailed about the changes you're attempting to implement!
## BONUS!
If you want an easy way to contribute, you can always create your own NetPuppy banner and add it to the `/banners` directory. Just review what the other banners look like and match the syntax, etc.. <3
