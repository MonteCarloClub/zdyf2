package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
)

var (
    gPort *int
    Version string
    IssuerName string
)

func init() {
    gPort  = flag.Int("port", 9000, "ra server port.")
    Version = "1.0"
    IssuerName = "RA-" + *flag.String("name", "1", "ra name")
}

func main() {
    http.HandleFunc("/ApplyForABSCertificate", ApplyForABSCertificate)
    http.HandleFunc("/VerifyABSCertificate", VerifyABSCertificate)

    if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
        log.Fatalln(err)
    }
}