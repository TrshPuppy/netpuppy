##################### TP U ARE HERE
####################### get rid of nargs?
###################### group. add_argument: does it take the same parameters
###################### as regular parser.add_argument?
##################### check out cool-retro-term

import argparse
import socket


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
    print(args)
