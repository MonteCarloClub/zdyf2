# ABS-DPKI
Decentralized Public Key Infrastructure Based on Attribute-Based Signature \
源自论文：https://doi.org/10.1360/SSI-2021-0177

## 编译 & 运行
需要有 Go 环境，fabric环境, 配置Nginx负载均衡，启动docker环境
> systemctl start docker
### 参数修改
修改./RA/config.yaml 中的 zdfy2 路径为相应的完整路径
#### make
./build.sh build

#### run CA & RA

服务器1运行./build.sh run 1 5\

1 代表 CA 起始序号，5代表 CA 截止序号

服务器2运行./build.sh run 6 10\

#### kill all processes
./build.sh stop

#### make clean
./build.sh clean

#### 仅控制后端

``` 
./build.sh makebe

./build.sh runbe 1 5

./build.sh killbe

./build.sh cleanbe
```



## 测试
### 模拟多用户并发申请 & 验证证书
./build.sh test

## 使用
路由可以改为{Nginx服务器}/dpki
当前Nginx为10.176.40.47
修改RA/ra.go redis服务器地址，默认为本机

### 证书申请

uid - 用户名

http://{Nginx服务器}/dpki/ApplyForABSCertificate?uid={}

### 非法申请

uid - 用户名

http://{Nginx服务器}/dpki/ApplyForIllegal?uid={}

### 证书验证

no - 证书serialNumber

cert - 完整证书文件

证书序列号验证http://{Nginx服务器}/dpki/VerifyABSCertificate?no={serialNumber}
完整证书验证, post请求 http://{Nginx服务器}/dpki/VerifyABSCert cert={cert}

### 证书撤销

no - 证书serialNumber

http:/{Nginx服务器}/dpki/RevokeABSCertificate?no={serialNumber} 

### 证书数量获取
http://{Nginx服务器}/dpki/GetCertificateNumber

### 证书获取

no - 证书serialNumber

http:/{Nginx服务器}/dpki/GetCertificate?no={serialNumber} \

### 简要证书获取

no - 证书serialNumber

http:/{Nginx服务器}/dpki/GetMetaCertificate?no={serialNumber} \

### 证书从区块链获取

no - 证书serialNumber

http:/{Nginx服务器}/dpki/GetCertificateFromFabric?no={serialNumber} \

### 证书列表

http:/{Nginx服务器}/dpki/IoTDevTest

### CA信誉审计

id - CA名称

http:/{Nginx服务器}/dpki/getScore?id=CA-1 \

### 黑名单操作：加入、移出、获取黑名单

uid - 用户名

http:/{Nginx服务器}/dpki/addToBlacklist?uid=

http:/{Nginx服务器}/dpki/removeFromBlacklist?uid=

http:/{Nginx服务器}/dpki/getBlacklist

### 证书发放记录

获取证书颁发记录，支持分页查询 

index - 页码，从0开始，0表示最新的记录

count - 记录数目，例如(10, 10)表示秩为10-19的10条记录

http:/{Nginx服务器}/dpki/getCertificates?index=0&count=10

