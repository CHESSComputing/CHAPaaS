#!/usr/bin/env python

import sys
import json
import requests

# The token is written on stdout when you start the notebook
notebook_path = '/Untitled.ipynb'
base = 'http://localhost:18889'
token=sys.argv[1]
headers = {'Authorization': 'Token %s' % token}

url = base + '/api/kernels'
response = requests.post(url,headers=headers)
kernel = json.loads(response.text)
print("kernel", kernel)
# if kernel["id"]:
#     url = base+"/api/kernels/"+kernel["id"]
#     response = requests.get(url,headers=headers)
#     data = json.loads(response.text)
#     print("kernel data", data)

url = base + '/api/contents'
response = requests.get(url,headers=headers)
print("response:", response.text)

# Load the notebook and get the code of each cell
# url = base + '/api/contents' + notebook_path
# response = requests.get(url,headers=headers)
# print("response:", response.text)
# file = json.loads(response.text)
# code = [ c['source'] for c in file['content']['cells'] if len(c['source'])>0 ]
# print("### code", code)

# create new notebook
url = base + '/api/contents/users/vkuznet'
# url = base + '/api/contents'
cells = [
        {"cell_type":"markdown", "source": "### Welcome to CHAP notebook.", 'metadata': {}},
        {"cell_type":"markdown", "source": "CHAP provides you data", 'metadata': {}}
        ]
content = {"cells":cells, 'metadata': {}, 'nbformat': 4, 'nbformat_minor': 4}
name = "notebook.ipynb"
body={"type": "notebook",
      "content": content,
      "name": name,
      "path": name,
      "format": "json",
      "created": "2023-06-07T17:34:22.234793Z",
      "last_modified": "2023-06-07T17:34:22.234793Z",
      'mimetype': None,
      'writable': True}
import pprint
pprint.pprint(body)
headers["Content-Type"] = "application/json"
# response = requests.post(url, headers=headers, )
response = requests.put(url + '/' + name, headers=headers, json=body)
print("response:", response.text, response.status_code)
# file = json.loads(response.text)
