package main

// handlers module holds all HTTP handlers functions
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/uptrace/bunrouter"
	"golang.org/x/oauth2"
)

// HTTPError represents HTTP error record
type HTTPError struct {
	Method         string `json:"method"`           // HTTP method
	HTTPCode       int    `json:"http_code"`        // HTTP error code
	Code           int    `json:"code"`             // server status code
	Timestamp      string `json:"timestamp"`        // timestamp of the error
	Path           string `json:"path"`             // URL path
	UserAgent      string `json:"user_agent"`       // http user-agent field
	XForwardedHost string `json:"x_forwarded_host"` // http.Request X-Forwarded-Host
	XForwardedFor  string `json:"x_forwarded_for"`  // http.Request X-Forwarded-For
	RemoteAddr     string `json:"remote_addr"`      // http.Request remote address
	Reason         string `json:"reason"`           // error message
}

// HTTPResponse rpresents HTTP JSON response
type HTTPResponse struct {
	Method         string `json:"method"`           // HTTP method
	Path           string `json:"path"`             // URL path
	UserAgent      string `json:"user_agent"`       // http user-agent field
	XForwardedHost string `json:"x_forwarded_host"` // http.Request X-Forwarded-Host
	XForwardedFor  string `json:"x_forwarded_for"`  // http.Request X-Forwarded-For
	RemoteAddr     string `json:"remote_addr"`      // http.Request remote address
	HTTPCode       int    `json:"http_code"`        // HTTP error code
	Code           int    `json:"code"`             // server status code
	Reason         string `json:"reason"`           // error code reason
	Timestamp      string `json:"timestamp"`        // timestamp of the error
	Response       string `json:"response"`         // response message
	Error          string `json:"error"`            // error message
	Data           string `json:"data"`             // HTTP response data
	ElapsedTime    string `json:"elapsed_time"`     // elapsed time of HTTP request
}

// helper function to parse given template and return HTML page
func tmplPage(tmpl string, tmplData TmplRecord) string {
	if tmplData == nil {
		tmplData = make(TmplRecord)
	}
	var templates Templates
	page := templates.Tmpl(tmpl, tmplData)
	//     tdir := fmt.Sprintf("%s/templates", Config.StaticDir)
	//     page := templates.TmplFile(tdir, tmpl, tmplData)
	return page
}

// helper function to generate JSON response
func httpResponse(w http.ResponseWriter, r *http.Request, tmpl TmplRecord) {
	httpCode := tmpl.GetInt("HttpCode")
	tmpl["EndTime"] = time.Now().Unix()
	elapsedTime := tmpl.GetElapsedTime()
	tmpl["ElapsedTime"] = elapsedTime
	// regenerate top part since we may
	tmpl["Top"] = tmplPage("top.tmpl", tmpl)
	top := tmpl.GetString("Top")
	bottom := tmpl.GetString("Bottom")
	tfile := tmpl.GetString("Template")
	if tfile == "" {
		if httpCode == 0 || httpCode == http.StatusOK {
			tfile = "success.tmpl"
		} else {
			tfile = "error.tmpl"
		}
	}
	page := tmplPage(tfile, tmpl)
	if httpCode != 0 {
		w.WriteHeader(httpCode)
	}
	w.Write([]byte(top + page + bottom))
}

// helper function to provide standard HTTP error reply
func httpError(w http.ResponseWriter, r *http.Request, tmpl TmplRecord, code int, err error, httpCode int) {
	tmpl["Code"] = code
	tmpl["Error"] = err
	tmpl["HttpCode"] = httpCode
	tmpl["Content"] = err.Error()
	tmpl["Template"] = "error.tmpl"
	httpResponse(w, r, tmpl)
}

// helper function to make initial template struct
func makeTmpl(title string) TmplRecord {
	tmpl := make(TmplRecord)
	tmpl["Title"] = title
	tmpl["User"] = ""
	tmpl["Base"] = Config.Base
	tmpl["ServerInfo"] = info()
	tmpl["Top"] = tmplPage("top.tmpl", tmpl)
	tmpl["Bottom"] = tmplPage("bottom.tmpl", tmpl)
	tmpl["StartTime"] = time.Now().Unix()
	return tmpl
}

