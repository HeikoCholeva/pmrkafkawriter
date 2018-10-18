package main

import (
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	pmrlib "github.com/HeikoCholeva/pmrlib-go"
)

var srv http.Server

func StartServer(port string, path string, tls bool, cert string, key string) {
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
	defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
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
