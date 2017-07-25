package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	certPool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		fmt.Println("ReadFile err:", err)
	}
	certPool.AppendCertsFromPEM(caCrt)

	server := &http.Server{
		Addr:    ":8080",
		Handler: &TestHandler{},
		TLSConfig: &tls.Config{
			ClientCAs:  certPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}
	err = server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		fmt.Println(err)
	}
}

type TestHandler struct {
	http.Handler
}

func (t *TestHandler) ServerHttp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is a example of test https service")
}
