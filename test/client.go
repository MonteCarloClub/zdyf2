package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"
    "sync"
    "time"
    "testing"
    "bufio"
    "os"
    "math/rand"
)

type Certificate struct {
    Version string        `json:"version"`
    SerialNumber string   `json:"serialNumber"`
    Signature string      `json:"signatureName"`
    Issuer string         `json:"issuer"`
    ValidityPeriod string `json:"validityPeriod"`
    ABSUID string         `json:"ABSUID"`
}

type CertificateResponse struct {
    CertificateContent Certificate `json:"certificate"`
    ABSSign string                 `json:"absSignature"`
}

func GenTest(uid string) string {
    client := &http.Client{Timeout: 10 * time.Second}

    resp, err := client.Get("http://127.0.0.1:8001/ApplyForABSCertificate?uid=" + uid)
    if err != nil {
        return err.Error()
    }
    defer resp.Body.Close()

    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err.Error()
    }

    cer := string(content)
    return cer
}

func VerifyTest(no string) bool {
    client := &http.Client{Timeout: 10 * time.Second}

    resp, err := client.Get("http://127.0.0.1:8001/VerifyABSCertificate?no=" + no)
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return true
}

func abs_test(num int) {
    fmt.Println("ABS test ---------------------")
    fmt.Println("ABS gen & verify: ")
    start := time.Now().UnixNano()
    var a [20000]float64
    var b [20000]float64

    for j := 0; j < num/100; j += 1 {
        var wg sync.WaitGroup

        for i := 0; i < 100; i += 1 {
            wg.Add(1)

            go func(uid string, j int, i int) {
                defer wg.Done()
                sBegin := time.Now().UnixNano()
                sign := GenTest(uid)
                sEnd := time.Now().UnixNano()
                a[j * 100 + i] = float64(sEnd - sBegin) / 1e9
                var cer CertificateResponse
                if err := json.Unmarshal([]byte(sign), &cer); err != nil {
                    return
                }
                VerifyTest(cer.CertificateContent.SerialNumber)
                sEnd = time.Now().UnixNano()
                b[j * 100 + i] = float64(sEnd - sBegin) / 1e9
            }(strconv.Itoa(i + 10000), j, i)

        }
        wg.Wait()
    }

    end := time.Now().UnixNano()
    fmt.Print("Total time: ")
    fmt.Println(float64(end - start) / 1e9)

    var avgGen float64
    var avgVerify float64
    for _, ga := range a {
        avgGen += ga
    }
    for _, gb := range b {
        avgVerify += gb
    }

    fmt.Print("Average time of ABS gen: ")
    fmt.Println(avgGen / float64(num))
    fmt.Print("Average time of ABS verify: ")
    fmt.Println(avgVerify / float64(num))
}

func Benchmark_Singletest(b *testing.B){  
    var n int
	for i := 0; i < b.N; i++ {
        uid := strconv.Itoa(i + 10000)
        sign := GenTest(uid)
        var cer CertificateResponse
        if err := json.Unmarshal([]byte(sign), &cer); err != nil {
         return
        }
        VerifyTest(cer.CertificateContent.SerialNumber)  
		n++
	}  
}

func applyTest(num int){
    fmt.Println("Apply test ---------------------")
    fmt.Println("ABS gen: ")
    start := time.Now().UnixNano()
    a := make([]float64, num, num)
    var valid int=0
    certs := make([]Certificate, num, num)
    for j := 0; j < num/100; j += 1 {
        var wg sync.WaitGroup
        for i := 0; i < 100; i += 1 {
            wg.Add(1)
            go func(uid string, j int, i int) {
                defer wg.Done()
                sBegin := time.Now().UnixNano()
                sign := GenTest(uid)
                sEnd := time.Now().UnixNano()
                a[j * 100 + i] = float64(sEnd - sBegin) / 1e9
                var cer CertificateResponse
                if err := json.Unmarshal([]byte(sign), &cer); err != nil {
                    return
                }else {
                    certs[j * 100 + i] = cer.CertificateContent;
                    valid += 1
                }
            }(strconv.Itoa(j * 100 + i + 10000), j, i)
        }
        wg.Wait()
    }
    end := time.Now().UnixNano()
    fmt.Print("Generate Total time: ")
    fmt.Println(float64(end - start) / 1e9)

    var avgGen float64
    var max float64=-1
    var min float64=10
    for _, ga := range a {
        avgGen += ga
        if ga < min {
            min = ga
        } 
        if ga > max {
            max = ga
        } 
    }
    fmt.Print("Max response time of Generate: ")
    fmt.Println(max)
    fmt.Print("Min response time of Generate: ")
    fmt.Println(min)
    fmt.Print("Average time of ABS gen: ")
    fmt.Println(avgGen / float64(num))
    var rate float64 = (float64(valid) / float64(num)) * 100
    fmt.Print("Total success rate: ", rate, "%\n")
    
    filePath := "./certs.txt"
    file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
    if err != nil {
        fmt.Println("文件打开失败", err)
    }
    defer file.Close()
    write := bufio.NewWriter(file)
    for _, cert := range certs {
        jsonCert, err := json.Marshal(cert)
        if err != nil {
            panic(err)
        }
        write.Write(jsonCert)
        write.WriteString("\n")
    }
    //Flush将缓存的文件真正写入到文件中
    write.Flush()
}

func verify(num int){
    var valid int=0
    a := make([]float64, num, num)
    certs := make([]Certificate, 0, num)
    filePath := "./certs.txt"
    file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
    if err != nil {
        fmt.Println("文件打开失败", err)
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    for {
        jsonCert, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        }
        var cert Certificate
        err = json.Unmarshal([]byte(jsonCert), &cert)
        if err != nil {
            fmt.Println("error:", err)
        }
        certs = append(certs, cert)
    }
    randCert := make([]int, num, num)
    for id, _ := range randCert{
        randCert[id] = rand.Intn(len(certs))
    }
    fmt.Println("Verify test ---------------------")
    start := time.Now().UnixNano()
    for j := 0; j < num/100; j += 1 {
        var wg sync.WaitGroup
        for i := 0; i < 100; i += 1 {
            wg.Add(1)
            go func(j int, i int) {
                defer wg.Done()
                sBegin := time.Now().UnixNano()
                cer := certs[randCert[j * 100 + i]]
                if VerifyTest(cer.SerialNumber) {
                    valid += 1
                }
                sEnd := time.Now().UnixNano()
                a[j * 100 + i] = float64(sEnd - sBegin) / 1e9
            }(j, i)
        }
        wg.Wait()
    }
    end := time.Now().UnixNano()
    fmt.Print("Verify Total time: ")
    fmt.Println(float64(end - start) / 1e9)

    var avgVer float64
    var max float64=-1
    var min float64=10
    for _, ga := range a {
        avgVer += ga
        if ga < min {
            min = ga
        } 
        if ga > max {
            max = ga
        } 
    }
    fmt.Print("Max response time of Verify: ")
    fmt.Println(max)
    fmt.Print("Min response time of Verify: ")
    fmt.Println(min)
    fmt.Print("Average time of Verify: ")
    fmt.Println(avgVer / float64(num))
    var rate float64 = (float64(valid) / float64(num)) * 100
    fmt.Print("Total success rate: ", rate, "%\n")
    
}

func main() {
    num := flag.Int("n", 1000, "number of test.")
    flag.Parse()
    // abs_test(*num)
    // applyTest(*num)
    // verify(*num)
}
