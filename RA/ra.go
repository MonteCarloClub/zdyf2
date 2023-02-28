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
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type FabricInfo struct {
	ABSUID       []byte
	SerialNumber []byte
	Cert         []byte
}

var (
	gPort      *int
	Version    string
	IssuerName string
	// CertificateMap map[string]CertificateResponse
	// UserID         map[string]string
	gRWLock   sync.RWMutex
	redisdb   *redis.Client
	CAPort    int
	CANum     int
	RAbase    int
	validHour time.Duration
	fbWorker  chan *FabricInfo
)

func init() {
	gPort = flag.Int("port", 8000, "ra server port.")
	Version = "1.0"
	raName := flag.Int("name", 0, "ra name")
	// CertificateMap = make(map[string]CertificateResponse)
	// UserID = make(map[string]string)
	flag.Parse()
	IssuerName = "RA-" + strconv.Itoa(*raName)
	RAbase = *raName
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 指定
		Password: "",
		DB:       0,
	})
	CAPort = 0      //CA的Port, 每调用一次CA， CAport + 1;
	CANum = 10      //每个RA对应的CA数量
	validHour = 720 //有效期时间小时数
	fbWorker = make(chan *FabricInfo, 10000)
	go fabricStore(fbWorker)
}

// 申请证书
func ApplyForABSCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = r.ParseForm()
	uid := r.Form.Get("uid")
	attribute := r.Form.Get("attribute")

	timeStr := time.Now().Add(time.Hour * validHour)
	serialNumber := uid + "-" + strconv.FormatInt(timeStr.UnixNano(), 10)
	c := Certificate{
		Version:        Version,
		SerialNumber:   serialNumber,
		Signature:      "Attribute-based Signature",
		Issuer:         IssuerName,
		IssuerCA:       "",
		ValidityPeriod: strconv.FormatInt(timeStr.UnixNano(), 10),
		ABSUID:         uid,
		ABSAttribute:   attribute,
	}
	b, _ := json.Marshal(c)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("http://127.0.0.1:"+strconv.Itoa(CAPort+9000+(RAbase-1)*CANum+1)+"/SingleGenerate", "application/json", bytes.NewReader(b))
	CAPort = (CAPort + 1) % CANum
	c.IssuerCA = IssuerName + "-CA-" + strconv.Itoa(CAPort+9000+(RAbase-1)*CANum+1)
	if err != nil {
		log.Printf("[Apply]CA SingleGenerate: %s 失败", c.SerialNumber)
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	sign, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Apply]CA SingleGenerate read info: %s 失败", c.SerialNumber)
		http.Error(w, err.Error(), 500)
		return
	}

	res := CertificateResponse{
		CertificateContent: c,
		ABSSign:            string(sign),
	}

	bData, _ := json.Marshal(res)
	_, _ = fmt.Fprintf(w, string(bData))
	// gRWLock.Lock()
	// gRWLock.Unlock()
	// CertificateMap[serialNumber] = res
	// UserID[uid] = serialNumber
	// start := time.Now().UnixNano()
	// redis store
	redisdb.Set(serialNumber, string(bData), time.Hour*validHour)
	// end := time.Now().UnixNano()
	// log.Printf("Redis Total time: %f", float64(end-start)/1e9)
	// fabric store
	fbInfo := FabricInfo{
		ABSUID:       []byte(c.ABSUID),
		SerialNumber: []byte(c.SerialNumber),
		Cert:         []byte(b),
	}
	fbWorker <- &fbInfo

}

func fabricStore(fbInfochan chan *FabricInfo) {
	var err error
	for {
		fbInfo := <-fbInfochan
		go func(fbInfo *FabricInfo) {
			if fbInfo == nil {
				return
			}
			args := make([][]byte, 0)
			args = append(args, fbInfo.ABSUID)
			args = append(args, fbInfo.SerialNumber)
			args = append(args, fbInfo.Cert)
			_, err = ChannelExecute("setCertificate", args)
			if err != nil {
				log.Printf("Fabric setCertificate: %s 失败 %s", fbInfo.SerialNumber, err.Error())
				time.Sleep(time.Millisecond * 1000)
				fbWorker <- fbInfo
			} else {
				log.Printf("[Apply]Fabric setCertificate: %s 成功", fbInfo.SerialNumber)
			}
		}(fbInfo)
	}
}

// 验证证书
func VerifyABSCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = r.ParseForm()
	serialNumber := r.Form.Get("no")

	rawData, err := redisdb.Get(serialNumber).Result()
	if err != nil {
		http.Error(w, "Certificate does not exist.", 500)
	} else {
		var res CertificateResponse
		if err := json.Unmarshal([]byte(rawData), &res); err != nil {
			return
		}
		valid := res.CertificateContent.ValidityPeriod
		if valid < strconv.FormatInt(time.Now().UnixNano(), 10) {
			http.Error(w, "The certificate has expired.", 500)
			return
		}
		sign := res.ABSSign
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post("http://127.0.0.1:"+strconv.Itoa(CAPort+9000+(RAbase-1)*CANum)+"/SingleVerify", "application/json", bytes.NewReader([]byte(sign)))
		CAPort = (CAPort + 1) % CANum
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()
		_, _ = fmt.Fprintf(w, "True")
	}
}

