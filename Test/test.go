package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	gRWLock      sync.RWMutex
	Certs        []Certificate
	Certificates []string
	NginxServer  string
)

type Certificate struct {
	Version        string `json:"version"`
	SerialNumber   string `json:"serialNumber"`
	Signature      string `json:"signatureName"`
	Issuer         string `json:"issuer"`
	IssuerCA       string `json:"issuerCA"`
	ValidityPeriod string `json:"validityPeriod"`
	ABSUID         string `json:"ABSUID"`
	ABSAttribute   string `json:"ABSAttribute"`
}

func init() {
	NginxServer = "http://10.176.40.46"
}

type CertificateResponse struct {
	CertificateContent Certificate `json:"certificate"`
	ABSSign            string      `json:"absSignature"`
}

func GenTest(uid string) string {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(NginxServer + "/dpki/ApplyForABSCertificate?uid=" + uid)
	// resp, err := client.Get("http://127.0.0.1:8001/ApplyForABSCertificate?uid=" + uid)
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
	resp, err := client.Get(NginxServer + "/dpki/VerifyABSCertificate?no=" + no)
	// resp, err := client.Get("http://127.0.0.1:8001/VerifyABSCertificate?no=" + no)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return true
}

func completeVerifyTest(bytesData string) bool {
	res, err := http.Post(NginxServer+"/dpki/VerifyABSCert", "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return true
}

func getCertTest(no string) bool {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(NginxServer + "/dpki/GetCertificate?no=" + no)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return true
}

func revokeCertTest(no string) bool {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(NginxServer + "/dpki/RevokeABSCertificate?no=" + no)
	if err != nil {
		return false
	}
	defer res.Body.Close()
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
				a[j*100+i] = float64(sEnd-sBegin) / 1e9
				var cer CertificateResponse
				if err := json.Unmarshal([]byte(sign), &cer); err != nil {
					return
				}
				VerifyTest(cer.CertificateContent.SerialNumber)
				sEnd = time.Now().UnixNano()
				b[j*100+i] = float64(sEnd-sBegin) / 1e9
			}(strconv.Itoa(i+10000), j, i)

		}
		wg.Wait()
	}

	end := time.Now().UnixNano()
	fmt.Print("Total time: ")
	fmt.Println(float64(end-start) / 1e9)

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

