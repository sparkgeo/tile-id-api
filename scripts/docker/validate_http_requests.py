import requests
from typing import Dict, Final, List
from random import randint
from time import sleep

request_root: Final[str] = "http://api"
tile_paths: Final[Dict[str, List[str]]] = {
    "xyz": [
        "0/0/0",
        "12/234/3245",
        "7/37/42",
        "11/1851/909",
        "22/3013344/1482687",
        "17/77342/34027",
    ],
    "tms": [
        "0/0/0",
        "12/234/850",
        "7/37/85",
        "11/1851/1138",
        "22/3013344/2711616",
        "17/77342/97044",
    ],
    "quadkey": [
        "",
        "220031303212",
        "0302121",
        "13320113213",
        "1213213113323231322222",
        "12010131022213132",
    ],
}


def wait() -> None:
    print("...waiting for server")
    sleep(1)

while True:
    try:
        response = requests.get(request_root)
        if response.status_code == 404:
            print("server available")
            break
        else:
            wait()
    except Exception:
        wait()

for tile_type, paths in tile_paths.items():
    for path in paths:
        headers: Dict[str, str] = {}
        query_string: str = ""
        extension_switch = randint(0, 3)
        extension: str
        extension_supports_opacity: bool
        if extension_switch == 0:
            extension = ""
            extension_supports_opacity = True
        elif extension_switch == 1:
            extension = ".png"
            extension_supports_opacity = True
        elif extension_switch == 2:
            extension = ".jpg"
            extension_supports_opacity = False
        else:
            extension = ".jpeg"
            extension_supports_opacity = False
        if extension_supports_opacity:
            opacity = randint(0, 100)
            opacity_type_switch = randint(0, 1)
            if opacity_type_switch == 0:
                query_string = f"?opacityPercent={opacity}"
            else:
                headers["X-Opacity-Percent"] = str(opacity)
        tile_url = f"{request_root}/{tile_type}/{path}{extension}{query_string}"
        print(f"requesting {tile_url} with headers {headers}")
        response = requests.get(
            tile_url,
            headers=headers,
        )
        if response.status_code != 200:
            raise Exception("Unexpected response code", response)