func VerifyABSCert(writer http.ResponseWriter, request *http.Request) {
	// 检查是否POST请求
	if request.Method != "POST" {
		writer.WriteHeader(405)
		return
	}
	// 解析form
	err := request.ParseForm()
	if err != nil {
		http.Error(writer, err.Error(), 400)
		return
	}

	rawCert, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("Read failed:", err)
	}
	defer request.Body.Close()
	var cert CertificateResponse
	if err := json.Unmarshal([]byte(rawCert), &cert); err != nil {
		http.Error(writer, err.Error(), 500)
	} else {
		SNumber := cert.CertificateContent.SerialNumber
		log.Println("[Verify]Complete Certificate verify:", SNumber)
		rawData, err := redisdb.Get(SNumber).Result()
		if err != nil {
			http.Error(writer, "The certificate is invalid.", 500)
			return
		}
		if rawData == string(rawCert) {
			// _, _ = fmt.Fprintf(writer, "True compare")
			valid := cert.CertificateContent.ValidityPeriod
			if valid < strconv.FormatInt(time.Now().UnixNano(), 10) {
				http.Error(writer, "The certificate has expired.", 500)
				return
			}
			sign := cert.ABSSign
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Post("http://127.0.0.1:"+strconv.Itoa(CAPort+9000+(RAbase-1)*CANum)+"/SingleVerify", "application/json", bytes.NewReader([]byte(sign)))
			CAPort = (CAPort + 1) % CANum
			if err != nil {
				http.Error(writer, err.Error(), 500)
				return
			}
			defer resp.Body.Close()
			_, _ = fmt.Fprintf(writer, "True")

		} else {
			log.Println("[Verify] Certificate verify failed,certificate is invalid:")
			http.Error(writer, "The certificate is invalid.", 500)
			return
			// _, _ = fmt.Fprintf(writer, "False compare")
		}
	}
}

// 通过序列号获取证书
func GetCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = r.ParseForm()
	serialNumber := r.Form.Get("no")

	// gRWLock.RLock()
	// defer gRWLock.RUnlock()

	rawData, err := redisdb.Get(serialNumber).Result()
	if err != nil {
		http.Error(w, "Certificate does not exist.", 500)
	} else {
		_, _ = fmt.Fprintf(w, rawData)
	}
	// if res, ok := CertificateMap[serialNumber]; !ok {
	// 	http.Error(w, "Certificate does not exist.", 500)
	// } else {
	// 	bData, _ := json.Marshal(res)
	// 	_, _ = fmt.Fprintf(w, string(bData))
	// }
}

// 撤销证书
func RevokeABSCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = r.ParseForm()
	serialNumber := r.Form.Get("no")

	// gRWLock.RLock()
	// defer gRWLock.RUnlock()
	_, err := redisdb.Del(serialNumber).Result()
	if err != nil {
		log.Printf("[Revoke]撤销证书: %s 失败 %s", serialNumber, err.Error())
		http.Error(w, "Certificate does not exist.", 500)
	} else {
		log.Printf("[Revoke]撤销证书: %s 成功", serialNumber)
		_, _ = fmt.Fprintf(w, "Revoke OK.")
	}
}

// 获取当前证书数量
func GetCertificateNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = redisdb.Scan(cursor, "*", 10).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)
		if cursor == 0 {
			break
		}
	}
	// gRWLock.RLock()
	http.Error(w, strconv.Itoa(n), 200)
	// gRWLock.RUnlock()
}

func GetCertificateFromFabric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = r.ParseForm()
	serialNumber := r.Form.Get("no")
	// gRWLock.RLock()
	// defer gRWLock.RUnlock()
	args := make([][]byte, 0)
	args = append(args, []byte(serialNumber))
	resp, err := ChannelExecute("getCertificate", args)
	if err != nil {
		http.Error(w, "Certificate does not exist.", 500)
		return
	} else {
		bData, _ := json.Marshal(resp)
		_, _ = fmt.Fprintf(w, string(bData))
	}
}

//func ConcurrencyTest(w http.ResponseWriter, r *http.Request) {
//    gRWLock.RLock()
//    http.Error(w, strconv.Itoa(len(CertificateMap)), 200)
//    gRWLock.RUnlock()
//}

