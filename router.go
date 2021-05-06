package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func ApplyForABSCertificate(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()

    m := r.Form.Get("m")
    sign := Generate(m)

    bdata, _ := json.Marshal(sign)
    _, _ = fmt.Fprintf(w, string(bdata))
}

func VerifyABSCertificate(w http.ResponseWriter, r *http.Request) {
    postdata, _ := ioutil.ReadAll(r.Body)

    var sign *ABSSignature
    if err := json.Unmarshal(postdata, &sign); err != nil {
        http.Error(w, "parse response error.", 500)
        return
    }

    if flag := Verify(sign); flag {
        http.Error(w, "OK", 200)
    } else {
        http.Error(w, "NO", 500)
    }
}
