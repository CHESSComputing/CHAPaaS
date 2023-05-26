package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Notebook represents jupyter notebook object
type Notebook struct {
	Host  string
	Token string
}

func (n *Notebook) Capture(fname string) (NotebookRecord, error) {
	var rec NotebookRecord
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	rurl := fmt.Sprintf("%s/api/contents/%s", n.Host, fname)
	req, err := http.NewRequest("GET", rurl, nil)
	if Config.Verbose > 0 {
		log.Printf("Jupyter notebook request %+v, error=%v", req, err)
	}
	if err != nil {
		return rec, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", n.Token))
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
