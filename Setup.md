To setup CHAPaaS we need the following:
- Download and install [Go](https://go.dev/doc/install) compiler

- compile code as following:
```
make
# if everything is fine it will produce the following executable
ls -al chapaas
```

- setup CHAP releated repositories:
```
# use either git or HTTPs URLs below

# ChessAnalysisPipeline
git clone git@github.com:CHESSComputing/ChessAnalysisPipeline.git

# CHAPBook
git clone git@github.com:CHAPUsers/CHAPBook.git
```

- download and install [jupyter](https://jupyter.org/install)

- obtain github token to access and write to CHAPBook repo

- properly configure jupyter notebook, e.g.
```
# see https://jupyter-notebook.readthedocs.io/en/stable/public_server.html

c.NotebookApp.tornado_settings = {
    "headers": {
        "Content-Security-Policy": "frame-ancestors 'self' http://localhost:8181"
    }
}
c.NotebookApp.token = "1234567890"
c.NotebookApp.open_browser = False
c.NotebookApp.port = 18889
c.NotebookApp.ip = '*'
```
The token ID should be put into CHAPaaS config too.

- start jupyter notebook server as following:
```
jupyter notebook --config /path/jupyter_config.py
```

- create configuration file for CHAPaaS (name it as `config-http.json`):
```
{
    "base": "",
    "verbose": 1,
    "chap": "/path/CHAPaaS/scripts/chap.sh",
    "chap_dir": "/path/ChessAnalysisPipeline/install",
    "user_dir": "/path/CHAPaaS/chap/users",
    "user_repo": "/path/CHAPBook",
    "scripts_dir": "/path/CHAPaaS/scripts",
    "redirect_url": "http://localhost:8181",
    "oauth": [
        {
            "provider": "github",
            "client_id" : "xyxyxyxyxyxy",
            "client_secret": "1234567890"
        }
    ],
    "jupyter_root": "/path/jupyter",
    "jupyter_token": "1234567890",
    "jupyter_host": "http://localhost:18889",
    "workflows_root": "/path/CHAPaaS/workflows",
    "github_token": "/path/CHAPaaS/token",
    "doi": "12345",
    "port": 8181
}
```
Here we used the following:
- jupyter token
- obtain CHAPBook repo token and store it into token file used in `github_token` configuration parameter
- obtain CHAP `workflows` and store it in `workflows_root` configureation
  paratemer
- obtain OAuth credentials from github for CHESSComputing organization
  - visit its
    [Setting](https://github.com/organizations/CHESSComputing/settings/applications)
    application page
  - go to developer settings, choose OAuth apps and setup your
  new OAuth app and obtain your `client_id` and `client_secret`
- finally, you may start CHAPaaS service as following:
```
./chapaas -config config-http.json
```
- on CHESS node you can start it as following:
```
./scripts/run.sh
```