// IoT 设备证书测试，获取所有证书
func IoTDevTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// var a []CertificateResponse
	// for _, v := range CertificateMap {
	// 	a = append(a, v)
	// }
	var a []string
	for {
		var cursor uint64
		var keys []string
		var err error
		keys, cursor, err = redisdb.Scan(cursor, "*", 10).Result()
		if err != nil {
			panic(err)
		}
		for _, v := range keys {
			rawData, err := redisdb.Get(v).Result()
			if err != nil {
				break
			}
			a = append(a, rawData)
		}
		if cursor == 0 {
			break
		}
	}
	bData, _ := json.Marshal(a)
	_, _ = fmt.Fprintf(w, string(bData))
}

// IoT 设备初始化
func IotDevInit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	a := [200]CertificateResponse{}
	for i := 0; i < 200; i++ {
		client := &http.Client{Timeout: 10 * time.Second}
		uid := "iotdevice" + strconv.Itoa(i)
		attribute := "tag" + strconv.Itoa(i) + "1,tag" + strconv.Itoa(i) + "2,tag" + strconv.Itoa(i) + "3"
		resp, err := client.Get("http://127.0.0.1:" + strconv.Itoa(8000+RAbase) + "/ApplyForABSCertificate?uid=" + uid + "&&attribute=" + attribute)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		s, _ := ioutil.ReadAll(resp.Body)

		var cer CertificateResponse
		_ = json.Unmarshal(s, &cer)
		// UserID[uid] = cer.CertificateContent.SerialNumber
		a[i] = cer
		resp.Body.Close()
	}

	bData, _ := json.Marshal(a)
	_, _ = fmt.Fprintf(w, string(bData))
}

// 通过UID撤销证书
// func RevokeABSCertificateByUID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	_ = r.ParseForm()
// 	uid := r.Form.Get("userid")
// 	serialNumber := UserID[uid]

// 	gRWLock.RLock()
// 	defer gRWLock.RUnlock()

// 	if _, ok := CertificateMap[serialNumber]; !ok {
// 		s := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

// 		resp := RevokeResponse{
// 			Status:       "Certificate does not exist.",
// 			Timestamp:    s,
// 			Tx:           Sha256(strconv.FormatInt(time.Now().UnixNano(), 10)),
// 			SerialNumber: serialNumber,
// 		}
// 		bData, _ := json.Marshal(resp)
// 		_, _ = fmt.Fprintf(w, string(bData))
// 	} else {
// 		delete(CertificateMap, serialNumber)
// 		delete(UserID, uid)
// 		s := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
// 		resp := RevokeResponse{
// 			Status:       "OK.",
// 			Timestamp:    s,
// 			Tx:           Sha256(strconv.FormatInt(time.Now().UnixNano(), 10)),
// 			SerialNumber: serialNumber,
// 		}
// 		bData, _ := json.Marshal(resp)
// 		_, _ = fmt.Fprintf(w, string(bData))
// 	}
// }

func Sha256(src string) string {
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

// 登录界面
// func login(w http.ResponseWriter, r *http.Request) {
// 	str, err := ioutil.ReadFile("./template/index.html")
// 	s, _ := os.Getwd()
// 	if err != nil {
// 		http.Error(w, s, 500)
// 		return
// 	}
// 	_, _ = w.Write([]byte(str))
// }

// // 通过用户 UID 获取证书
// func GetCertificateByUID(w http.ResponseWriter, r *http.Request) {
// 	_ = r.ParseForm()
// 	uid := r.Form.Get("uid")
// 	serialNumber := UserID[uid]

// 	gRWLock.RLock()
// 	defer gRWLock.RUnlock()

// 	if res, ok := CertificateMap[serialNumber]; !ok {
// 		http.Error(w, "Certificate does not exist.", 500)
// 	} else {
// 		bData, _ := json.Marshal(res)
// 		_, _ = fmt.Fprintf(w, string(bData))
// 	}
// }

func main() {
	http.HandleFunc("/ApplyForABSCertificate", ApplyForABSCertificate)
	http.HandleFunc("/VerifyABSCertificate", VerifyABSCertificate)
	http.HandleFunc("/VerifyABSCert", VerifyABSCert)

	http.HandleFunc("/RevokeABSCertificate", RevokeABSCertificate)

	http.HandleFunc("/GetCertificateNumber", GetCertificateNumber)
	http.HandleFunc("/GetCertificate", GetCertificate)
	http.HandleFunc("/GetCertificateFromFabric", GetCertificateFromFabric)

	//todo
	http.HandleFunc("/IoTDevTest", IoTDevTest)
	http.HandleFunc("/IotDevInit", IotDevInit)

	// http.HandleFunc("/RevokeABSCertificateByUID", RevokeABSCertificateByUID)
	// http.HandleFunc("/GetCertificateByUID", GetCertificateByUID)
	// http.HandleFunc("/login", login)
	// http.Handle("/", http.FileServer(http.Dir("template")))

	//http.HandleFunc("/ConcurrencyTest", ConcurrencyTest)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
		log.Fatalln(err)
	}
}
