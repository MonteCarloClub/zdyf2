package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

func GenTest(str string) {
    client := &http.Client{Timeout: 10 * time.Second}

    resp, err := client.Get("http://127.0.0.1:9000/ApplyForABSCertificate?m=" + str)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }

    var sign *ABSSignature
    if err := json.Unmarshal(content, &sign); err != nil {
        return
    }

    fmt.Println(string(content))
    VerifyTest(sign)
}

func VerifyTest(sign *ABSSignature) {
    body, err := json.Marshal(sign)
    if err != nil {
        return
    }

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Post("http://127.0.0.1:9000/VerifyABSCertificate", "application/json", bytes.NewReader(body))
    if err != nil {
        return
    }
    defer resp.Body.Close()

    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }

    fmt.Println(string(content))
}

func main() {
    GenTest("123")
}
