package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"flag"

)

var (
    gPort *int
    Version string
    // IssuerName string
    // CertificateMap map[string]CertificateResponse
    // UserID map[string]string
    // gRWLock             sync.RWMutex
)

func init() {
    gPort  = flag.Int("port", 8090, "ra server port.")
    Version = "1.0"
    // IssuerName = "RA-" + *flag.String("name", "1", "ra name")
    // CertificateMap = make(map[string]CertificateResponse)
    // UserID = make(map[string]string)
    flag.Parse()
}

func GetCertificate(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()
    serialNumber := r.Form.Get("serialNumber")
	args := make([][]byte, 0)
	args = append(args, []byte(serialNumber))
	resp, err := ChannelExecute("getCertificate", args)
	if err != nil {
	 	http.Error(w, "Certificate does not exist.", 500)
	 	return
	}else {
        bData, _ := json.Marshal(resp)
        _, _ = fmt.Fprintf(w, string(bData))
    }
}

func main() {
    http.HandleFunc("/GetCertificate", GetCertificate)
    // http.HandleFunc("/ApplyForABSCertificate", ApplyForABSCertificate)
    // http.HandleFunc("/VerifyABSCertificate", VerifyABSCertificate)
    // http.HandleFunc("/RevokeABSCertificate", RevokeABSCertificate)
    // http.HandleFunc("/RevokeABSCertificateByUID", RevokeABSCertificateByUID)

    // http.HandleFunc("/GetCertificateNumber", GetCertificateNumber)

    // http.HandleFunc("/GetCertificateByUID", GetCertificateByUID)

    // http.HandleFunc("/IoTDevTest", IoTDevTest)
    // http.HandleFunc("/IotDevInit", IotDevInit)

    // http.HandleFunc("/login", login)
    // http.Handle("/", http.FileServer(http.Dir("template")))

    //http.HandleFunc("/ConcurrencyTest", ConcurrencyTest)

    if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
        log.Fatalln(err)
    }
}