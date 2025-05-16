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

- repo management from service account:
```
# add new ssh key
ssh-keygen -t ed25519 -C "chapaas-deploy" -f ~/.ssh/chapaas_deploy_key

# add new deploy key on repo site
# visit https://github.com/CHESSComputing/CHAPaaS/settings/keys

# adjust .git/config
[core]
        repositoryformatversion = 0
        filemode = true
        bare = false
        logallrefupdates = true
        sshCommand = ssh -i ~/.ssh/chapaas_deploy_key -o IdentitiesOnly=yes
[remote "origin"]
        url = git@github.com:CHESSComputing/CHAPaaS.git
        fetch = +refs/heads/*:refs/remotes/origin/*
[branch "main"]
        remote = origin
        merge = refs/heads/main
[pull]
        rebase = false


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
	Visit github page: https://github.com/settings/tokens?type=beta
	and go to developers settings where we can generate fine-grained access token:
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
