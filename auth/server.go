package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var usersPassword = map[string][]byte{
	"joe":  []byte("$2a$12$aMfFQpGSiPiYkekov7LOsu63pZFaWzmlfm1T8lvG6JFj2Bh4SZPWS"),
	"mary": []byte("$2a$12$l398tX477zeEBP6Se0mAv.ZLR8.LZZehuDgbtw2yoQeMjIyCNCsRW"),
}

func verifyUserPass(username, password string) bool {
	wantPass, hasUser := usersPassword[username]
	if !hasUser {
		return false
	}

	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr != nil {
		return true
	}
	return false
}

func main() {
	addr := flag.String("addr", ":4567", "HTTPS network address")
	certFile := flag.String("certfile", "cert.pem", "certificate pem file")
	keyFile := flag.String("keyfile", "key.pem", "key pem file")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Go serves at ....")
	})

	mux.HandleFunc("/secret/", func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok && verifyUserPass(user, pass) {
			fmt.Fprintf(w, "hwere is joni\n")
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})

	srv := &http.Server{
		Addr:    *addr,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	log.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServeTLS(*certFile, *keyFile)
	log.Fatal(err)
}
