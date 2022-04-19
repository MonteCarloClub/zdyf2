package main

import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"
)

var (
    gPort *int
    Version string
    IssuerName string
    CertificateMap map[string]CertificateResponse
    UserID map[string]string
    gRWLock             sync.RWMutex
)

func init() {
    gPort  = flag.Int("port", 8000, "ra server port.")
    Version = "1.0"
    IssuerName = "RA-" + *flag.String("name", "1", "ra name")
    CertificateMap = make(map[string]CertificateResponse)
    UserID = make(map[string]string)
    flag.Parse()
}

// 申请证书
func ApplyForABSCertificate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    _ = r.ParseForm()
    uid := r.Form.Get("uid")
    attribute := r.Form.Get("attribute")

    timeStr := time.Now()
    serialNumber := uid + "-" + strconv.FormatInt(timeStr.UnixNano(), 10)
    c := Certificate{
        Version: Version,
        SerialNumber: serialNumber,
        Signature: "Attribute-based Signature",
        Issuer: IssuerName,
        ValidityPeriod: strconv.FormatInt(timeStr.UnixNano(), 10),
        ABSUID: uid,
        ABSAttribute: attribute,
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
    UserID[uid] = serialNumber
    gRWLock.Unlock()

    bData, _ := json.Marshal(res)
    _, _ = fmt.Fprintf(w, string(bData))
}

// 验证证书
func VerifyABSCertificate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

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

// 撤销证书
func RevokeABSCertificate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

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

// 获取当前证书数量
func GetCertificateNumber(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    gRWLock.RLock()
    http.Error(w, strconv.Itoa(len(CertificateMap)), 200)
    gRWLock.RUnlock()
}

//func ConcurrencyTest(w http.ResponseWriter, r *http.Request) {
//    gRWLock.RLock()
//    http.Error(w, strconv.Itoa(len(CertificateMap)), 200)
//    gRWLock.RUnlock()
//}

// IoT 设备证书测试，获取所有证书
func IoTDevTest(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    var a []CertificateResponse
    for _, v := range CertificateMap {
        a = append(a, v)
    }
    bData, _ := json.Marshal(a)
    _, _ = fmt.Fprintf(w, string(bData))
}

// 通过序列号获取证书
func GetCertificate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    _ = r.ParseForm()
    serialNumber := r.Form.Get("no")

    gRWLock.RLock()
    defer gRWLock.RUnlock()

    if res, ok := CertificateMap[serialNumber]; !ok {
        http.Error(w, "Certificate does not exist.", 500)
    } else {
        bData, _ := json.Marshal(res)
        _, _ = fmt.Fprintf(w, string(bData))
    }
}

// IoT 设备初始化
func IotDevInit(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    a := [200]CertificateResponse{}
    for i := 0; i < 200; i++ {
        client := &http.Client{Timeout: 10 * time.Second}
        uid := "iotdevice" + strconv.Itoa(i)
        attribute := "tag" +  strconv.Itoa(i) + "1,tag" +  strconv.Itoa(i) + "2,tag" + strconv.Itoa(i) + "3"
        resp, err := client.Get("http://127.0.0.1:8001/ApplyForABSCertificate?uid=" + uid + "&&attribute=" + attribute)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        s, _ := ioutil.ReadAll(resp.Body)

        var cer CertificateResponse
        _ = json.Unmarshal(s, &cer)
        UserID[uid] = cer.CertificateContent.SerialNumber
        a[i] = cer
        resp.Body.Close()
    }

    bData, _ := json.Marshal(a)
    _, _ = fmt.Fprintf(w, string(bData))
}

// 通过序列号撤销证书
func RevokeABSCertificateByUID(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    _ = r.ParseForm()
    uid := r.Form.Get("userid")
    serialNumber := UserID[uid]

    gRWLock.RLock()
    defer gRWLock.RUnlock()

    if _, ok := CertificateMap[serialNumber]; !ok {
        s := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

        resp := RevokeResponse{
            Status: "Certificate does not exist.",
            Timestamp: s,
            Tx: Sha256(strconv.FormatInt(time.Now().UnixNano(), 10)),
            SerialNumber: serialNumber,
        }
        bData, _ := json.Marshal(resp)
        _, _ = fmt.Fprintf(w, string(bData))
    } else {
        delete(CertificateMap, serialNumber)
        delete(UserID, uid)
        s := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
        resp := RevokeResponse{
            Status: "OK.",
            Timestamp: s,
            Tx: Sha256(strconv.FormatInt(time.Now().UnixNano(), 10)),
            SerialNumber: serialNumber,
        }
        bData, _ := json.Marshal(resp)
        _, _ = fmt.Fprintf(w, string(bData))
    }
}

func Sha256(src string) string {
    m := sha256.New()
    m.Write([]byte(src))
    res := hex.EncodeToString(m.Sum(nil))
    return res
}

// 登录界面
func login(w http.ResponseWriter, r *http.Request) {
    str, err := ioutil.ReadFile("./template/index.html")
    s, _ := os.Getwd()
    if err != nil {
        http.Error(w, s, 500)
        return
    }
    _, _ = w.Write([]byte(str))
}

// 通过用户 UID 获取证书
func GetCertificateByUID(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()
    uid := r.Form.Get("uid")
    serialNumber := UserID[uid]

    gRWLock.RLock()
    defer gRWLock.RUnlock()

    if res, ok := CertificateMap[serialNumber]; !ok {
        http.Error(w, "Certificate does not exist.", 500)
    } else {
        bData, _ := json.Marshal(res)
        _, _ = fmt.Fprintf(w, string(bData))
    }
}

func main() {
    http.HandleFunc("/ApplyForABSCertificate", ApplyForABSCertificate)
    http.HandleFunc("/VerifyABSCertificate", VerifyABSCertificate)
    http.HandleFunc("/RevokeABSCertificate", RevokeABSCertificate)
    http.HandleFunc("/RevokeABSCertificateByUID", RevokeABSCertificateByUID)

    http.HandleFunc("/GetCertificateNumber", GetCertificateNumber)
    http.HandleFunc("/GetCertificate", GetCertificate)
    http.HandleFunc("/GetCertificateByUID", GetCertificateByUID)

    http.HandleFunc("/IoTDevTest", IoTDevTest)
    http.HandleFunc("/IotDevInit", IotDevInit)

    http.HandleFunc("/login", login)
    http.Handle("/", http.FileServer(http.Dir("template")))

    //http.HandleFunc("/ConcurrencyTest", ConcurrencyTest)

    if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
        log.Fatalln(err)
    }
}