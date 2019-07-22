package my_charles

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		"Hi, This is an example of https service in golang!")
}

func Init() {
	//normal()
	//return
	http.HandleFunc("/", handler)
	certFile := "/home/peterq/dev/env/cert/server.bundle.crt"
	certKey := "/home/peterq/dev/env/cert/server.nopass.key"
	certKey, certFile = "", ""
	server := &http.Server{
		Addr:    ":8081",
		Handler: nil,
		TLSConfig: &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				log.Println(info.ServerName, info.Conn.LocalAddr(), info.Conn.RemoteAddr())
				return getCert(info)
			},
		},
	}
	server.ListenAndServeTLS(certFile, certKey)
}

func normal() {
	http.HandleFunc("/", handler)
	certFile := "/home/peterq/dev/env/cert/server.bundle.crt"
	certKey := "/home/peterq/dev/env/cert/server.nopass.key"
	http.ListenAndServeTLS(":8081", certFile, certKey, nil)
}

var certMap = make(map[string]*tls.Certificate)
var certMapLock = new(sync.Mutex)

func getCert(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cn := info.ServerName
	if cn == "" {
		cn = info.Conn.LocalAddr().String()
		cn = strings.Split(cn, ":")[0]
	}
	certMapLock.Lock()
	fmt.Scanln()
	defer certMapLock.Unlock()
	if cert, ok := certMap[cn]; ok {
		return cert, nil
	}
	i := CertInformation{
		Country:            []string{"CN"},
		Organization:       []string{"PeterQ Info Tech .Ltd"},
		OrganizationalUnit: []string{"cert for " + cn},
		EmailAddress:       []string{"me@peterq.cn"},
		Province:           []string{"HuNan"},
		Locality:           nil,
		CommonName:         cn,
		CrtName:            "",
		KeyName:            "",
		IsCA:               false,
		Names:              nil,
	}
	certBuff, keyBuff, err := genCertWithCa(i)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(certBuff.Bytes(), keyBuff.Bytes())
	if err != nil {
		return nil, err
	}
	certMap[cn] = &cert
	return &cert, nil
}
