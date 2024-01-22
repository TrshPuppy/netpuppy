from netpuppy.utils import banner, user_selection_update
import argparse
import socket
import sys
from netpuppy.classes import SocketConnection


def network_port(value: str) -> int:
    ivalue = int(value)
    if ivalue <= 0 or ivalue > 65535:
        raise argparse.ArgumentTypeError("%s is an invalid port number" % value)
    return ivalue


def main() -> None:
    # Make the CLI arg parser:
    parser = argparse.ArgumentParser(
        prog="netpuppy",  # ./netpuppy
        description="Launch a puppy to sneef and fetch data for you!",
        epilog="Tell netpuppy he was a good boi.",
    )

    # Add arguments
    #   First group is mutually exclusive (listen vs host-ip/ connect)
    exclusive_group = parser.add_mutually_exclusive_group(required=True)
    exclusive_group.add_argument("-l", "--listen", action="store_true")
    exclusive_group.add_argument("-H", "--host-ip", action="store")

    parser.add_argument(
        "-p", "--port", action="store", nargs=1, required=True, type=network_port
    )

    # Get the list of arguments:
    try:
        # If netpuppy was executed w/ valid args, print the banner:
        args = parser.parse_args()
        print(banner())
    except argparse.ArgumentError:
        sys.exit(1)

    # Print the args:
    print(user_selection_update(args.host_ip, args.port[0], args.listen))

    # Create the socket and connection depending on the input:
    print("Creating soocket...")
    connection: socket.socket | None = None

    # Server mode
    if args.listen:
        SERVER_IP: str = "0.0.0.0"
        SERVER_PORT: int = int(args.port[0])

        # Create a socket:
        ls: socket.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        ls.bind((SERVER_IP, SERVER_PORT))
        ls.listen()

        # Set connection:
        print("Trying connection...")
        connection, addr = ls.accept()

    # Client mode
    elif args.host_ip:
        HOST_IP: str = args.host_ip
        HOST_PORT: int = args.port[0]

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
        current_connection: SocketConnection = SocketConnection(connection, peer[0])

        # Update user:
        print(f"Connection established to: {current_connection.address} port {peer[1]}")

        try:
            while connection:
                current_connection.read_stream()
                current_connection.write_stream()

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
