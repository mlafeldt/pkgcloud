#!/usr/bin/env python
"""
 Generate Go map of distros supported by Packagecloud API
 By generating the list once, we save an expensive API call.
 See https://packagecloud.io/docs/api#resource_distributions
"""

import base64
import json
import os
import sys
from pathlib import Path
from urllib.request import Request, urlopen

VAR = "supportedDistros"
if len(sys.argv) > 1:
    VAR = sys.argv[1]
token = os.environ["PACKAGECLOUD_TOKEN"]

request = Request("https://packagecloud.io/api/v1/distributions.json")
b64 = base64.b64encode(bytes(f"{token}:", "ascii"))
request.add_header("Authorization", f"Basic {b64.decode('utf-8')}")
with urlopen(request) as resp:
    data = json.loads(resp.read())

result = {}
for distros in data.values():
    for d in distros:
        for v in d["versions"]:
            k = d["index_name"]
            if "index_name" in v:
                k = "/".join([k, v["index_name"]])
            v = v["id"]
            result[k] = v

print(f"// Generated with {Path( __file__ ).name}")
print("\npackage pkgcloud\n")
print(f"var {VAR} = map[string]int{{")
for k, v in sorted(result.items(), key=lambda x: x[1]):
    print(f'\t"{k}": {v},')
print("}")
