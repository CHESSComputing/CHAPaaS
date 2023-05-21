#!/usr/bin/env python

import sys
import json
import requests

# The token is written on stdout when you start the notebook
notebook_path = '/Untitled.ipynb'
base = 'http://localhost:8888'
token=sys.argv[1]
headers = {'Authorization': 'Token %s' % token}

url = base + '/api/kernels'
response = requests.post(url,headers=headers)
kernel = json.loads(response.text)

# Load the notebook and get the code of each cell
url = base + '/api/contents' + notebook_path
response = requests.get(url,headers=headers)
print("response:", response.text)
file = json.loads(response.text)
code = [ c['source'] for c in file['content']['cells'] if len(c['source'])>0 ]

print("### code", code)
