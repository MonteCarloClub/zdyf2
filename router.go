package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)

func ApplyForABSCertificate(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()

    uid := r.Form.Get("uid")
    timeStr := time.Now()
    c := Certificate{
        Version: Version,
        SerialNumber: uid + "-" + strconv.FormatInt(timeStr.UnixNano(), 10),
        Signature: "Attribute-based Signature",
        Issuer: IssuerName,
        ValidityPeriod: strconv.FormatInt(timeStr.UnixNano(), 10),
        ABSUID: uid,
    }

    b, _ := json.Marshal(c)
    sign := Generate(string(b))

    res := CertificateResponse{
        CertificateContent: c,
        ABSSign: *sign,
    }

    bdata, _ := json.Marshal(res)
    _, _ = fmt.Fprintf(w, string(bdata))
}

func VerifyABSCertificate(w http.ResponseWriter, r *http.Request) {
    postdata, _ := ioutil.ReadAll(r.Body)

    var resp *CertificateResponse
    if err := json.Unmarshal(postdata, &resp); err != nil {
        http.Error(w, "parse response error.", 500)
        return
    }

    cer := resp.CertificateContent
    valid := cer.ValidityPeriod
    fmt.Println(valid)
    fmt.Println(time.Now().UnixNano())
    if valid < strconv.FormatInt(time.Now().UnixNano(), 10) {
        http.Error(w, "The certificate has expired.", 500)
        return
    }

    sign := &resp.ABSSign
    if flag := Verify(sign); flag {
        http.Error(w, "OK.", 200)
    } else {
        http.Error(w, "Failed to Verify.", 500)
    }
}
