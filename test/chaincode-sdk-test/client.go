package main

import (
 	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
 	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
 	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 配置信息
var (
 sdk           *fabsdk.FabricSDK                              // Fabric SDK
 channelName   = "appchannel"                                 // 通道名称
 username      = "Admin"                                      // 用户
 chainCodeName = "dpki"                              // 链码名称
 endpoints     = []string{"peer0.CA1", "peer1.CA1"} // 要发送交易的节点
)

// init 初始化
func init() {
 	// 通过配置文件初始化SDK
	var err error
 	sdk, err = fabsdk.New(config.FromFile("config.yaml"))
 	if err != nil {
  		panic(err)
 	}
}

// ChannelExecute 区块链交互
func ChannelExecute(fcn string, args [][]byte) (string, error) {
 	// 创建客户端，表明在通道的身份
 	ctx := sdk.ChannelContext(channelName, fabsdk.WithUser(username))
 	cli, err := channel.New(ctx)
 	if err != nil { 
  		return "error", err
	}
 	// 对区块链账本的写操作（调用了链码的invoke）
 	resp, err := cli.Execute(channel.Request{
  		ChaincodeID: chainCodeName,
  		Fcn:         fcn,
  		Args:        args,
 	}, channel.WithTargetEndpoints(endpoints...))
 	if err != nil {
  		return "error", err
 	}
 	//返回链码执行后的结果
 	return string(resp.Payload), nil
}

// ChannelQuery 区块链查询
func ChannelQuery(fcn string, args [][]byte) (string, error) {
 	// 创建客户端，表明在通道的身份
 	ctx := sdk.ChannelContext(channelName, fabsdk.WithUser(username))
 	cli, err := channel.New(ctx)
 	if err != nil {
  		return "error", err
 	}
 	// 对区块链账本查询的操作（调用了链码的invoke），只返回结果
 	resp, err := cli.Query(channel.Request{
  		ChaincodeID: chainCodeName,
  		Fcn:         fcn,
  		Args:        args,
 	}, channel.WithTargetEndpoints(endpoints...))
 	if err != nil {
  		return "error", err
 	}
	//返回链码执行后的结果
 	return string(resp.Payload), nil
}