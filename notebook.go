package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

/*
 * The JupyterRoot defines top level directory where we run Jupyter app
 * Within this area we must create /users where we'll store each individual
 * user notebooks. Therefore, in a code we use /api/contents/users API
 * which includes /users path which should exist in JupyterRoot
 * Create() api of Notebook struct will properly creats /users area within JupyterRoot area
 */

// Notebook represents jupyter notebook object
type Notebook struct {
	Host     string // jupyter hostname
	Token    string // jupyter server token
	Root     string // jupyter root area
	User     string // jupyter user name
	FileName string // notebook file name
}

// Create creates notebook user area and notebook file
func (n *Notebook) Create() error {
	// ensure that new user's area exists under JupyterRoot
	path := fmt.Sprintf("%s/users/%s", Config.JupyterRoot, n.User)
	if Config.Verbose > 0 {
		log.Println("create notebook", path)
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	// create notebook file
	// https://jupyter-server.readthedocs.io/en/latest/developers/rest-api.html
	var jsonData = []byte(`{"type": "notebook"}`)
	//     rurl := fmt.Sprintf("%s/api/contents/%s/%s", n.Host, n.Root, n.User)
	rurl := fmt.Sprintf("%s/api/contents/users/%s", n.Host, n.User)
	if Config.Verbose > 0 {
		log.Printf("jupyter request HTTP POST %s %+v", rurl, string(jsonData))
	}
	rec, err := notebookCall("POST", rurl, n.Token, bytes.NewBuffer(jsonData))
	if Config.Verbose > 0 {
		log.Printf("jupyter response %+v, error %v", rec, err)
	}
	return err
}

// Capture fetches content of notebook file
func (n *Notebook) Capture() (NotebookRecord, error) {
	//     rurl := fmt.Sprintf("%s/api/contents/%s/%s/%s", n.Host, n.Root, n.User, n.FileName)
	rurl := fmt.Sprintf("%s/api/contents/users/%s/%s", n.Host, n.User, n.FileName)
	rec, err := notebookCall("GET", rurl, n.Token, nil)
	return rec, err
}

// helper function to make API call to jupyter notebook
func notebookCall(method, rurl, token string, body io.Reader) (NotebookRecord, error) {
	var rec NotebookRecord
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, rurl, body)
	if Config.Verbose > 0 {
		log.Printf("Jupyter notebook request %+v, error=%v", req, err)
	}
	if err != nil {
		return rec, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	req.Header.Add("X-Frame-Options", "SAMEORIGIN")
	resp, err := client.Do(req)
	if Config.Verbose > 0 {
		log.Printf("Jupyter notebook response %+v, error=%v", resp, err)
	}
	if err != nil {
		return rec, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&rec)
	return rec, nil
}

/*
{"name": "Untitled.ipynb", "path": "Untitled.ipynb", "last_modified": "2023-05-19T13:55:49.526770Z", "created": "2023-05-19T13:55:49.526770Z", "content": {"cells": [{"cell_type": "code", "execution_count": 1, "id": "4dc583a6", "metadata": {"trusted": false}, "outputs": [], "source": "a=1"}, {"cell_type": "code", "execution_count": 2, "id": "a5ae958c", "metadata": {"trusted": false}, "outputs": [{"name": "stdout", "output_type": "stream", "text": "1\n"}], "source": "print(a)"}, {"cell_type": "code", "execution_count": null, "id": "e7fdf954", "metadata": {"trusted": false}, "outputs": [], "source": ""}], "metadata": {"kernelspec": {"display_name": "Python 3 (ipykernel)", "language": "python", "name": "python3"}, "language_info": {"codemirror_mode": {"name": "ipython", "version": 3}, "file_extension": ".py", "mimetype": "text/x-python", "name": "python", "nbconvert_exporter": "python", "pygments_lexer": "ipython3", "version": "3.11.3"}}, "nbformat": 4, "nbformat_minor": 5}, "format": "json", "mimetype": null, "size": 989, "writable": true, "type": "notebook"}
*/
type NotebookRecord struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	LastModified string `json:"last_modified"`
	Created      string `json:"created"`
	Content      NotebookContent
}
type NotebookContent struct {
	Cells []Cell
}
type Cell struct {
	CellType         string
	ExecutionCounter int
	Id               string
	Source           string
}
