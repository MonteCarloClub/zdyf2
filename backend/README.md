# ABS-DPKI
Decentralized Public Key Infrastructure Based on Attribute-Based Signature \
源自论文：https://doi.org/10.1360/SSI-2021-0177

## 编译 & 运行
需要有 Go 环境，fabric环境, 配置Nginx负载均衡，启动docker环境
> systemctl start docker
### 参数修改
修改./RA/config.yaml 中的 zdfy2 路径为相应的完整路径
### make
./build.sh build

### run CA & RA
./build.sh run 100 10\
100 代表 CA 数量，10 代表 RA 数量，CA数量应为RA10倍

### kill all processes
./build.sh stop

### make clean
./build.sh clean

## 测试
### 模拟多用户并发申请 & 验证证书
./build.sh test

## 使用
路由可以改为{Nginx服务器}/dpki
当前Nginx为10.176.40.46
修改RA/ra.go redis服务器地址，默认为本机

### 证书申请
http://{Nginx服务器}/dpki/ApplyForABSCertificate?uid={}&&attribute={} \
也可以在页面操作

### 证书验证
证书序列号验证http://{Nginx服务器}/dpki/VerifyABSCertificate?no={serialNumber}
完整证书验证, post请求 http://{Nginx服务器}/dpki/VerifyABSCert cert={cert}

### 证书撤销
通过证书序列号：http:/{Nginx服务器}/dpki/RevokeABSCertificate?no={serialNumber} \

### 证书数量获取
http://{Nginx服务器}/dpki/GetCertificateNumber

### 证书获取
通过证书序列号：http:/{Nginx服务器}/dpki/GetCertificate?no={serialNumber} \
