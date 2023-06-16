# CHAP documentation
CHAP stans for CHess Analysis Pipeline framework. It is Python based
framework with the following components:
- **Reader** represents an object which performs input read operation
from specific resource
- **Writer** represents an object which performs output write operation
of provided data to specific resource
- **Processor** represents an object which performs specific algorithm
over passed data
  - **Fitter** represents specific algorithm designed to fit the data
  - **Visualizer** represents specfic object to transform given
  data into specific visualized form
  - **Transformer** represents specific algorithm to transform given
  data from one form to another.

### CHAPBook
The architecture of CHAPBook (CHAP as a Service) has the following structure:
![Architecture](/images/CHAPaas_architecture.png)
It provides simple and minimalistic web interface for end-user
who can write their own implementation of given **Processor**.
It hides complexity of CHAP framework from end-user, and capture
the user code (presented in Jupyter notebook), execute it within
CHAP pipeline and present back results to end-user.

Here are the main components:
- CHAP users area, defined via `user_dir` in configuration file, contains
all user based workflows (user processors)
- CHAP install area, defined via `chap_dir` in configuration file, contains
installation of CHAP codebase
- CHAP workflows area, defined via `workflows_root` in configuration file,
  contains CHAP pre-defined workflows, e.g. SAXSWAXS, Tomo, etc.
- jupyter area, defined via `jupyter_root` in configuration file, contains
location of user's notebooks
  - we also have `jupyter_host` and `jupyter_port` to appropriately define
  jupyter hostname and port number

CHAPBook workflow consists of the following:
- each user is authenticated with CHAP via github OAuth
  - therefore, we capture user name and some profile info
- on **Notebook** page we create for a given user a pre-defined CHAP notebook
- when user fill out a notebook logic and hit either **Run** or **Profile**
  button his/her code is captured by CHAP server and
  - CHAP creates a new user processor in CHAP user area
  - it creates correspoding pipeline yaml configuration file
  - it executes user's processor via CHAP pipeline
