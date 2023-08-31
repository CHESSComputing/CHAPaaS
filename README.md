# CHAPaaS
CHAP as a Service (a.k.a CHAPBook) provides a web based service (interface)
for rapid development of CHESS analysis workflows within
[CHESS Analysis Pipeline (CHAP)](https://github.com/CHESSComputing/ChessAnalysisPipeline)
framework.

[![DOI](https://zenodo.org/badge/642907044.svg)](https://zenodo.org/badge/latestdoi/642907044)

### Requirements
CHAPaaS is written in GoLang and can be build using standard GoLang
build procedure. Please use provided `Makefile` and build it using
```
make
```
But in order to run CHAPaaS we need the following:
- [jupyter](https://jupyter.org/install)
notebook where CHAP code will be executed
- [miniconda](https://conda.io/en/latest/miniconda-install.html) or
[conda](https://conda.io/projects/conda/en/stable/)
- corresponding python environments for supported workflows
  - [tomo/conda.yml](https://github.com/CHESSComputing/CHAPBookWorkflows/blob/main/tomo/conda.yml)
  - [saxswaxs/conda.yml](https://github.com/CHESSComputing/CHAPBookWorkflows/blob/main/saxswaxs/conda.yml)
  To install given environment please do the following:
```
conda env create -f <workflow>/conda.yml
```

### Running CHAPaaS
To run the service please follow these steps:
1. Setup jupyter notebook server elsewhere, it can be local installation or
   remote one
   - configure your notebook to allow access from CHAPaaS. Here is an example
   of local jupyter setup:
```
c.NotebookApp.tornado_settings = {
    "headers": {
        "Content-Security-Policy": "frame-ancestors 'self' http://localhost:8182"
    }
}
c.NotebookApp.token = "jupyter-token"
c.NotebookApp.open_browser = False
c.NotebookApp.port = 18889
c.NotebookApp.ip = '*'
```
Feel free to provide your own token value, and adjust ports as necessary.
In this example port `8182` refers to CHAPaaS port, while `18889`
is notebook server port.
2. Install either miniconda or coda tool
   - install all supported workflow environments, e.g.
```
git clone git@github.com:CHESSComputing/CHAPBookWorkflows.git
conda env create -f CHAPBookWorkflows/tomo/conda.yml
conda env create -f CHAPBookWorkflows/saxswaxs/conda.yml
...
```
3. Install [ChessAnalysisPipeline](https://github.com/CHESSComputing/ChessAnalysisPipeline)
and make it accessible in your PATH.
4. Create proper configuration file, e.g. `config.json`
```
{
    "base": "",
    "verbose": 1,
    "chap": "/path/CHAPaaS/scripts/chap.sh",
    "chap_dir": "/path/ChessAnalysisPipeline",
    "user_dir": "/path/CHAPaaS/chap/users",
    "user_repo": "/tmp/CHAPBook",
    "scripts_dir": "/path/CHAPaaS/scripts",
    "redirect_url": "http://localhost:8182",
    "oauth": [
        {
            "provider": "github",
            "client_id" : "123",
            "client_secret": "xyz"
        }
    ],
    "jupyter_root": "/tmp/jupyter",
    "jupyter_token": "jupyter-token",
    "jupyter_host": "http://localhost:18889",
    "workflows_root": "/path/CHAPaaS/examples",
    "github_token": "/path/CHAPaaS/token",
    "doi": "666131920",
    "development_mode": true,
    "port": 8182
}
```
5. Start the service
```
./chapaas -config config.json
```
