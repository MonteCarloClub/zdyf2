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
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type FabricInfo struct {
	ABSUID       []byte
	SerialNumber []byte
	Cert         []byte
}

const (
	/**
	 * 一个遥远的时间戳：2106-01-02 15:04:05，Go百岁生日
	 * 令priority=DistantTimestamp-当前时间，即可实现在zset关于时间戳的降序排序
	 * todo 使此时间戳更大，并更新redis存储的值，避免千年虫问题
	 * @KofClubs 2023-04-05
	 */
	DistantTimestamp = 4291859045000000000.0
)

var (
	gPort      *int
	Version    string
	IssuerName string
	// CertificateMap map[string]CertificateResponse
	// UserID         map[string]string
	gRWLock          sync.RWMutex
	redisdb          *redis.Client
	redisMetaDb      *redis.Client
	redisRecordDb    *redis.Client
	redisBlacklistDb *redis.Client
	CAPort           int
	CANum            int
	RAbase           int
	validHour        time.Duration
	fbWorker         chan *FabricInfo
)

func init() {
	gPort = flag.Int("port", 8000, "ra server port.")
	Version = "1.0"
	raName := flag.Int("name", 0, "ra name")
	// CertificateMap = make(map[string]CertificateResponse)
	// UserID = make(map[string]string)
	flag.Parse()
	IssuerName = "CA-" + strconv.Itoa(*raName)
	RAbase = *raName
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "10.176.40.28:6379", // 指定
		Password: "",
		DB:       0,
	})
	redisMetaDb = redis.NewClient(&redis.Options{
		Addr:     "10.176.40.28:6379",
		Password: "",
		DB:       1,
	})
	redisRecordDb = redis.NewClient(&redis.Options{
		Addr:     "10.176.40.28:6379",
		Password: "",
		DB:       2,
	})
	redisBlacklistDb = redis.NewClient(&redis.Options{
		Addr:     "10.176.40.28:6379",
		Password: "",
		DB:       3,
	})

	// TODO redis 证书发放记录
	CAPort = 0      //CA的Port, 每调用一次CA， CAport + 1;
	CANum = 1       //每个RA对应的CA数量
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
	blacklist := getBlacklist()
	for _, v := range blacklist {
		if uid == v {
			log.Printf("[Apply] UID: %s is in blacklist.", uid)
			http.Error(w, "UID在黑名单中.", http.StatusBadRequest)
			return
		}
	}
	timeStr := time.Now()
	validTimeStr := time.Now().Add(time.Hour * validHour)
	serialNumber := uid + "-" + strconv.FormatInt(timeStr.UnixNano(), 10)
	// TODO 时间戳
	c := Certificate{
		Version:      Version,
		ABSUID:       uid,
		SerialNumber: serialNumber,
		Signature:    "Attribute-based Signature",
		// Issuer:       IssuerName,
		IssuerCA:       IssuerName,
		IssueTime:      timeStr.Format("2006-01-02 15:03:04"),
		ValidityPeriod: strconv.FormatInt(validTimeStr.UnixNano(), 10),
		ABSAttribute:   attribute,
	}
	b, _ := json.Marshal(c)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("http://127.0.0.1:"+strconv.Itoa(CAPort+9000+(RAbase-1)*CANum+1)+"/SingleGenerate", "application/json", bytes.NewReader(b))
	CAPort = (CAPort + 1) % CANum
	// c.IssuerCA = IssuerName + "-CA-" + strconv.Itoa(CAPort+9000+(RAbase-1)*CANum+1)
	if err != nil {
		log.Printf("[Apply] CA SingleGenerate: %s failed.", c.SerialNumber)
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	sign, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Apply] CA SingleGenerate read info: %s failed.", c.SerialNumber)
		http.Error(w, err.Error(), 500)
		return
	}
	res := CertificateResponse{
		CertificateContent: c,
		Hash:               Sha256(string(b)),
		ABSSign:            string(sign),
	}
	// 证书发放记录插入
	bData, _ := json.Marshal(res)
	_, _ = fmt.Fprintf(w, string(bData))
	redisdb.Set(serialNumber, string(bData), time.Hour*validHour)
	b, _ = json.Marshal(c)
	redisMetaDb.Set(serialNumber, string(b), time.Hour*validHour)

	// 证书颁发记录排序
	/**
	 * 此处在旧实现的基础上，将证书颁发信息存在redis的zset中，具体操作为：
	 * ZADD {IssuerName} {priority} {serialNumber}
	 * IssuerName - RA对应的全局唯一标识符
	 * priority - 优先级，zset对证书按priority升序排列，证书越新优先级越高，priority=大数-时间戳
	 * serialNumber - 证书对应的全局唯一标识符
	 * 需求指出，要求分页、从新到旧地获取任何RA对应颁发的证书，因此设计此关于时间戳降序的有序集合
	 * @KofClubs 2023-04-05
	 */
	zKey := IssuerName
	zValue := redis.Z{
		Score:  DistantTimestamp - float64(time.Now().UnixNano()),
		Member: serialNumber,
	}
	_, err = redisRecordDb.ZAdd(zKey, zValue).Result()
	_, err = redisRecordDb.ZAdd("Total", zValue).Result()
	if err != nil {
		fmt.Printf("[ERROR] fail to add certificate %v to meta zset", serialNumber)
	}
	log.Printf("[Apply] Apply success: %s.", c.SerialNumber)
	fbInfo := FabricInfo{
		ABSUID:       []byte(c.ABSUID),
		SerialNumber: []byte(c.SerialNumber),
		Cert:         []byte(b),
	}
	fbWorker <- &fbInfo

}

