package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
	tmpl["Lines"] = lines
	tmpl["UserProcessor"] = processor
	tfile := "processor.tmpl"
	var templates Templates
	content := templates.TextTmpl(tfile, tmpl)
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
		log.Println("genChapConfog:\n", config)
	}
	return config
}
func getChapWorkflows() []Workflow {
	var workflows []Workflow
	w := Workflow{
		UserName:    "user",
		Name:        "SAXSWAX",
		Type:        "saxswaxs workflow type",
		Group:       "beamline-X",
		Version:     "v0.0.1",
		Description: "bla-bla-bla",
		Reference:   "http://some.site/tomo",
	}
	workflows = append(workflows, w)
	w = Workflow{
		UserName:    "user",
		Name:        "TOMO",
		Type:        "tomo workflow type",
		Group:       "beamlineA",
		Version:     "v1.2.3",
		Description: "bla-bla-bla",
		Reference:   "http://some.site/saxswax",
	}
	workflows = append(workflows, w)
	return workflows
}
