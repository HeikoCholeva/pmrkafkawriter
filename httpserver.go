package main

import (
	"log"
	"context"
	"strings"
	"net/http"
	"io/ioutil"
	pmrlib "github.com/HeikoCholeva/pmrlib-go"
)

var srv http.Server

func stopServer() {
	err := srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func startServer(port string, path string, tls bool, cert string, key string) {
	srv.Addr = port
	log.Println("Starting server on port", srv.Addr)

	http.HandleFunc(path, handler)
	if tls {
		log.Fatal(srv.ListenAndServeTLS(cert, key))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
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
                err = report.FromJSON(body)
                if err != nil {
			scode = 400
			w.WriteHeader(scode)
			w.Header().Set("Connection", "close")
                } else {
                        jsonstr, err := pmrlib.ToJSON(report)
                        if err != nil {
                                panic(err)
                        } else {
                               	writeToKafka(jsonstr)
				scode = 202
				w.WriteHeader(scode)
                                w.Header().Set("Connection", "close")
                        }
                }
        } else {
		scode = 405
        	w.WriteHeader(scode)
        	w.Header().Set("Connection", "close")
	}
	log.Println(r.RemoteAddr,"\t", r.Method, "-", scode, "\t", r.RequestURI, "\t", len(body))
}
