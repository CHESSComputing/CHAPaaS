package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func initUserDir(user string) error {
	// create and initialize user dir
	err := os.MkdirAll(fmt.Sprintf("%s", Config.UserDir), 0755)
	if err != nil {
		log.Printf("ERROR: unable to create used dir %s, error %v", Config.UserDir, err)
		return err
	}
	fname := fmt.Sprintf("%s/__init__.py", Config.UserDir)
	file, err := os.Create(fname)
	if err != nil {
		log.Printf("ERROR: unable to create %s, error %v", fname, err)
		return err
	}
	defer file.Close()
	file.Write([]byte("# auto-generated file to load user processors\n"))
	return nil
}

/*
// helper function to add user module processor to user dir __init__.py file
func addUserProcessor(user, module, processor string) {
	fname := fmt.Sprintf("%s/__init__.py", Config.UserDir)
	if Config.Verbose > 0 {
		log.Printf("update %s with module=%s and processor=%s", fname, module, processor)
	}
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, user) && strings.Contains(line, processor) {
			found = true
			break
		}
	}
	if !found {
		line := fmt.Sprintf("from %s.%s import %s\n", user, module, processor)
		file.Write([]byte(line))
	}
}
*/

// helper function to run CHAP pipeline
func runCHAP(user, config string) ([]byte, error) {
	var out []byte
	var err error
	fname := fmt.Sprintf("%s/%s/chap.yaml", Config.UserDir, user)
	if Config.Verbose > 0 {
		log.Println("writing config file", fname)
	}
	if _, err := os.Stat(fname); errors.Is(err, os.ErrExist) {
		err = os.Remove(fname)
		if err != nil {
			log.Println("runCHAP os.remove", err)
			return out, err
		}
	}
	file, err := os.Create(fname)
	if err != nil {
		log.Println("runCHAP os.create", err)
	}
	defer file.Close()
	file.Write([]byte(config))
	// run CHAP pipeline
	cmd := fmt.Sprintf("%s/chap.sh", Config.ScriptsDir)
	// user dir in configuration contains /users suffix
	// but for running chap.sh script we should strip it off
	userDir := strings.Replace(Config.UserDir, "/users", "", -1)
	log.Println("### runCHAP:", cmd, fname, Config.ChapDir, userDir)
	out, err = exec.Command(cmd, fname, Config.ChapDir, userDir).Output()
	return out, err
}

// helper function to generate user code
func genUserCode(user, module, processor string, lines []string) {

	// initialize user dir
	err := initUserDir(user)
	if err != nil {
		log.Println("ERROR: gen user code", err)
	}
	tmpl := make(TmplRecord)
	var templates Templates
	var content string
	userCode := strings.Join(lines, "\n")
	if strings.Contains(userCode, "class UserProcessor") {
		content = userCode
	} else {
		var newLines []string
		for _, line := range lines {
			newLine := strings.Replace(line, "\n", "\n        ", -1)
			newLines = append(newLines, newLine)
		}
		tmpl["Lines"] = newLines
		tmpl["UserProcessor"] = processor
		tfile := "processor.tmpl"
		content = templates.TextTmpl(tfile, tmpl)
	}
	tdir := fmt.Sprintf("%s/%s", Config.UserDir, user)
	err = os.MkdirAll(fmt.Sprintf("%s/%s", Config.UserDir, user), 0755)
	if err != nil {
		log.Println("genUserCode os.MkdirAll", err)
	}
	fname := fmt.Sprintf("%s/userprocessor.py", tdir)
	err = os.Remove(fname)
	if err != nil {
		log.Println("genUserCode os.remove", err)
	}
	file, err := os.Create(fname)
	if err != nil {
		log.Println("genUserCode", err)
	}
	defer file.Close()
	file.Write([]byte(content))

	// update user dir __init__.py file
	//addUserProcessor(user, module, processor)
}

// helper function to generate CHAP config based on user workflow
func genWorkflowConfig(user, module, workflow string) string {
	var config string
	for _, w := range chapWorkflows.getWorkflows() {
		if w.Name == workflow {
			fname := filepath.Join(w.Directory, w.Config)
			file, err := os.Open(fname)
			if err != nil {
				log.Printf("ERROR: unable to locate %s", fname)
				break
			}
			defer file.Close()
			if body, err := io.ReadAll(file); err == nil {
				config = string(body)
				break
			} else {
				log.Printf("ERROR: unable to read %s, error %v", fname, err)
			}
		}
	}
	if Config.Verbose > 0 {
		log.Println("workflow config\n", config)
	}
	// TODO: properly handle yaml data by unmarshal it to pipeline
	// struct and add new config entry to it, e.g.
	/*
		// convert config to yaml data
		var pipeline map[string]interface{}
		if err = yaml.Unmarshal(body, &p); err == nil {
		}
	*/
	config += fmt.Sprintf("  - users.%s.%s.UserProcessor: {}\n", user, module)
	return config
}

type Pipeline struct {
}

// helper function to generate CHAP config based on given reader/writer
func genChapConfig(user, module, reader, writer string) string {
	config := "pipeline:\n"
	reader = strings.ToLower(reader)
	writer = strings.ToLower(writer)
	if reader == "yaml" {
		config += "  - common.YAMLReader: {}"
	} else if reader == "nexus" {
		config += "  - common.NexuReader: {}"
	}
	//config += "  - UserProcessor: {}\n  - common.PrintProcessor: {}"
	config += fmt.Sprintf("  - users.%s.%s.UserProcessor: {}\n", user, module)
	config += "  - common.PrintProcessor: {}\n"
	if writer == "yaml" {
		config += "  - common.YAMLWriter: {}"
	} else if writer == "nexus" {
		config += "  - common.NexuWriter: {}"
	}
	/*
	   	config := `
	   pipeline:
	     - UserProcessor: {}
	     - common.PrintProcessor: {}
	   `
	*/
	if Config.Verbose > 0 {
		log.Println("genChapConfig:\n", config)
	}
	return config
}

// helper function to get CHAP workflows
func getChapWorkflows() []Workflow {
	var workflows []Workflow
	if entries, err := os.ReadDir(Config.WorkflowsRoot); err == nil {
		for _, entry := range entries {
			if info, err := entry.Info(); err == nil {
				dname := info.Name()
				// list specific workflow directory
				dir := filepath.Join(Config.WorkflowsRoot, dname)
				if Config.Verbose > 0 {
					log.Println("reading dir", dir)
				}
				if files, err := os.ReadDir(dir); err == nil {
					for _, fentry := range files {
						if finfo, err := fentry.Info(); err == nil {
							fname := filepath.Join(dir, finfo.Name())
							if finfo.Name() == "chap.yaml" {
								if Config.Verbose > 0 {
									log.Println("reading chap spec", fname)
								}
								file, err := os.Open(fname)
								if err != nil {
									log.Printf("ERROR: unable to open file %s, error %v", fname, err)
								}
								defer file.Close()
								if body, err := io.ReadAll(file); err == nil {
									var w Workflow
									err = yaml.Unmarshal(body, &w)
									if err == nil {
										w.Directory = filepath.Join(Config.WorkflowsRoot, entry.Name())
										workflows = append(workflows, w)
									} else {
										log.Printf("ERROR: unable to unmarshal body of file %s, error %v", fname, err)
									}
								} else {
									log.Printf("ERROR: unable to read file %s, error %v", fname, err)
								}
							}
						}
					}
				} else {
					log.Printf("ERROR: unable to read files from %s, error %v", dir, err)
				}
			}
		}
	} else {
		log.Println("ERROR:", err)
	}
	return workflows
}