func ApplyForIllegal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = r.ParseForm()
	uid := r.Form.Get("uid")

	attribute := r.Form.Get("attribute")
	timeStr := time.Now().Add(time.Hour * validHour)
	validTimeStr := time.Now().Add(time.Hour * validHour)
	serialNumber := uid + "-" + strconv.FormatInt(timeStr.UnixNano(), 10)
	// TODO 时间戳
	c := Certificate{
		Version:      Version,
		SerialNumber: serialNumber,
		Signature:    "Attribute-based Signature",
		// Issuer:       IssuerName,
		IssuerCA:       IssuerName,
		IssueTime:      timeStr.Format("2006-01-02 15:03:04"),
		ValidityPeriod: validTimeStr.Format("2006-01-02 15:03:04"),
		ABSUID:         uid,
		ABSAttribute:   attribute,
	}
	b, _ := json.Marshal(c)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("http://127.0.0.1:"+strconv.Itoa(CAPort+9000+(RAbase-1)*CANum+1)+"/SingleGenerate", "application/json", bytes.NewReader(b))
	CAPort = (CAPort + 1) % CANum
	// c.IssuerCA = IssuerName + "-CA-" + strconv.Itoa(CAPort+9000+(RAbase-1)*CANum+1)
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
	// TODO 计算哈希
	res := CertificateResponse{
		CertificateContent: c,
		Hash:               Sha256(string(b)),
		ABSSign:            string(sign),
	}
	// 证书发放记录插入
	bData, _ := json.Marshal(res)
	_, _ = fmt.Fprintf(w, string(bData))
	redisdb.Set(serialNumber, string(bData), time.Hour*validHour)
	b, _ = json.Marshal(c)
	redisMetaDb.Set(serialNumber, string(b), time.Hour*validHour)

	// 证书颁发记录排序
	/**
	 * 此处在旧实现的基础上，将证书颁发信息存在redis的zset中，具体操作为：
	 * ZADD {IssuerName} {priority} {serialNumber}
	 * IssuerName - RA对应的全局唯一标识符
	 * priority - 优先级，zset对证书按priority升序排列，证书越新优先级越高，priority=大数-时间戳
	 * serialNumber - 证书对应的全局唯一标识符
	 * 需求指出，要求分页、从新到旧地获取任何RA对应颁发的证书，因此设计此关于时间戳降序的有序集合
	 * @KofClubs 2023-04-05
	 */
	zKey := IssuerName
	zValue := redis.Z{
		Score:  DistantTimestamp - float64(time.Now().UnixNano()),
		Member: serialNumber,
	}
	_, err = redisRecordDb.ZAdd(zKey, zValue).Result()
	_, err = redisRecordDb.ZAdd("Total", zValue).Result()
	if err != nil {
		fmt.Printf("[ERROR] fail to add certificate %v to meta zset", serialNumber)
	}

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
				// log.Printf("[Apply]Fabric setCertificate: %s failed %s", fbInfo.SerialNumber, err.Error())
				time.Sleep(time.Millisecond * 1000)
				fbWorker <- fbInfo
			} else {
				log.Printf("[Apply] Fabric setCertificate: %s success", fbInfo.SerialNumber)
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
		log.Println("[Verify] Read certificate failed:", err)
	}
	defer request.Body.Close()
	var cert CertificateResponse
	if err := json.Unmarshal([]byte(rawCert), &cert); err != nil {
		http.Error(writer, err.Error(), 500)
	} else {
		SNumber := cert.CertificateContent.SerialNumber
		b, _ := json.Marshal(cert.CertificateContent)
		requestHash := Sha256(string(b))
		rawData, err := redisdb.Get(SNumber).Result()
		var certFromRedis CertificateResponse
		err = json.Unmarshal([]byte(rawData), &certFromRedis)
		if err != nil {
			log.Println("[Verify] Redis get certificate failed:", SNumber)
			http.Error(writer, "Redis certificate failed.", 500)
			return
		}
		if requestHash == certFromRedis.Hash {
			valid := cert.CertificateContent.ValidityPeriod
			if valid < strconv.FormatInt(time.Now().UnixNano(), 10) {
				log.Println("[Verify] The certificate has expired:", SNumber)
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
			log.Println("[Verify] Complete Certificate verify:", SNumber)
			_, _ = fmt.Fprintf(writer, "True")

		} else {
			// log.Println(rawData)
			// log.Println(string(rawCert))
			log.Println("[Verify] Certificate verify failed,certificate is invalid:", SNumber)
			http.Error(writer, "The certificate is invalid.", 500)
			return
			// _, _ = fmt.Fprintf(writer, "False compare")
		}

		// if err != nil {
		// 	http.Error(writer, "The certificate is invalid.", 500)
		// 	return
		// }
		// log.Println("[Verify]Recieve Verify request:", SNumber)
		// // TODO 比对所有信息改为比对哈希
		// if rawData == string(rawCert) {
		// 	// _, _ = fmt.Fprintf(writer, "True compare")
		// 	valid := cert.CertificateContent.ValidityPeriod
		// 	if valid < strconv.FormatInt(time.Now().UnixNano(), 10) {
		// 		log.Println("[Verify]The certificate has expired:", SNumber)
		// 		http.Error(writer, "The certificate has expired.", 500)
		// 		return
		// 	}
		// 	sign := cert.ABSSign
		// 	client := &http.Client{Timeout: 10 * time.Second}
		// 	resp, err := client.Post("http://127.0.0.1:"+strconv.Itoa(CAPort+9000+(RAbase-1)*CANum)+"/SingleVerify", "application/json", bytes.NewReader([]byte(sign)))
		// 	CAPort = (CAPort + 1) % CANum
		// 	if err != nil {
		// 		http.Error(writer, err.Error(), 500)
		// 		return
		// 	}
		// 	defer resp.Body.Close()
		// 	log.Println("[Verify] Complete Certificate verify:", SNumber)
		// 	_, _ = fmt.Fprintf(writer, "True")

		// } else {
		// 	// log.Println(rawData)
		// 	// log.Println(string(rawCert))
		// 	log.Println("[Verify] Certificate verify failed,certificate is invalid:", SNumber)
		// 	http.Error(writer, "The certificate is invalid.", 500)
		// 	return
		// 	// _, _ = fmt.Fprintf(writer, "False compare")
		// }
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
		log.Printf("[Query] Query certificate %s is not found, it does not exist or has been revoked.", serialNumber)
		http.Error(w, "Certificate does not exist or has been revoked.", 500)
	} else {
		log.Printf("[Query] Query certificate %s success.", serialNumber)
		_, _ = fmt.Fprintf(w, rawData)
	}
	// if res, ok := CertificateMap[serialNumber]; !ok {
	// 	http.Error(w, "Certificate does not exist.", 500)
	// } else {
	// 	bData, _ := json.Marshal(res)
	// 	_, _ = fmt.Fprintf(w, string(bData))
	// }
}

func GetMetaCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = r.ParseForm()
	serialNumber := r.Form.Get("no")

	// gRWLock.RLock()
	// defer gRWLock.RUnlock()

	rawData, err := redisMetaDb.Get(serialNumber).Result()
	if err != nil {
		// log.Printf("[Query] Certificate %s is not found, it does not exist or has been revoked.", serialNumber)
		http.Error(w, "Certificate does not exist.", 500)
	} else {
		// log.Printf("[Query] Certificate %s success.", serialNumber)
		_, _ = fmt.Fprintf(w, rawData)
	}
}

// 撤销证书
func RevokeABSCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = r.ParseForm()
	serialNumber := r.Form.Get("no")
	_, err := redisdb.Del(serialNumber).Result()
	_, err = redisMetaDb.Del(serialNumber).Result()
	if err != nil {
		log.Printf("[Revoke] Certificate %s is not found, it does not exist or has been revoked.", serialNumber)
		http.Error(w, "Certificate does not exist or has been revoked.", 500)
	} else {
		log.Printf("[Revoke] Revoke certificate: %s success", serialNumber)
		_, _ = fmt.Fprintf(w, "Revoke OK.")
	}

	/**
	 * 此处对应地将证书颁发信息从redis的zset删除，具体操作为：
	 * ZREM {IssuerName} {serialNumber}
	 * IssuerName - RA对应的全局唯一标识符
	 * serialNumber - 证书对应的全局唯一标识符
	 * @KofClubs 2023-04-05
	 */
	zKey := IssuerName
	// _, err = redisdb.ZRem(zKey, serialNumber).Result()
	// if err != nil {
	// 	fmt.Printf("[ERROR] fail to remove certificate %v from zset", serialNumber)
	// }
	_, err = redisRecordDb.ZRem(zKey, serialNumber).Result()
	_, err = redisRecordDb.ZRem("Total", serialNumber).Result()
	if err != nil {
		fmt.Printf("[ERROR] fail to remove certificate %v from meta zset", serialNumber)
	}
	args := make([][]byte, 0)
	args = append(args, []byte(serialNumber))
	resp, err := ChannelExecute("getCertificate", args)
	if err != nil {
		return
	} else {
		bData, _ := json.Marshal(resp)
		fbInfo := FabricInfo{
			ABSUID:       []byte("[Revoked]"),
			SerialNumber: []byte(serialNumber),
			Cert:         []byte("[Revoked at " + time.Now().Format("2006-01-02 15:03:04") + "] " + string(bData)),
		}
		fbWorker <- &fbInfo
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
		keys, cursor, err = redisMetaDb.Scan(cursor, "*", 10).Result()
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
		// bData, _ := json.Marshal(resp)
		_, _ = fmt.Fprintf(w, string(resp))
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
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = redisMetaDb.Scan(cursor, "*", 10).Result()
		if err != nil {
			log.Printf("[IoTDevTest]获取证书列表失败: %s", err.Error())
			_, _ = fmt.Fprintf(w, err.Error(), 500)
			// panic(err)
			return
		}
		for _, v := range keys {
			rawData, err := redisMetaDb.Get(v).Result()
			if err != nil {
				log.Printf("[IoTDevTest]获取证书列表失败: %s", err.Error())
				_, _ = fmt.Fprintf(w, err.Error(), 500)
				return
			}
			a = append(a, rawData)
		}
		if cursor == 0 {
			break
		}
	}
	bData, _ := json.Marshal(a)
	log.Printf("[IoTDevTest]获取证书列表成功.")
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

