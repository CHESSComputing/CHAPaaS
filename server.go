package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/uptrace/bunrouter"

	gologin "github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

// content is our static web server content.
//go:embed static
var StaticFs embed.FS

// The OAuth parts are based on
// https://github.com/dghubble/gologin
// package where we explid github authentication, see
// https://github.com/dghubble/gologin/blob/main/examples/github

// helper function to get base path
func basePath(s string) string {
	if Config.Base != "" {
		if strings.HasPrefix(s, "/") {
			s = strings.Replace(s, "/", "", 1)
		}
		if strings.HasPrefix(Config.Base, "/") {
			return fmt.Sprintf("%s/%s", Config.Base, s)
		}
		return fmt.Sprintf("/%s/%s", Config.Base, s)
	}
	return s
}

// bunrouter implementation of the compatible (with net/http) router handlers
func bunRouter() *bunrouter.CompatRouter {
	router := bunrouter.New(
		bunrouter.Use(bunrouterLoggingMiddleware),
		bunrouter.Use(bunrouterLimitMiddleware),
	).Compat()
	base := Config.Base
	router.GET(base+"/", IndexHandler)
	router.GET(base+"/favicon.ico", FaviconHandler)

	// server routes
	router.GET(base+"/docs", DocsHandler)
	router.GET(base+"/login", LoginHandler)
	router.GET(base+"/access", AccessHandler)
	router.GET(base+"/publish", PublishHandler)
	router.GET(base+"/notebook", NotebookHandler)
	router.GET(base+"/workflows", WorkflowsHandler)

	// chap routes
	router.GET(base+"/chap/run", ChapRunHandler)
	router.GET(base+"/chap/profile", ChapProfileHandler)

	// auth end-points
	// github OAuth routes
	var arec OAuthRecord
	var err error
	arec, err = Config.Credentials("github")
	if err != nil {
		log.Println("WARNING:", err)
	}
	config := &oauth2.Config{
		ClientID:     arec.ClientID,
		ClientSecret: arec.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/github/callback", Config.RedirectURL),
		// RedirectURL:  fmt.Sprintf("http://localhost:%d%s/github/callback", Config.Port, Config.Base),
		Endpoint: githubOAuth2.Endpoint,
	}
	stateConfig := gologin.DebugOnlyCookieConfig
	githubLogin := gologinHandler(config, github.StateHandler(stateConfig, github.LoginHandler(config, nil)))
	githubCallback := gologinHandler(config, github.StateHandler(
		stateConfig,
		github.CallbackHandler(config, issueSession("github"), nil),
	))
	router.Router.GET(base+"/github/login", bunrouter.HTTPHandler(githubLogin))
	router.Router.GET(base+"/github/callback", bunrouter.HTTPHandler(githubCallback))

	// static handlers
	for _, dir := range []string{"js", "css", "images"} {
		filesFS, err := fs.Sub(StaticFs, "static/"+dir)
		if err != nil {
			panic(err)
		}
		m := fmt.Sprintf("%s/%s", Config.Base, dir)
		fileServer := http.FileServer(http.FS(filesFS))
		hdlr := http.StripPrefix(m, fileServer)
		router.Router.GET(m+"/*path", bunrouter.HTTPHandler(hdlr))
	}

	return router
}

// Server implements MLaaS server
func Server() {

	// initialize server middleware
	initLimiter(Config.LimiterPeriod)

	// setup server router
	router := bunRouter()

	// start HTTPs server
	if len(Config.DomainNames) > 0 {
		server := LetsEncryptServer(Config.DomainNames...)
		log.Println("Start HTTPs server with LetsEncrypt", Config.DomainNames)
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else if Config.ServerCrt != "" && Config.ServerKey != "" {
		tlsConfig := &tls.Config{
			RootCAs: RootCAs(),
		}
		server := &http.Server{
			Addr:      ":https",
			TLSConfig: tlsConfig,
			Handler:   router,
		}
		log.Printf("Start HTTPs server with %s and %s on :%d", Config.ServerCrt, Config.ServerKey, Config.Port)
		log.Fatal(server.ListenAndServeTLS(Config.ServerCrt, Config.ServerKey))
	} else {
		log.Printf("Start HTTP server on :%d", Config.Port)
		http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), router)
	}
}
