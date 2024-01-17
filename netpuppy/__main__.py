##################### TP U ARE HERE
####################### get rid of nargs?
###################### group. add_argument: does it take the same parameters
###################### as regular parser.add_argument?
##################### check out cool-retro-term

import argparse
import socket
import ipaddress
import sys


def format_ip(ip_string):
    # Check for ipv4 vs ipv6
    print(f"ip_string = {ip_string}")
    try:
        address = ipaddress.ip_address(ip_string)
    except ValueError:
        print("Invalid IP address")
        return False


def main() -> None:
    # Make the CLI arg parser:
    parser = argparse.ArgumentParser(
        prog="netpuppy",  # ./netpuppy
        description="Launch a puppy to sneef and fetch data for you",
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

    connection = None

    # Create the socket and connection depending on the input:
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
            cs = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)
        except:
            print("Failed to create socket")
            sys.exit(1)

        # Set connection:
        connection = cs.connect((HOST_IP, int(HOST_PORT)))

        # Receive and send data:
        SEND_DATA = b""
        RECEIVE_DATA = b""

        try:
            rdata = connection.recv(1024)
            sdata = connection.send("sent tiddies")

            if rdata:
                RECEIVE_DATA += rdata
            elif sdata:
                SEND_DATA += sdata

        except KeyboardInterrupt:
            connection.close()
            sys.exit(1)
        except Exception as err:
            print(f"Unknown error: {err}")
            sys.exit(1)

    #     # Create a socket:
    #     ls = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    #     ls.bind((SERVER_IP, SERVER_PORT))
    #     ls.listen()

    #     try:
    #         connection, addr = ls.accept()

    #         data = connection.recv(1024)

    #     except KeyboardInterrupt:
    #         connection.close()
    #         ls.close()

    #         # while connection:
    #         #     CLIENT_DATA += data
    #         # print(f"DATA from CLIENT: {str(CLIENT_DATA)}")

    #     # print(f"Listening on port {args.port[0]}")

    # # When the user wants a connecting client:
    # elif args.host_ip:
    #     # BRAINDUMP
    #     #        - know the host of both machines
    #     #        - receive commands to run on the client
    #     #        - send the output of those commands back to the host

    #     HOST_IP = args.host_ip
    #     HOST_PORT = args.port[0]
    #     RECEIVED_DATA = b""

    #     # Create a connection socket:
    #     cs = None

    #     try:
    #         # Try to create an ipv4 socket and connection:
    #         cs = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    #         cs.connect((HOST_IP, int(HOST_PORT)))

    #     except socket.error as err:
    #         cs = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)

    #         # Connect socket to host:
    #         print(f"Connecting to {args.host_ip} on port {args.port[0]}")
    #         cs.connect((HOST_IP, int(HOST_PORT)))
    #     except:
    #         print("Failed to connect to host")
    #         sys.exit(1)

    #     print(f"Connecting to {HOST_IP} on port {HOST_PORT}")
    #     # Receive data from host:
    #     try:
    #         while True:
    #             host_data = cs.recv(1024)
    #             RECEIVED_DATA += host_data
    #             print(str(RECEIVED_DATA))
    #     except KeyboardInterrupt:
    #         cs.close()