// HandleGetScore 基于黑名单策略对其他CA打分
/**
 * 请求例子：/getScore?id=
 * 响应例子："score": "10"
 * @KofClubs 2023-04-05
 */
func HandleGetScore(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")

	type GetScoreResp struct {
		Score string `json:"score"`
	}
	q := req.URL.Query()
	id := q.Get("id")
	if id == "" {
		fmt.Printf("[ERROR] /getScore: id not found\n")
		http.Error(resp, "id not found", http.StatusBadRequest)
		return
	}
	serialNumbers, err := redisRecordDb.ZRange(id, 0, 10).Result()
	if err != nil {
		fmt.Printf("[ERROR] /getScore: fail to get certificates dealt by %v\n", id)
		http.Error(resp, "illegal id", http.StatusBadRequest)
		return
	}
	score := calculateScore(id, serialNumbers)
	respBody := &GetScoreResp{
		Score: strconv.Itoa(score),
	}
	respBodyBytes, err := json.Marshal(respBody)
	if err != nil {
		fmt.Printf("[ERROR] /getScore: fail to marshal response of %v\n", id)
		http.Error(resp, "fail to marshal", http.StatusInternalServerError)
	}
	resp.Write(respBodyBytes)
}

func HandleAddToBlacklist(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	q := req.URL.Query()
	uid := q.Get("uid")
	_, err := redisBlacklistDb.Get(uid).Result()
	if err != nil {
		redisBlacklistDb.Set(uid, "0", time.Hour*validHour)
		_, _ = fmt.Fprintf(resp, "Add to blacklist success.")
		return
	}
	_, _ = fmt.Fprintf(resp, uid+" is already illegal.")
	return
}

func HandleRemoveFromBlacklist(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	q := req.URL.Query()
	uid := q.Get("uid")
	_, err := redisBlacklistDb.Get(uid).Result()
	if err != nil {
		_, _ = fmt.Fprintf(resp, uid+" is not in blacklist.")
		return
	}
	_, err = redisBlacklistDb.Del(uid).Result()
	_, _ = fmt.Fprintf(resp, "Remove from blacklist success.")
	return
}

