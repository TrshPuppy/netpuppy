def sum(a: int, b: int) -> int:
    return a + b


def banner() -> str:
    banner: str = """
TrshPuppy brings you...

|8PPPPe                    ___      .++.
|8    |8  |eeee |eeeee  __/_  `.  .'    `.
|8e   |8  |8      |8    \_,  | \_'  /   )`-')
|88   |8  |8eee   |8e    U ) `-`    \  ((`\"`
|88   |8  |88     |88    ___Y  ,    .'7 /| 
|88___|8__|88ee___|88___(_,___/___.'_(_/_/_

|8PPPPe
|8    |8 |e   .e  |eeeee  |eeeee  |e    .e
|8eeee8  |8   |8  |8   |8 |8   |8 |8    |8
|88      |8e  |8  |8eee8  |8eee8  |8eeee8
|88      |88  |8  |88     |88      |88
|88______|88ee8___|88_____|88______|88____

           Launch a puppy to
         ~ sneef  and  fetch ~
           data   for   you!
"""
    return banner


def user_selection_update(h: str, p: str, l: str) -> str:
    update: str = """
           bork!
      __  /  
 (___()'`;       |Host: {host}
 / )   /`        |Port: {port}
 /\\'--/\\         |Mode: {mode}
    """.format(
        host=h,
        port=p,
        mode="Server" if l else "Client",
    )

    return update