// gologinHandler provides wrapper for gologin handlers
// it gets HTTP request referrer and adds this information to oauth2 RedirectURL
func gologinHandler(config *oauth2.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get HTTP request referrer
		referer := r.Referer()
		if referer != "" && strings.Contains(referer, "redirect=") {
			// modify oauth config RedirectURL with our referrer value
			// if it does not contain redirect part
			if !strings.Contains(config.RedirectURL, "redirect=") {
				arr := strings.Split(referer, "redirect=")
				api := arr[1]
				config.RedirectURL = fmt.Sprintf("%s?redirect=%s", config.RedirectURL, api)
			}
		}
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// FaviconHandler handles favicon icon
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("%s/images/favicon.ico", Config.StaticDir))
}

// helper function to check user's authorization
func checkAuthz(tmpl TmplRecord, w http.ResponseWriter, r *http.Request) error {

	// check if we use dev mode and then do not check for user session
	if Config.DevMode {
		log.Println("WARNING: server use development mode, the checkAuthz is off")
		tmpl["User"] = "dev-user"
		tmpl["Token"] = "dev-token"
		tmpl["Provider"] = "dev-provider"
		return nil
	}

	// get our session cookies
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		return err
	}

	// extract user context from OAuth
	user, ok := session.GetOk(sessionUserName)
	if !ok {
		return errors.New("web session does not present user name")
	}
	token, ok := session.GetOk(sessionToken)
	if !ok {
		return errors.New("web session does not present access token")
	}
	provider, ok := session.GetOk(sessionProvider)
	if !ok {
		return errors.New("web session does not present access token")
	}
	tmpl["User"] = user
	tmpl["Token"] = token
	tmpl["Provider"] = provider
	return nil
}

// helper function to get user name from web session
func getUser(r *http.Request) (string, error) {
	// check if we use dev mode and then do not check for user session
	if Config.DevMode {
		log.Println("WARNING: server use development mode, the checkAuthz is off")
		user := "dev-user"
		return user, nil
	}

	// get our session cookies
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		return "", err
	}

	// extract user context from OAuth
	user, ok := session.GetOk(sessionUserName)
	if !ok {
		return "", errors.New("web session does not present user name")
	}
	return fmt.Sprintf("%v", user), nil
}

// NotebookHandler handles notebook page
func NotebookHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP notebook")
	var userName string
	var err error

	// user HTTP call should present either valid token or it will be
	// redirected to /login end-point
	if err := checkAuthz(tmpl, w, r); err != nil {
		rpath := fmt.Sprintf("%s/login?redirect=%s", Config.Base, r.URL.Path)
		// get our session cookies
		session, err := sessionStore.Get(r, sessionName)
		if err != nil {
			log.Printf("NotebookHandler, session %s redirect to %s due to error %v", sessionName, rpath, err)
			http.Redirect(w, r, rpath, http.StatusTemporaryRedirect)
			return
		}
		// check if ser has been authenticated with any OAuth providers
		user, ok := session.GetOk(sessionUserName)
		if !ok {
			log.Printf("NotebookHandler, unable to identify username due to error %v", err)
			http.Redirect(w, r, rpath, http.StatusTemporaryRedirect)
			return
		}
		userID, _ := session.GetOk(sessionUserID)
		provider, _ := session.GetOk(sessionProvider)
		userName = user.(string)
		tmpl["User"] = userName
		tmpl["UserID"] = userID.(string)
		tmpl["Provider"] = provider.(string)
	} else {
		userName, err = getUser(r)
		if err != nil {
			tmpl["Error"] = err
			tmpl["HttpCode"] = http.StatusBadRequest
			httpResponse(w, r, tmpl)
			return
		}
	}

	// we need to check if given notebook exists and if not we should create it
	notebook := Notebook{
		Host:     Config.JupyterHost,
		Token:    Config.JupyterToken,
		Root:     Config.JupyterRoot,
		User:     userName,
		FileName: "userprocessor.ipynb",
	}
	if Config.Verbose > 0 {
		log.Printf("Notebook %+v", notebook)
	}
	err = notebook.Create()
	if err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusInternalServerError
		httpResponse(w, r, tmpl)
		return
	}
	tmpl["JupyterToken"] = Config.JupyterToken
	tmpl["JupyterHost"] = Config.JupyterHost
	tmpl["Notebook"] = fmt.Sprintf("/users/%s/%s", userName, notebook.FileName)
	tmpl["Workflows"] = chapWorkflows.getWorkflows()
	tmpl["Readers"] = []string{"YamlReader", "NexusReader", "CSVReader"}
	tmpl["Writers"] = []string{"YamlWriter", "NexusWriter", "CSVWriter"}
	tmpl["Processors"] = []string{
		"AsyncProcessor",
		"IntegrationProcessor",
		"IntegrateMapProcessor",
		"MapProcessor",
		"NexusToNumpyProcessor",
		"NexusToXarrayProcessor",
		"PrintProcessor",
		"StrainAnalysisProcessor",
		"XarrayToNexusProcessor",
		"XarrayToNumpyProcessor",
	}
	tmpl["Template"] = "notebook.tmpl"
	httpResponse(w, r, tmpl)
}

