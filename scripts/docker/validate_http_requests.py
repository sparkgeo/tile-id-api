from random import randint
from time import sleep
from typing import Dict, Final, List

import requests

api_request_root: Final[str] = "http://api:8080"
wms_request_root: Final[str] = "http://mapproxy:8080"
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
    print("...waiting for servers")
    sleep(1)


api_ready = False
wms_ready = False
while True:
    try:
        api_response = requests.get(f"{api_request_root}/nonexistent")
        if api_response.status_code == 404:
            print("API server available")
            api_ready = True
        wms_response = requests.get(wms_request_root)
        if wms_response.status_code == 200:
            print("WMS server available")
            wms_ready = True
    except Exception:
        pass

    if api_ready and wms_ready:
        break
    else:
        wait()


if requests.get(f"{api_request_root}/openapi.yml").status_code != 200:
    raise Exception("openapi.yml not found")
root_response = requests.get(f"{api_request_root}/", allow_redirects=False)
if root_response.status_code < 300 or root_response.status_code >= 400:
    raise Exception("root redirect not found")
if root_response.headers["location"] != "/docs/":
    raise Exception("docs not found")
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
        tile_url = f"{api_request_root}/{tile_type}/{path}{extension}{query_string}"
        print(f"requesting {tile_url} with headers {headers}")
        api_response = requests.get(
            tile_url,
            headers=headers,
        )
        if api_response.status_code != 200:
            raise Exception("Unexpected response code", api_response)

for wms_layer in ["xyz", "tms", "quadkey"]:
    wms_url = "".join(
        [
            wms_request_root,
            "/wms",
            "?SERVICE=WMS",
            "&VERSION=1.3.0",
            "&REQUEST=GetMap",
            "&BBOX=-14091244.36530674994,7425609.323962658644,-14090591.39260812104,7426122.497220084071",
            "&CRS=EPSG:3857",
            "&WIDTH=1135&HEIGHT=893",
            "&LAYERS=",
            wms_layer,
            "&STYLES=",
            "&FORMAT=image/png",
            "&TRANSPARENT=TRUE",
        ]
    )
    print(f"requesting WMS {wms_url}")
    wms_response = requests.get(wms_url)
    if wms_response.status_code != 200:
        raise Exception(
            f"WMS response unexpected status {wms_response.status_code} for {wms_url}"
        )
    response_type = wms_response.headers["Content-Type"]
    if response_type != "image/png":
        raise Exception(
            f"WMS response unexpected content-type {response_type} for {wms_url}"
        )
