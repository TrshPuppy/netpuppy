# U ARE HERE:
#      focusing on creating connection which can send and receive data
#      fix sending data portion/ logic
#      touch up receiving data portion/ logic

import argparse
import socket
import sys


def main():
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

    parser.add_argument("-p", "--port", action="store", nargs=1, required=True)

    # Get the list of arguments:
    args = parser.parse_args()
    print(f" args = {args}")

    # Create the socket and connection depending on the input:
    connection = None

    if args.listen:
        SERVER_IP = "0.0.0.0"
        SERVER_PORT = 44440

        # Create a socket:
        ls = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        ls.bind((SERVER_IP, SERVER_PORT))
        ls.listen()

        # Set connection:
        connection, addr = ls.accept()

    elif args.host_ip:
        HOST_IP = args.host_ip
        HOST_PORT = args.port[0]

        # Create a socket:
        cs = None

        try:
            # Try to create an ipv4 socket and connection:
            cs = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        except socket.error as err:
            # If IPv4 fails, try an  IPv6 socket:
            cs = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)

        except:
            print("Failed to create socket")
            sys.exit(1)

        # Set connection:
        connection = cs.connect((HOST_IP, int(HOST_PORT)))

    # Receive and send data:
    SEND_DATA = b""
    RECEIVE_DATA = b""

    while True:
        try:
            rdata = connection.recv(1024)
            sdata = "send test data"

            if rdata:
                RECEIVE_DATA += rdata
                print(f"Received data: {str(RECEIVE_DATA)}")
            else:
                SEND_DATA += sdata

        except KeyboardInterrupt:
            connection.close()
            sys.exit(1)
        except Exception as err:
            print(f"Unknown error: {err}")
            sys.exit(1)