func Benchmark_Singletest(b *testing.B) {
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

func applyTest(num int) {
	fmt.Println("----------------------证书申请test ---------------------")
	fmt.Println("Apply test ---------------------")
	fmt.Println("ABS gen: ")
	start := time.Now().UnixNano()
	a := make([]float64, num, num)
	var valid int = 0
	Certs = make([]Certificate, num, num)
	Certificates = make([]string, num, num)
	certificates := make([]CertificateResponse, num, num)
	for j := 0; j < num/100; j += 1 {
		var wg sync.WaitGroup
		for i := 0; i < 100; i += 1 {
			wg.Add(1)
			go func(uid string, j int, i int) {
				defer wg.Done()
				sBegin := time.Now().UnixNano()
				sign := GenTest(uid)
				sEnd := time.Now().UnixNano()
				a[j*100+i] = float64(sEnd-sBegin) / 1e9
				var cer CertificateResponse
				if err := json.Unmarshal([]byte(sign), &cer); err != nil {
					return
				} else {
					certificates[j*100+i] = cer
					Certs[j*100+i] = cer.CertificateContent
					CertStr, _ := json.Marshal(cer)
					Certificates[j*100+i] = string(CertStr)
					gRWLock.Lock()
					valid += 1
					gRWLock.Unlock()
				}
			}(strconv.Itoa(j*100+i+10000), j, i)
		}
		wg.Wait()
	}
	end := time.Now().UnixNano()
	fmt.Print("Generate Total time: ")
	fmt.Println(float64(end-start) / 1e9)

	var avgGen float64
	var max float64 = -1
	var min float64 = 10
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
	write := bufio.NewWriter(file)
	for _, cert := range Certs {
		jsonCert, err := json.Marshal(cert)
		if err != nil {
			panic(err)
		}
		write.Write(jsonCert)
		write.WriteString("\n")
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	file.Close()

	filePath2 := "./certificates.txt"
	file2, err := os.OpenFile(filePath2, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	write2 := bufio.NewWriter(file2)
	for _, cert := range certificates {
		jsonCert, err := json.Marshal(cert)
		if err != nil {
			panic(err)
		}
		write2.Write(jsonCert)
		write2.WriteString("\n")
	}
	//Flush将缓存的文件真正写入到文件中
	write2.Flush()
	file2.Close()
}

func verify(num int) {
	fmt.Println("----------------------证书验证test ---------------------")
	fmt.Println("证书载入中 ---------------------")
	var valid int = 0
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
	for id, _ := range randCert {
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
				cer := certs[randCert[j*100+i]]
				if VerifyTest(cer.SerialNumber) {
					gRWLock.Lock()
					valid += 1
					gRWLock.Unlock()
				}
				sEnd := time.Now().UnixNano()
				a[j*100+i] = float64(sEnd-sBegin) / 1e9
			}(j, i)
		}
		wg.Wait()
	}
	end := time.Now().UnixNano()
	fmt.Print("Verify Total time: ")
	fmt.Println(float64(end-start) / 1e9)

	var avgVer float64
	var max float64 = -1
	var min float64 = 10
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

func completeVerify(num int) {
	fmt.Println("----------------------证书验证test ---------------------")
	fmt.Println("证书载入中 ---------------------")
	var valid int = 0
	a := make([]float64, num, num)
	certs := make([]string, 0, num)
	filePath := "./certificates.txt"
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	reader := bufio.NewReader(file)
	for {
		jsonCert, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		certs = append(certs, jsonCert)
	}
	file.Close()
	//随机选取保存的证书进行验证
	randCert := make([]int, num, num)
	for id, _ := range randCert {
		randCert[id] = rand.Intn(len(certs))
	}
	fmt.Println("完整证书验证 test ---------------------")
	start := time.Now().UnixNano()
	for j := 0; j < num/100; j += 1 {
		var wg sync.WaitGroup
		for i := 0; i < 100; i += 1 {
			wg.Add(1)
			go func(j int, i int) {
				defer wg.Done()
				sBegin := time.Now().UnixNano()
				cer := certs[randCert[j*100+i]]
				if completeVerifyTest(cer) {
					gRWLock.Lock()
					valid += 1
					gRWLock.Unlock()
				}
				sEnd := time.Now().UnixNano()
				a[j*100+i] = float64(sEnd-sBegin) / 1e9
			}(j, i)
		}
		wg.Wait()
	}
	end := time.Now().UnixNano()
	fmt.Print("Verify Total time: ")
	fmt.Println(float64(end-start) / 1e9)

	var avgVer float64
	var max float64 = -1
	var min float64 = 10
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

func HighConcurrencyTest() {
	num := 100000
	var valid int = 0
	a := make([]float64, num, num)
	randCert := make([]int, num, num)
	for id, _ := range randCert {
		randCert[id] = rand.Intn(len(Certificates))
	}
	fmt.Println("-------------------------十万并发测试 test ---------------------")
	start := time.Now().UnixNano()
	for j := 0; j < num/100; j += 1 {
		var wg sync.WaitGroup
		// 申请请求 10%
		for i := 0; i < 10; i += 1 {
			wg.Add(1)
			go func(uid string, j int, i int) {
				defer wg.Done()
				sBegin := time.Now().UnixNano()
				sign := GenTest(uid)
				sEnd := time.Now().UnixNano()
				a[j*100+i] = float64(sEnd-sBegin) / 1e9
				var cer CertificateResponse
				if err := json.Unmarshal([]byte(sign), &cer); err != nil {
					fmt.Println("证书申请解析错误：", sign)
				} else {
					gRWLock.Lock()
					valid += 1
					gRWLock.Unlock()
				}
			}(strconv.Itoa(j*100+i+10000), j, i)
		}
		// 验证请求 50%
		for i := 0; i < 50; i += 1 {
			wg.Add(1)
			go func(j int, i int) {
				defer wg.Done()
				sBegin := time.Now().UnixNano()
				cer := Certificates[randCert[j*100+i]]
				if completeVerifyTest(cer) {
					gRWLock.Lock()
					valid += 1
					gRWLock.Unlock()
				}
				sEnd := time.Now().UnixNano()
				a[j*100+i] = float64(sEnd-sBegin) / 1e9
			}(j, i)
		}
		// 查询请求 30%
		for i := 0; i < 30; i += 1 {
			wg.Add(1)
			go func(j int, i int) {
				defer wg.Done()
				sBegin := time.Now().UnixNano()
				cer := Certs[randCert[j*100+i]]
				if getCertTest(cer.SerialNumber) {
					gRWLock.Lock()
					valid += 1
					gRWLock.Unlock()
				}
				sEnd := time.Now().UnixNano()
				a[j*100+i] = float64(sEnd-sBegin) / 1e9
			}(j, i)
		}
		// 撤销请求 10%
		for i := 0; i < 10; i += 1 {
			wg.Add(1)
			go func(j int, i int) {
				defer wg.Done()
				sBegin := time.Now().UnixNano()
				cer := Certs[randCert[j*100+i]]
				if revokeCertTest(cer.SerialNumber) {
					gRWLock.Lock()
					valid += 1
					gRWLock.Unlock()
				}
				sEnd := time.Now().UnixNano()
				a[j*100+i] = float64(sEnd-sBegin) / 1e9
			}(j, i)
		}
		wg.Wait()
	}
	end := time.Now().UnixNano()
	var avgVer float64
	var max float64 = -1
	var min float64 = 10
	for _, ga := range a {
		avgVer += ga
		if ga < min && ga > 0 {
			min = ga
		}
		if ga > max {
			max = ga
		}
	}
	fmt.Print("Total time: ")
	fmt.Println(float64(end-start) / 1e9)
	fmt.Print("Max response time: ")
	fmt.Println(max)
	fmt.Print("Min response time: ")
	fmt.Println(min)
	fmt.Print("Average time: ")
	fmt.Println(avgVer / float64(num))
	var rate float64 = (float64(valid) / float64(num)) * 100
	fmt.Print("Total success rate: ", rate, "%\n")
}

func main() {
	num := flag.Int("n", 1000, "number of test.")
	flag.Parse()
	// abs_test(*num)
	applyTest(*num)
	// verify(*num)
	completeVerify(*num)
	// HighConcurrencyTest()
}
