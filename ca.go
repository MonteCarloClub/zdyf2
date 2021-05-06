package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
)

var (
    gPort *int
)

func init() {
    gPort  = flag.Int("port", 9000, "ra server port.")
}

func main() {
    http.HandleFunc("/ApplyForABSCertificate", ApplyForABSCertificate)
    http.HandleFunc("/VerifyABSCertificate", VerifyABSCertificate)

    if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
        log.Fatalln(err)
    }
}