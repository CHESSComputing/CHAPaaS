package main

// config module
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// OAuthRecord defines OAuth provider's credentials
type OAuthRecord struct {
	Provider     string `json:"provider"`      // name of the provider
	ClientID     string `json:"client_id"`     // client id
	ClientSecret string `json:"client_secret"` // client secret
}

// Configuration stores server configuration parameters
type Configuration struct {
	// web server parts
	Base        string `json:"base"`             // base URL
	LogFile     string `json:"log_file"`         // server log file
	Port        int    `json:"port"`             // server port number
	Verbose     int    `json:"verbose"`          // verbose output
	StaticDir   string `json:"static_dir"`       // speficy static dir location
	RedirectURL string `json:"redirect_url"`     // redirect URL for OAuth provider
	DevMode     bool   `json:"development_mode"` // turn off authz in development mode

	// OAuth parts
	OAuth []OAuthRecord `json:"oauth"` // oauth configurations

	// proxy parts
	XForwardedHost      string `json:"X-Forwarded-Host"`       // X-Forwarded-Host field of HTTP request
	XContentTypeOptions string `json:"X-Content-Type-Options"` // X-Content-Type-Options option

	// server parts
	RootCAs         string   `json:"rootCAs"`      // server Root CAs path
	ServerCrt       string   `json:"server_cert"`  // server certificate
	ServerKey       string   `json:"server_key"`   // server certificate
	DomainNames     []string `json:"domain_names"` // LetsEncrypt domain names
	LimiterPeriod   string   `json:"rate"`         // limiter rate value
	FoxdenPublicKey string   `json:"foxden_public_key"`

	// storage parts
	StorageDir string `json:"storage_dir"` // storage directory

	// CHAP parts
	ChapDir       string `json:"chap_dir"`       // CHAP install area
	UserDir       string `json:"user_dir"`       // user directory
	DocDir        string `json:"doc_dir"`        // CHAP doc directory
	UserRepo      string `json:"user_repo"`      // user repo to use, e.g. CHAPUsers
	ScriptsDir    string `json:"scripts_dir"`    // scripts dir area
	JupyterToken  string `json:"jupyter_token"`  // jupyter token
	JupyterHost   string `json:"jupyter_host"`   // jupyter host:port
	JupyterRoot   string `json:"jupyter_root"`   // jupyter root directory
	WorkflowsRoot string `json:"workflows_root"` // workflows directory
	GithubToken   string `json:"github_token"`   // github token to use for publication
	DOI           string `json:"doi"`            // CHAPUsers/CHAPBook DOI reference
}

// Credentials returns provider OAuth credential record
func (c Configuration) Credentials(provider string) (OAuthRecord, error) {
	for _, rec := range c.OAuth {
		if rec.Provider == provider {
			return rec, nil
		}
	}
	msg := fmt.Sprintf("No OAuth provider %s is found", provider)
	return OAuthRecord{}, errors.New(msg)
}

// Config variable represents configuration object
var Config Configuration

// helper function to parse server configuration file
func parseConfig(configFile string) error {
	data, err := os.ReadFile(filepath.Clean(configFile))
	if err != nil {
		log.Println("Unable to read", err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println("Unable to parse", err)
		return err
	}

	// default values
	if Config.Port == 0 {
		Config.Port = 8181
	}
	if Config.LimiterPeriod == "" {
		Config.LimiterPeriod = "100-S"
	}
	if Config.StaticDir == "" {
		cdir, err := os.Getwd()
		if err == nil {
			Config.StaticDir = fmt.Sprintf("%s/static", cdir)
		} else {
			Config.StaticDir = "static"
		}
	}
	if Config.StorageDir == "" {
		Config.StorageDir = "/tmp"
	}
	if Config.JupyterHost == "" {
		log.Fatal("Empty JupyterHost, please adjust your configuration")
	}
	if Config.UserDir == "" {
		log.Fatal("Empty UserDir, please adjust your configuration")
	}
	if Config.DocDir == "" {
		log.Fatal("Empty DocDir, please adjust your configuration")
	}
	if Config.UserRepo == "" {
		log.Fatal("Empty UserRepo, please adjust your configuration")
	}
	if Config.ScriptsDir == "" {
		Config.ScriptsDir = "scripts"
	}
	if Config.JupyterRoot == "" {
		log.Fatal("Empty JupyterRoot, please adjust your configuration")
	}
	if Config.ChapDir == "" {
		log.Fatal("Empty ChapDir, please adjust your configuration")
	}
	if Config.RedirectURL == "" {
		if host, err := os.Hostname(); err == nil {
			Config.RedirectURL = fmt.Sprintf("http://%s:%d%s/github/callback", host, Config.Port, Config.Base)
		} else {
			Config.RedirectURL = fmt.Sprintf("http://localhost:%d%s/github/callback", Config.Port, Config.Base)
		}
		log.Println("RedirectURL", Config.RedirectURL)
	}
	return nil
}
