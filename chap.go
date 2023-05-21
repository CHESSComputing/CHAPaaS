package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// helper function to run CHAP pipeline
func runCHAP(user, config string) ([]byte, error) {
	fname := fmt.Sprintf("/tmp/chap-%s.yaml", user)
	err := os.Remove(fname)
	if err != nil {
		log.Println("runCHAP os.remove", err)
	}
	file, err := os.Create(fname)
	if err != nil {
		log.Println("runCHAP os.create", err)
	}
	defer file.Close()
	file.Write([]byte(config))

	// run CHAP pipeline
	cmd := fmt.Sprintf("%s", Config.CHAP)
	log.Println("### runCHAP:", cmd, fname)
	out, err := exec.Command(cmd, fname).Output()
	return out, err
}

// helper function to generate user code
func genUserCode(user string, lines []string) {
	tmpl := make(TmplRecord)
	tmpl["Lines"] = lines
	tfile := "processor.tmpl"
	var templates Templates
	content := templates.TextTmpl(tfile, tmpl)
	tdir := fmt.Sprintf("/tmp/%s", user)
	err := os.MkdirAll(fmt.Sprintf("/tmp/%s", user), 0755)
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
}

// helper function to generate CHAP config based on user workflow
func genChapConfig(user string) string {
	config := `
pipeline:
  - UserProcessor: {}
  - common.PrintProcessor: {}
`
	return config
}
