package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	pmrlib "github.com/HeikoCholeva/pmrlib-go"
	auth "github.com/abbot/go-http-auth"
)

var srv http.Server
var users []string

func stopServer() {
	err := srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func startServer(listen string, path string, tls bool, cert string, key string, ba bool) {
	srv.Addr = listen
	log.Printf("Starting server on port %v. TLS enabled: %v, Basic Auth enabled: %v\n", srv.Addr, tls, ba)

	if ba {
		authenticator := auth.NewBasicAuthenticator("localhost", getSecret)
		http.HandleFunc(path, authenticator.Wrap(handleAuth))
		initAuthFile()
	} else {
		http.HandleFunc(path, handler)
	}

	if tls {
		log.Fatal(srv.ListenAndServeTLS(cert, key))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}

func initAuthFile() {
	log.Printf("Using file \"%v\" for http basic auth", cfg.WebServer.BasicAuthFile)

	f, err := os.OpenFile(cfg.WebServer.BasicAuthFile, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening file \"%v\": %v", cfg.WebServer.BasicAuthFile, err)
	}
	defer f.Close()

	buf, err := ioutil.ReadFile(cfg.WebServer.BasicAuthFile)
	if err != nil {
		panic(err)
	}

	userlist := strings.Replace(string(buf), " ", "", -1)
	users = strings.Split(userlist, "\n")
}

func handleAuth(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	r.Close = true
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var scode int

	if strings.ToUpper(r.Method) == "POST" {
		var report pmrlib.Report
		if err := json.Unmarshal(body, &report); err != nil {
			scode = 400
			w.WriteHeader(scode)
			w.Header().Set("Connection", "close")
		} else {
			writeToKafka(string(body))
			scode = 202
			w.WriteHeader(scode)
			w.Header().Set("Connection", "close")
		}
	} else {
		scode = 405
		w.WriteHeader(scode)
		w.Header().Set("Connection", "close")
	}
	log.Println(r.RemoteAddr, "\t", r.Method, "-", scode, "\t", r.RequestURI, "\t", len(body), "\tUser: ", r.Username)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var scode int

	if strings.ToUpper(r.Method) == "POST" {
		var report pmrlib.Report
		if err := json.Unmarshal(body, &report); err != nil {
			scode = 400
			w.WriteHeader(scode)
			w.Header().Set("Connection", "close")
		} else {
			b, err := json.Marshal(report)
			if err != nil {
				scode = 400
			} else {
				scode = 202
				writeToKafka(string(b))
			}
			w.WriteHeader(scode)
			w.Header().Set("Connection", "close")
		}
	} else {
		scode = 405
		w.WriteHeader(scode)
		w.Header().Set("Connection", "close")
	}
	log.Println(r.RemoteAddr, "\t", r.Method, "-", scode, "\t", r.RequestURI, "\t", len(body))
}

func getSecret(user, realm string) string {
	for i := 0; i < len(users); i++ {
		curUser := strings.Split(users[i], "=")
		if curUser[0] == user {
			return curUser[1]
		}
	}
	return ""
}
