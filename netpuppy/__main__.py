##################### TP U ARE HERE
# Get a hostname, throw it to a IP-lookup library.
# A good library should just return the IP-address,
# if you give it an IP-address instead of a hostname.
# (Note: One hostname can have multiple IPv4 and/or IPv6 addresses,
#  so prepare to try all of them until you can make a connection.)
# find library to help with domain name resolution
#   OTHER:
#       mypy complains about types (check this/ K.I.M.)
#       unit tests for argparse
#       dev branch
#           protect pushed to main (need to validatae)\
#           CICD w/ github actions
#           formatting tests

import argparse
import socket
import ipaddress


def main() -> None:
    # Get the Args from command line:
    cli_args = get_parse_args()
    print(f"args in main: {cli_args}")

    return 0


def get_parse_args():
    parser = argparse.ArgumentParser(
        prog="netpuppy",  # ./netpuppy
        description="Launch a puppy to sneef and fetch data for you",
        epilog="Tell netpuppy he was a good boi.",
    )

    # This class handles input for the --port flag which has multiple valid formats
    class CustomPortAction(argparse.Action):
        print(f"custom action outside call")
        # PORT NOTES:
        #   Limit individual port range
        #   Allow for a list of ports
        #       should be: range, or list (ie. 1-10 vs 1,5,7)

        #    parser = ArgumentParser(prog='netpuppy', usage=None, description='Launch a puppy to sneef and fetch data for you', formatter_class=<class 'argparse.HelpFormatter'>, conflict_handler='error', add_help=True): namespace = Namespace(listen=True, host_ip=None, port=None): values= ['20-55']: option_string = -p
        #    Namespace(listen=True, host_ip=None, port=None)

        # ArgParser enters here:
        def __call__(self, parser, namespace, values, option_string=None):
            print(f"custom actiont")

            # Set self.values to values for --port flag stored from argparse
            setattr(self, "values", values)

            # ULTIMATEL GOAL OF THIS CLASS:
            #   Set values on namespace object to correct values
            setattr(namespace, "port", "tiddies")

            return 0

    # Add arguments
    #   First group is mutually exclusive (listen vs host-ip/ connect)
    exclusive_group = parser.add_mutually_exclusive_group(required=True)
    exclusive_group.add_argument("-l", "--listen", action="store_true")
    exclusive_group.add_argument("-H", "--host-ip", action="store")
    # HOST NOTES:
    #   Limit to real IP addresses
    #   Return and error if the IP is not valid
    #   be able to take both IP addresses and/ or hostnames
    #       IP vs Hostname
    #           argparse probs cant check the input
    #           INSTEAD:
    #              regex the input for IP format (x.x.x.x) OR (x:x:x:x:x:x:x:x)
    #                   handle ipv4 and ipv6 formats (multiple)
    #                   handle ipv4 embedded into ip4 6 ETC
    #              + (if the input fails that regex check, treat it as
    #                   a possible domain/ hostname instead)

    # Regular arguments (not grouped)
    parser.add_argument(
        "-p", "--port", action=CustomPortAction, required=True, nargs="+"
    )

    parser.add_argument("-v", "--version", action="version", version="NetPuppy v1.0")

    # Get the list of arguments:
    args = parser.parse_args()
    print(args)

    return args
