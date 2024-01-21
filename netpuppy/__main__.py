from netpuppy.utils import banner, user_selection_update
import argparse
import socket
import sys


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
            # print(f"IPv4 failed: {err}")
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

    # Receive and send data:
    SEND_DATA: bytes = b""
    RECEIVE_DATA: bytes = b""

    if connection:
        peer = connection.getpeername()
        peername: str = peer[0]
        peer_port: str = peer[1]

        try:
            print(f"Connection established to: {peername} port {peer_port}")
            while connection:
                # Data received:
                rdata: bytes = connection.recv(1024)

                # Send data (inout from user):
                input_chunk: str = input()

                SEND_DATA += input_chunk.encode("utf-8")
                input_chunk = ""
                # sudo code:
                #       for every loop of the while loop
                #           check for user input (data to send)
                #       * probs can't send and receive at the same time?

                if len(rdata) > 0:
                    RECEIVE_DATA += rdata
                    rdata = b""
                    print(f"Received data: {str(RECEIVE_DATA)}")
                elif len(SEND_DATA) > 0:
                    # Send the data

                    print("tiddies else thing line 96")
                    # SEND_DATA += sdata

        except KeyboardInterrupt:
            print(f"Keyboard interrupt: {KeyboardInterrupt}")

        except Exception as err:
            print("Connection failed.")
            print(f"Unknown error: {err}")

        finally:
            if connection:
                connection.close()

    else:
        print("Connection failed.")
        sys.exit(1)