// ChapDocHandler handles individual workflow page/API
func ChapDocHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP documentation")
	params := bunrouter.ParamsFromContext(r.Context())
	topic := params.ByName("topic")
	fname := fmt.Sprintf("%s/%s.md", Config.DocDir, topic)
	content, err := mdToHTML(fname)
	if err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusBadRequest
		httpResponse(w, r, tmpl)
		return
	}
	w.Write([]byte(content))
}

// ChapWorkflowHandler handles individual workflow page/API
func ChapWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP workflow")
	// get user name from web session
	user, err := getUser(r)
	if err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusBadRequest
		httpResponse(w, r, tmpl)
		return
	}
	params := bunrouter.ParamsFromContext(r.Context())
	workflow := params.ByName("workflow")
	module := "userprocessor" // it is irrelevant in this case
	config := genWorkflowConfig(user, module, workflow)
	tmpl["Config"] = config
	tmpl["Workflow"] = workflow
	tmpl["Template"] = "workflow_config.tmpl"
	httpResponse(w, r, tmpl)
}

// ChapConfigHandler handles individual workflow configuration
func ChapConfigHandler(w http.ResponseWriter, r *http.Request) {
	// get user name from web session
	user, err := getUser(r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	params := bunrouter.ParamsFromContext(r.Context())
	workflow := params.ByName("workflow")
	if r.Method == "POST" {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		if Config.Verbose > 0 {
			log.Printf("### POST record=%+v error=%v", string(data), err)
		}
		// write provided config content back to user's area into chap.yaml file
		fname := fmt.Sprintf("%s/users/%s/%s/chap.yaml", Config.UserDir, user, workflow)
		file, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.Write(data)
		w.Write([]byte("success"))
	}
	module := "userprocessor" // it is irrelevant in this case
	config := genWorkflowConfig(user, module, workflow)
	w.Write([]byte(config))
}

// ChapRunHandler handles CHAP run page
func ChapRunHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP pipeline")

	// get user name from web session
	user, err := getUser(r)
	if err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusBadRequest
		httpResponse(w, r, tmpl)
		return
	}
	tmpl["User"] = user
	tmpl["Title"] = fmt.Sprintf("CHAP pipeline (%s)", user)

	// prepare notebook
	notebook := Notebook{
		Host:     Config.JupyterHost,
		Token:    Config.JupyterToken,
		Root:     Config.JupyterRoot,
		User:     user,
		FileName: "userprocessor.ipynb"}
	tmpl["Notebook"] = notebook.FileName
	// capture notebook content
	rec, err := notebook.Capture()
	if err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusBadRequest
		httpResponse(w, r, tmpl)
		return
	}
	var lines []string
	for _, cell := range rec.Content.Cells {
		lines = append(lines, cell.Source)
	}
	if Config.Verbose > 0 {
		log.Printf("### CHAP %+v, error %v", rec, err)
	}

	// get reader, writer parameters
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusBadRequest
		httpResponse(w, r, tmpl)
		return
	}
	var workflow string
	if values, ok := params["chapworkflow"]; ok {
		workflow = values[0]
	}

	content := fmt.Sprintf("\n<h1>Workflow: %s</h1><br/>\n", workflow)
	content += fmt.Sprintf("\n<b>Output area: &nbsp;&rArr;&nbsp; <span class=\"blue\"><a href=\"%s/users/%s/%s\">workflow files</a></span></b>", Config.Base, user, workflow)
	content += fmt.Sprintf("\n<b> <span class=\"blue\"><a href=\"%s/users/%s/%s/chap.log\">chap.log</a></span></b>", Config.Base, user, workflow)
	content += "\n<h2>Config:</h2><br/>\n"

	// TODO: think about how dynamically pass module and processor
	// possibly via params to avoid hardcoding
	module := "userprocessor"
	processor := "UserProcessor"
	// generate user code
	genUserCode(user, module, processor, lines)

	// generate user config
	var config string
	if Config.Verbose > 0 {
		log.Printf("### GENERATE NOTEBOOK workflow=%s", workflow)
	}
	if workflow != "" {
		config = genWorkflowConfig(user, module, workflow)
	} else {
		config = genChapConfig(user, module, "yaml", "yaml")
	}
	content += fmt.Sprintf("\n<pre>\n%s\n</pre><br/>\n", config)

	// run CHAP pipeline
	out, err := runCHAP(user, config, workflow)
	if err != nil {
		tmpl["Error"] = err
		tmpl["Template"] = "error.tmpl"
	} else {
		tmpl["Template"] = "success.tmpl"
	}

	// prepare web response
	userInput := strings.Trim(strings.Join(lines, "\n"), " ")
	content += fmt.Sprintf("\n<h2>Input:</h2>\n<pre>\n%s\n</pre><br/>\n", userInput)
	content += fmt.Sprintf("\n<h2>Log:</h2>\n<pre>\n%s\n</pre><br/>\n", out)
	if err != nil {
		content += fmt.Sprintf("\n<h2>Error:</h2>\n<pre>\n%v\n</pre><br/>\n", err)
	}
	if Config.Verbose > 0 {
		log.Println("### CHAP content\n", content)
	}
	tmpl["Content"] = template.HTML(content)

	// check if profile output is present
	// TODO: to handle multiple users we should probably provide
	// ability to specify user's profile file name
	// so far we'll use default one
	fname := "profile.dat"
	if _, err := os.Stat(fname); errors.Is(err, os.ErrExist) {
		file, err := os.Open(fname)
		if err != nil {
			log.Println("ERROR:", err)
		}
		defer file.Close()
		if data, err := io.ReadAll(file); err == nil {
			content += fmt.Sprintf("\n<h2>Profile info:</h2>\n<pre>\n%v</pre><br/>\n", string(data))
		} else {
			log.Println("ERROR:", err)
		}
	}

	httpResponse(w, r, tmpl)
}

