package main

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "time"
)

var (
    gPort *int
    Version string
    IssuerName string
    CertificateMap map[string]CertificateResponse
)

func init() {
    gPort  = flag.Int("port", 9000, "ra server port.")
    Version = "1.0"
    IssuerName = "RA-" + *flag.String("name", "1", "ra name")
    CertificateMap = make(map[string]CertificateResponse)
}

func ApplyForABSCertificate(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()
    uid := r.Form.Get("uid")

    timeStr := time.Now()
    serialNumber := uid + "-" + strconv.FormatInt(timeStr.UnixNano(), 10)
    c := Certificate{
        Version: Version,
        SerialNumber: serialNumber,
        Signature: "Attribute-based Signature",
        Issuer: IssuerName,
        ValidityPeriod: strconv.FormatInt(timeStr.UnixNano(), 10),
        ABSUID: uid,
    }

    b, _ := json.Marshal(c)
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Post("http://127.0.0.1:9001/SingleGenerate", "application/json", bytes.NewReader(b))
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer resp.Body.Close()

    sign, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    res := CertificateResponse{
        CertificateContent: c,
        ABSSign: string(sign),
    }
    CertificateMap[serialNumber] = res
    bData, _ := json.Marshal(res)
    _, _ = fmt.Fprintf(w, string(bData))
}

func VerifyABSCertificate(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()
    serialNumber := r.Form.Get("no")

    if res, ok := CertificateMap[serialNumber]; !ok {
        http.Error(w, "Certificate does not exist.", 500)
    } else {
        valid := res.CertificateContent.ValidityPeriod
        fmt.Println(valid)
        fmt.Println(time.Now().UnixNano())
        if valid == strconv.FormatInt(time.Now().UnixNano(), 10) {
            http.Error(w, "The certificate has expired.", 500)
        } else {
            sign := res.ABSSign
            client := &http.Client{Timeout: 10 * time.Second}
            resp, err := client.Post("http://127.0.0.1:9001/SingleVerify", "application/json", bytes.NewReader([]byte(sign)))
            if err != nil {
                http.Error(w, err.Error(), 500)
                return
            }
            defer resp.Body.Close()

            _, _ = fmt.Fprintf(w, "OK.")
        }
    }
}

func main() {
    http.HandleFunc("/ApplyForABSCertificate", ApplyForABSCertificate)
    http.HandleFunc("/VerifyABSCertificate", VerifyABSCertificate)

    if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
        log.Fatalln(err)
    }
}