func HandleGetBlacklist(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	blacklist := getBlacklist()
	type blacklistResp struct {
		Blacklist []string `json:"certificates"`
	}
	respBody := &blacklistResp{
		Blacklist: blacklist,
	}
	respBosyBytes, err := json.Marshal(respBody)
	if err != nil {
		fmt.Printf("[ERROR] /GetBlacklist: fail to get blacklist")
		http.Error(resp, "fail to get blacklist", http.StatusBadRequest)
		return
	}
	resp.Write(respBosyBytes)
	return
}

// HandleGetCertificate 获取证书颁发记录，支持分页查询
/**
 * 请求字段：
 * index - 页码，从0开始，0表示最新的记录
 * count - 记录数目，例如(1, 10)表示秩	为10-19的10条记录
 *
 */
func HandleGetCertificates(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	type GetCertificateResp struct {
		Certificates []string `json:"certificates"`
	}

	q := req.URL.Query()
	index := q.Get("index")
	if index == "" {
		fmt.Printf("[ERROR] /getCertificates: index not found\n")
		http.Error(resp, "index not found", http.StatusBadRequest)
		return
	}
	indexInt64, err := strconv.ParseInt(index, 10, 64)
	if err != nil {
		fmt.Printf("[ERROR] /getCertificates: index not number\n")
		http.Error(resp, "index not number", http.StatusBadRequest)
		return
	}

	count := q.Get("count")
	if count == "" {
		fmt.Printf("[ERROR] /getCertificates: count not found\n")
		http.Error(resp, "count not found", http.StatusBadRequest)
		return
	}
	countInt64, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		fmt.Printf("[ERROR] /getCertificates: count not number\n")
		http.Error(resp, "count not number", http.StatusBadRequest)
		return
	}

	zKey := "Total"
	lo := countInt64 * indexInt64
	hi := countInt64 * (indexInt64 + 1)
	certificates, err := redisRecordDb.ZRange(zKey, lo, hi).Result()
	if err != nil {
		fmt.Printf("[ERROR] /getCertificates: fail to get certificates dealt by %v\n", zKey)
		http.Error(resp, "illegal id", http.StatusBadRequest)
		return
	}
	respBody := &GetCertificateResp{
		Certificates: certificates,
	}
	respBosyBytes, err := json.Marshal(respBody)
	if err != nil {
		fmt.Printf("[ERROR] /getCertificates: fail to marshal response of %v\n", zKey)
		http.Error(resp, "fail to marshal", http.StatusInternalServerError)
	}
	resp.Write(respBosyBytes)
}

func HandleGetCAName(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Write([]byte(IssuerName))
}

func main() {
	http.HandleFunc("/ApplyForABSCertificate", ApplyForABSCertificate)
	http.HandleFunc("/ApplyForIllegal", ApplyForIllegal)

	http.HandleFunc("/VerifyABSCertificate", VerifyABSCertificate)
	http.HandleFunc("/VerifyABSCert", VerifyABSCert)

	http.HandleFunc("/RevokeABSCertificate", RevokeABSCertificate)

	http.HandleFunc("/GetCertificateNumber", GetCertificateNumber)
	http.HandleFunc("/GetCertificate", GetCertificate)
	http.HandleFunc("/GetMetaCertificate", GetMetaCertificate)
	http.HandleFunc("/GetCertificateFromFabric", GetCertificateFromFabric)

	http.HandleFunc("/IoTDevTest", IoTDevTest)
	// http.HandleFunc("/IotDevInit", IotDevInit)

	http.HandleFunc("/getScore", HandleGetScore)
	http.HandleFunc("/addToBlacklist", HandleAddToBlacklist)
	http.HandleFunc("/removeFromBlacklist", HandleRemoveFromBlacklist)
	http.HandleFunc("/getBlacklist", HandleGetBlacklist)
	http.HandleFunc("/getCertificates", HandleGetCertificates)

	http.HandleFunc("/getCAName", HandleGetCAName)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
		log.Fatalln(err)
	}
}

func Sha256(src string) string {
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func calculateScore(id string, serialNumbers []string) int {
	score := 10
	blacklist := getBlacklist()
	for _, serialNumber := range serialNumbers {
		res := strings.Split(serialNumber, "-")
		if len(res) != 2 {
			return -1
		}
		for _, v := range blacklist {
			if res[0] == v {
				score = score - 1
			}
		}
	}
	return score
}

func getBlacklist() []string {
	var a []string
	var cursor uint64
	var keys []string
	for {
		var err error
		keys, cursor, err = redisBlacklistDb.Scan(cursor, "*", 10).Result()
		if err != nil {
			log.Printf("[getBlacklist]获取黑名单列表失败: %s", err.Error())
			return a
		}
		a = append(a, keys...)
		if cursor == 0 {
			break
		}
	}
	return a
}