// ChapProfileHandler handles CHAP profile page
func ChapProfileHandler(w http.ResponseWriter, r *http.Request) {
	// set profiler HTTP header
	r.Header.Set("profile", "true")
	ChapRunHandler(w, r)
}

// ChapCommitHandler handles publishing page
func ChapCommitHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP commit")
	var userName, msg string
	var err error
	userName, err = getUser(r)
	if err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusBadRequest
		httpResponse(w, r, tmpl)
		return
	}
	cmd := fmt.Sprintf("%s/commit.sh", Config.ScriptsDir)
	notebook := filepath.Join(Config.UserDir, userName)
	log.Printf("shell # %s %s %s", cmd, notebook, Config.UserRepo)
	out, err := exec.Command(cmd, notebook, Config.UserRepo).CombinedOutput()
	log.Println("shell output # ", string(out), err)
	content := fmt.Sprintf("\n<b>Commit status: </b>")
	status := "SUCCESS"
	if err != nil {
		tmpl["Error"] = err
		tmpl["Template"] = "error.tmpl"
		status = "ERROR"
		msg = fmt.Sprintf("Fail to commit user codebase to %s, error %v\n<pre>%s</pre>", Config.UserRepo, err, string(out))
		msg += fmt.Sprintf("\nPlease open ticket at <a href=\"https://github.com/CHESSComputing/CHASaaS/issues\">CHESSComputing/CHASaaS</a> repository")
	} else {
		tmpl["Template"] = "success.tmpl"
		button := fmt.Sprintf("<a href=\"%s/chap/publish\" class=\"button button-small button-round\">Publish</a>", Config.Base)
		msg = fmt.Sprintf("If you reade you may %s your code", button)
	}
	content += fmt.Sprintf("<b>%s</b>\n\n<pre>\n%s\n</pre><br/>\n", status, msg)
	tmpl["Content"] = template.HTML(content)
	httpResponse(w, r, tmpl)
}

