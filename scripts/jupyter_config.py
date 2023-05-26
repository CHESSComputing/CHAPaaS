# see https://jupyter-notebook.readthedocs.io/en/stable/public_server.html
c.NotebookApp.tornado_settings = {
    "headers": {
        "Content-Security-Policy": "frame-ancestors 'self' http://localhost:8181"
    }
}
c.NotebookApp.token = "47e67734f6221fec0f18fab5c501c8bef133b14195fdbc08"
c.NotebookApp.open_browser = False
c.NotebookApp.port = 8888

c.JupyterHub.tornado_settings = {
    'headers': {
        'Content-Security-Policy': "frame-ancestors 'self' http://localhost:8181"
    }
}
