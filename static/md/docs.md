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

### CHAP as a Service
The architecture of CHAP as a Service has the following structure:
![Architecture](/images/CHAPaas_architecture.png)
It provides simple and minimalistic web interface for end-user
who can write their own implementation of given **Processor**.
It hides complexity of CHAP framework from end-user, and capture
the user code (presented in Jupyter notebook), execute it within
CHAP pipeline and present back results to end-user.