// ChapPublishHandler handles publishing page
func ChapPublishHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP publish")
	var userName, msg string
	var err error
	userName, err = getUser(r)
	if err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusBadRequest
		httpResponse(w, r, tmpl)
		return
	}
	userTag := "0" // force publish.sh script to create new tag
	token := getToken()
	releaseNotes := fmt.Sprintf("CHAPBook release %s by %s", userTag, userName)
	cmd := fmt.Sprintf("%s/publish.sh", Config.ScriptsDir)
	log.Printf("shell# %s %s %s %s \"%s\"", cmd, Config.UserRepo, token, userTag, releaseNotes)
	out, err := exec.Command(cmd, Config.UserRepo, token, userTag, releaseNotes).CombinedOutput()
	log.Println("shell output # ", string(out), err)
	content := fmt.Sprintf("\n<b>Publication status: </b>")
	status := "SUCCESS"
	if err != nil {
		tmpl["Error"] = err
		tmpl["Template"] = "error.tmpl"
		msg = fmt.Sprintf("release %s fail to be publised with error %v\n<pre>%v</pre>", userTag, err, string(out))
		msg += fmt.Sprintf("\nPlease open ticket at <a href=\"https://github.com/CHAPUsers/CHAPBook/issues\">CHAPUsers/CHAPBook</a> repository")
		status = "ERROR"
	} else {
		tmpl["Template"] = "success.tmpl"
		msg = fmt.Sprintf("release sucessfully published, DOI: %s", getDOI())
	}
	content += fmt.Sprintf("<b>%s</b>\n<pre>\n%s\n</pre><br/>\n", status, msg)
	tmpl["Content"] = template.HTML(content)
	httpResponse(w, r, tmpl)
}

// WorkflowsHandler handles workflow page
func WorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP workflows")
	// TODO: get list of workflows from user repository and
	// present them on a web
	workflows := getChapWorkflows()
	tmpl["Workflows"] = workflows
	tmpl["Template"] = "workflows.tmpl"
	httpResponse(w, r, tmpl)
}

// LoginHandler handles login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP login")
	tmpl["GithubLogin"] = fmt.Sprintf("%s/github/login", Config.Base)
	tmpl["Template"] = "login.tmpl"
	httpResponse(w, r, tmpl)
}

// AccessHandler handles access page
func AccessHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP access")

	// user HTTP call should present either valid token or it will be
	// redirected to /login end-point
	if err := checkAuthz(tmpl, w, r); err != nil {
		tmpl["Error"] = err
		tmpl["HttpCode"] = http.StatusBadRequest
		httpResponse(w, r, tmpl)
		return
	}
	user := tmpl.GetString("User")
	token := tmpl.GetString("Token")
	if Config.Verbose > 0 {
		log.Printf("AccessHandler: user %s token %s", user, token)
	}

	// HTTP response with user info
	content := fmt.Sprintf("User %s, access token: %s", user, token)
	tmpl["Content"] = template.HTML(content)
	tmpl["Template"] = "success.tmpl"
	httpResponse(w, r, tmpl)
}

// DocsHandler handles documentation page
func DocsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAP documentation")
	fname := fmt.Sprintf("%s/md/docs.md", Config.StaticDir)
	content, err := mdToHTML(fname)
	if err != nil {
		httpError(w, r, tmpl, FileIOError, err, http.StatusInternalServerError)
		return
	}
	tmpl["Content"] = template.HTML(content)
	tmpl["Template"] = "docs.tmpl"
	httpResponse(w, r, tmpl)
}

// IndexHandler handles CHAPBook index (main) page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := makeTmpl("CHAPBook")
	tmpl["Base"] = Config.Base
	tmpl["Token"] = ""
	top := tmplPage("top.tmpl", tmpl)
	bottom := tmplPage("bottom.tmpl", tmpl)
	page := tmplPage("index.tmpl", tmpl)
	w.Write([]byte(top + page + bottom))
}
