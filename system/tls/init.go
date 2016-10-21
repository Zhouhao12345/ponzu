package tls

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bosssauce/ponzu/system/db"

	"golang.org/x/crypto/acme/autocert"
)

var m autocert.Manager

// setup attempts to locate or create the cert cache directory and the certs for TLS encryption
func setup() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Couldn't find working directory to locate or save certs.")
	}

	cache := autocert.DirCache(filepath.Join(pwd, "system", "tls", "certs"))
	if _, err := os.Stat(string(cache)); os.IsNotExist(err) {
		err := os.MkdirAll(string(cache), os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatalln("Couldn't create cert directory at", cache)
		}
	}

	host, err := db.Config("domain")
	if err != nil {
		log.Fatalln("No 'domain' field set in Configuration. Please add a domain before attempting to make certificates.")
	}

	email, err := db.Config("admin_email")
	if err != nil {
		log.Fatalln("No 'admin_email' field set in Configuration. Please add an admin email before attempting to make certificates.")
	}

	m = autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       cache,
		HostPolicy:  autocert.HostWhitelist(string(host)),
		RenewBefore: time.Hour * 24 * 30,
		Email:       string(email),
	}

}

// Enable runs the setup for creating or locating certificates and starts the TLS server
func Enable() {
	setup()

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}

	go server.ListenAndServeTLS("", "")
}
