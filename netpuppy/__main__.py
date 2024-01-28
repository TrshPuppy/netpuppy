# YOU ARE HERE:
#  So the only issue is bottom window (listening process (we didn't give it a subprocess))
# stdin not being input into the opened subprocess
# maybe use rev instead of date


from netpuppy.utils import banner, user_selection_update
import argparse
import socket
import sys
from netpuppy.classes import SocketConnection
import asyncio


def network_port(value: str) -> int:
    ivalue = int(value)
    if ivalue <= 0 or ivalue > 65535:
        raise argparse.ArgumentTypeError("%s is an invalid port number" % value)
    return ivalue


async def main() -> None:
    # Make the CLI arg parser:
    parser = argparse.ArgumentParser(
        prog="netpuppy",  # ./netpuppy
        description="Launch a puppy to sneef and fetch data for you!",
        epilog="Tell netpuppy he was a good boi.",
    )

    # Add arguments as individual subparsed groups - this gains us the ability to easily extend our 'commands' in the future
    # as well as fine grained control over the arguments for each command. This does come with a slight drawback of
    # having to repeat shared argument definitions (like port), but seems worth the tradeoff
    subparsers = parser.add_subparsers()

    sp = subparsers.add_parser("connect", help="Connect to a host (Client Mode)")
    sp.set_defaults(cmd="connect")
    sp.add_argument("host_ip", help="Host IP Address", type=str)  # required input
    sp.add_argument(
        "port", help="Host Port", type=network_port, nargs="?", default="44440"
    )  # optional with default

    sp = subparsers.add_parser("listen", help="Listen on a port (Server Mode)")
    sp.set_defaults(cmd="listen")
    sp.set_defaults(host_ip="0.0.0.0")  # default to all interfaces
    sp.add_argument(
        "port", help="Listen Port", type=network_port, nargs="?", default="44440"
    )  # optional with default

    # Get the list of arguments:
    try:
        # If netpuppy was executed w/ valid args, print the banner:
        args = parser.parse_args()
        print(banner())
    except argparse.ArgumentError:
        sys.exit(1)

    # Print the args:
    print(user_selection_update(args.host_ip, args.port, args.cmd))

    # Create the socket and connection depending on the input:
    print("Creating soocket...")
    connection: socket.socket | None = None
    sub_proc: None | int = 0

    # Server mode
    if args.cmd == "listen":
        SERVER_IP: str = "0.0.0.0"
        SERVER_PORT: str = args.port
        PEER_TYPE: str = "offense"

        # Create a socket:
        ls: socket.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        ls.bind((SERVER_IP, int(SERVER_PORT)))
        ls.listen()

        # Set connection:
        print("Trying connection...")
        connection, addr = ls.accept()

    # Client mode
    elif args.cmd == "connect":
        HOST_IP: str = args.host_ip
        HOST_PORT: str = args.port
        PEER_TYPE: str = "connect_back"

        # Create a socket:
        cs: socket.socket | None = None

        try:
            # Try to create an ipv4 socket and connection:
            cs = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        except socket.error as err:
            # If IPv4 fails, try an  IPv6 socket:
            try:
                cs = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)
                print("connection ipv 6 block")
            except socket.error as errrrrr:
                print(f"IPv6 failed: {errrrrr}")
                sys.exit(1)

        except:
            print("Failed to create socket")
            sys.exit(1)

        # Set connection:
        print("Trying connection...")
        connection = cs

        connection.connect((HOST_IP, HOST_PORT))

    if connection:
        # Get peer name and port, create SocketConnection object:
        peer = connection.getpeername()  # [peername, peer port]
        current_connection: SocketConnection = SocketConnection(
            connection, peer[0], PEER_TYPE
        )
        socket_file_no: int = (
            current_connection.socket.fileno()
        )  # the file descriptor of the socket

        # Depending on the socket "type", start a subprocess:
        if (
            current_connection.type == "offense"
        ):  #  An 'offense' peer is the one listening for a connect_back (starts as the listener)
            # sub_proc = subprocess.call
            print("The peer type is offense")
        elif (
            current_connection.type == "connect_back"
        ):  # A 'connect_back' peer is sending data from the machine it's on back to the offense peer
            print("The peer type is connect_back")

        print(
            f"Connection established to: {current_connection.ILoveOCaml} port {peer[1]}"
        )

        try:
            while connection:
                # asyncio.eventloop
                gather = await asyncio.gather(
                    current_connection.read_stream(), current_connection.write_stream()
                )

                asyncio.run(gather)

                # await current_connection.read_stream()
                # await current_connection.write_stream()

                # print(f"in main: {current_connection.sub_proc.stdout.read()}")

        except KeyboardInterrupt:
            print(f"Keyboard interrupt: {KeyboardInterrupt}")
            current_connection.running = False

        except Exception as err:
            print("Connection failed.")
            print(f"Unknown error: {err}")
            current_connection.running = False

        finally:
            if connection:
                connection.close()

    else:
        print("Connection failed.")
        sys.exit(1)


def lil_wayne():
    asyncio.run(main())


lil_wayne()

# SUDO CODE: (what OS)
# nmap fingerprint
# WHAT DOES TP WANT!??!!??
#       - prompt > (local user): current directory
#       - message on connection which tells the user:
#           - what os
#           - what directory?
#  SHELL OPTIONS:
#           Linux: files (if the target is linux)
#               'open' a /bin/bash process file
#                   return stdout back
#                       pipe our input into stdin
#          Python: check for python on the system
#               make a python shell
#
# ""

# BASH's JOB?
#   - find bash
#   - call him up
#   - subprocess lib?


# PYTHON'S JOB?
#


#        SERVERS JOB?
# start a shell/ subprocess thing to gather info about the target based on the
# input from the client
#

#
# 100,000 feet:
#   connect to the other computer target
#   report back some env info like the os,
#
#   50, 000 feet
#   either:"
#         a command which will work regardless of the OS (python)
#           if else:
#               if its not linux
#                    its something shitty
#
# STRUCTURE:
