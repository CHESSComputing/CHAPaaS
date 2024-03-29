package main

// utils module
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"archive/tar"
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	mhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/crypto/acme/autocert"
)

// RootCAs returns cert pool of our root CAs
func RootCAs() *x509.CertPool {
	log.Println("Load RootCAs from", Config.RootCAs)
	rootCAs := x509.NewCertPool()
	files, err := ioutil.ReadDir(Config.RootCAs)
	if err != nil {
		log.Printf("Unable to list files in '%s', error: %v\n", Config.RootCAs, err)
		return rootCAs
	}
	for _, finfo := range files {
		fname := fmt.Sprintf("%s/%s", Config.RootCAs, finfo.Name())
		caCert, err := os.ReadFile(filepath.Clean(fname))
		if err != nil {
			if Config.Verbose > 1 {
				log.Printf("Unable to read %s\n", fname)
			}
		}
		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			if Config.Verbose > 1 {
				log.Printf("invalid PEM format while importing trust-chain: %q", fname)
			}
		}
		if Config.Verbose > 1 {
			log.Println("Load CA file", fname)
		}
	}
	return rootCAs
}

// LetsEncryptServer provides HTTPs server with Let's encrypt for
// given domain names (hosts)
func LetsEncryptServer(hosts ...string) *http.Server {
	// setup LetsEncrypt cert manager
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Cache:      autocert.DirCache("certs"),
	}

	tlsConfig := &tls.Config{
		// Set InsecureSkipVerify to skip the default validation we are
		// replacing. This will not disable VerifyPeerCertificate.
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequestClientCert,
		RootCAs:            RootCAs(),
		GetCertificate:     certManager.GetCertificate,
	}

	// start HTTP server with our rootCAs and LetsEncrypt certificates
	server := &http.Server{
		Addr:      ":https",
		TLSConfig: tlsConfig,
	}
	// start cert Manager goroutine
	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	log.Println("Starting LetsEncrypt HTTPs server")
	return server
}

// LogName return proper log name based on Config.LogName and either
// hostname or pod name (used in k8s environment).
func LogName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("unable to get hostname", err)
	}
	if os.Getenv("MY_POD_NAME") != "" {
		hostname = os.Getenv("MY_POD_NAME")
	}
	logName := Config.LogFile + "_%Y%m%d"
	if hostname != "" {
		logName = fmt.Sprintf("%s_%s", Config.LogFile, hostname) + "_%Y%m%d"
	}
	return logName
}

// helper function to parse given markdown file and return HTML content
func mdToHTML(fname string) (string, error) {
	file, err := os.Open(fname)
	if err != nil {
		log.Println("ERROR: unable to open", fname, err)
		return "", err
	}
	defer file.Close()
	var md []byte
	md, err = io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := mhtml.CommonFlags | mhtml.HrefTargetBlank
	opts := mhtml.RendererOptions{Flags: htmlFlags}
	renderer := mhtml.NewRenderer(opts)
	content := markdown.Render(doc, renderer)
	//     return html.EscapeString(string(content)), nil
	return string(content), nil
}

// helper function to return CHAPUsers/CHAPBook github token
func getToken() string {
	fname := Config.GithubToken
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	body, err := io.ReadAll(file)
	if err != nil {
		log.Println("ERROR: unable to read github token file", fname, err)
		return ""
	}
	token := string(body)
	return strings.Replace(token, "\n", "", -1)
}

// get latest DOI badge
func getDOI() string {
	doi := Config.DOI
	return fmt.Sprintf("<a href=\"https://zenodo.org/badge/latestdoi/%s\"><img src=\"https://zenodo.org/badge/%s.svg\" alt=\"DOI\"></a>", doi, doi)
}

// Tar function creates tar ball from the source and target
// https://golangdocs.com/tar-gzip-in-golang
func Tar(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

// Gzip function gzip given source and produce output at target destination
func Gzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}
