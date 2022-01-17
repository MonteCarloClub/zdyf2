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
    "sync"
    "time"
)

var (
    gPort *int
    Version string
    IssuerName string
    CertificateMap map[string]CertificateResponse
    gRWLock             sync.RWMutex
)

func init() {
    gPort  = flag.Int("port", 8000, "ra server port.")
    Version = "1.0"
    IssuerName = "RA-" + *flag.String("name", "1", "ra name")
    CertificateMap = make(map[string]CertificateResponse)
    flag.Parse()
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
        http.Error(w, err.Error(), 500)
        return
    }
    defer resp.Body.Close()

    sign, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    res := CertificateResponse{
        CertificateContent: c,
        ABSSign: string(sign),
    }

    gRWLock.Lock()
    CertificateMap[serialNumber] = res
    gRWLock.Unlock()

    bData, _ := json.Marshal(res)
    _, _ = fmt.Fprintf(w, string(bData))
}

func VerifyABSCertificate(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()
    serialNumber := r.Form.Get("no")

    gRWLock.RLock()
    defer gRWLock.RUnlock()

    if res, ok := CertificateMap[serialNumber]; !ok {
        http.Error(w, "Certificate does not exist.", 500)
    } else {
        valid := res.CertificateContent.ValidityPeriod
        //fmt.Println(valid)
        //fmt.Println(time.Now().UnixNano())
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

func RevokeABSCertificate(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()
    serialNumber := r.Form.Get("no")

    gRWLock.RLock()
    defer gRWLock.RUnlock()

    if _, ok := CertificateMap[serialNumber]; !ok {
        http.Error(w, "Certificate does not exist.", 500)
    } else {
        delete(CertificateMap, serialNumber)
        _, _ = fmt.Fprintf(w, "OK.")
    }
}

func GetCertificateNumber(w http.ResponseWriter, r *http.Request) {
    gRWLock.RLock()
    http.Error(w, strconv.Itoa(len(CertificateMap)), 200)
    gRWLock.RUnlock()
}

//func ConcurrencyTest(w http.ResponseWriter, r *http.Request) {
//    gRWLock.RLock()
//    http.Error(w, strconv.Itoa(len(CertificateMap)), 200)
//    gRWLock.RUnlock()
//}

func main() {
    http.HandleFunc("/ApplyForABSCertificate", ApplyForABSCertificate)
    http.HandleFunc("/VerifyABSCertificate", VerifyABSCertificate)
    http.HandleFunc("/RevokeABSCertificate", RevokeABSCertificate)
    http.HandleFunc("/GetCertificateNumber", GetCertificateNumber)
    //http.HandleFunc("/ConcurrencyTest", ConcurrencyTest)

    if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
        log.Fatalln(err)
    }
}