import random
import os


def sum(a: int, b: int) -> int:
    return a + b


def get_random_banner() -> str:
    working_dir: str = os.path.dirname(os.path.abspath(__file__))
    banner_dir: str = os.path.join(working_dir, "banners")
    banners: list[str] = [
        file for file in os.listdir(banner_dir) if file.endswith(".txt")
    ]

    if not banners:
        print("No banners found :~(")
        return ""

    chosen_one: str = os.path.join(banner_dir, random.choice(banners))

    with open(chosen_one, "r") as file:
        banner_to_return: str = file.read()
        return banner_to_return


def banner() -> str:
    random_banner: str = get_random_banner()
    banner: str = f"""
TrshPuppy brings you...

{random_banner.strip()}

           Launch a puppy to
         ~ sneef  and  fetch ~
           data   for   you!
"""
    return banner


def user_selection_update(h: str, p: str, c: str) -> str:
    update: str = ""

    if c == "connect":
        update = """
           bork!
      __  /  
 (___()'`;      |Host: {host}
 / )   /`       |Port: {port}
 /\\'--/\\        |Mode: {mode}
    """.format(
            host=h,
            port=p,
            mode="Client",
        )

    else:
        update = """
    .-.  *sneef sneef*
   / (_   
  ( "  6\\___o   |Host: {host}
  /  (  ___/    |Port: {port}
 /     /  U     |Mode: {mode}
    """.format(
            host=h,
            port=p,
            mode="Server",
        )

    return update
