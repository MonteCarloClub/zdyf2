# ABS-DPKI
Decentralized Public Key Infrastructure Based on Attribute-Based Signature \
源自论文：https://doi.org/10.1360/SSI-2021-0177

## 编译 & 运行
需要有 Go 环境
### make
./build.sh build

### run CA & RA
./build.sh run 100 10\
100 代表 CA 数量，10 代表 RA 数量

### kill all processes
./build.sh stop

### make clean
./build.sh clean

## 使用
### IoT 设备初始化
发送请求：http://127.0.0.1:8001/IotDevInit

### 主页面
http://127.0.0.1:8001/login

### 证书管理页面
http://127.0.0.1:8001/login#/certificates

### 证书申请
http://127.0.0.1:8001/ApplyForABSCertificate?uid={}&&attribute={} \
也可以在页面操作

### 证书验证
http://127.0.0.1:8001/VerifyABSCertificate?no={serialNumber}

### 证书撤销
通过证书序列号：http://127.0.0.1:8001/RevokeABSCertificate?no={serialNumber} \
通过用户 UID：http://127.0.0.1:8001/RevokeABSCertificateByUID?userid={uid} \
也可以在页面操作

### 证书数量获取
http://127.0.0.1:8001/GetCertificateNumber

### 证书获取
通过证书序列号：http://127.0.0.1:8001/GetCertificate?no={serialNumber} \
通过用户 UID：http://127.0.0.1:8001/GetCertificateByUID?uid={uid} \
也可以在页面操作

## 测试
### 模拟多用户并发申请 & 验证证书
./Client/abs_client --n=10000 \
n 表示用户数量