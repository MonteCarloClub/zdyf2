#!/bin/bash

echo "区块链 ： 启动"
docker-compose up -d
echo "正在等待节点的启动完成，等待5秒"
sleep 5

CA1Peer0Cli="CORE_PEER_ADDRESS=peer0.CA1:7051 CORE_PEER_LOCALMSPID=CA1MSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/CA1/users/Admin@CA1/msp"
CA1Peer1Cli="CORE_PEER_ADDRESS=peer1.CA1:7051 CORE_PEER_LOCALMSPID=CA1MSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/CA1/users/Admin@CA1/msp"
CA2Peer0Cli="CORE_PEER_ADDRESS=peer0.CA2:7051 CORE_PEER_LOCALMSPID=CA2MSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/CA2/users/Admin@CA2/msp"
CA2Peer1Cli="CORE_PEER_ADDRESS=peer1.CA2:7051 CORE_PEER_LOCALMSPID=CA2MSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/CA2/users/Admin@CA2/msp"

echo "创建通道"
docker exec cli bash -c "$CA1Peer0Cli peer channel create -o orderer.com:7050 -c appchannel -f /etc/hyperledger/config/appchannel.tx"

echo "将所有节点加入通道"
docker exec cli bash -c "$CA1Peer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$CA1Peer1Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$CA2Peer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$CA2Peer1Cli peer channel join -b appchannel.block"

echo "更新锚节点"
docker exec cli bash -c "$CA1Peer0Cli peer channel update -o orderer.com:7050 -c appchannel -f /etc/hyperledger/config/CA1Anchor.tx"
docker exec cli bash -c "$CA2Peer0Cli peer channel update -o orderer.com:7050 -c appchannel -f /etc/hyperledger/config/CA2Anchor.tx"

echo "安装链码"
docker exec cli bash -c "$CA1Peer0Cli peer chaincode install -n dpki -v 1.0.3 -l golang -p chaincode"
docker exec cli bash -c "$CA1Peer1Cli peer chaincode install -n dpki -v 1.0.3 -l golang -p chaincode"
docker exec cli bash -c "$CA2Peer0Cli peer chaincode install -n dpki -v 1.0.3 -l golang -p chaincode"
docker exec cli bash -c "$CA2Peer1Cli peer chaincode install -n dpki -v 1.0.3 -l golang -p chaincode"

# 只需要其中一个节点实例化
# -n 对应上一步安装链码的名字
# -v 版本号
# -C 是通道，在fabric的世界，一个通道就是一条不同的链
# -c 为传参，传入init参数
echo "实例化链码"
docker exec cli bash -c "$CA1Peer0Cli peer chaincode instantiate -o orderer.com:7050 -C appchannel -n dpki -l golang -v 1.0.3 -c '{\"Args\":[]}' -P \"OR ('CA1MSP.member','CA2MSP.member')\""

echo "正在等待链码实例化完成，等待5秒"
sleep 5

# 进行链码交互，验证链码是否正确安装及区块链网络能否正常工作
echo "验证链码"
docker exec cli bash -c "$CA1Peer0Cli peer chaincode invoke -C appchannel -n dpki -c '{\"Args\":[\"setCertificate\",\"UUID1\",\"serialNumber1\",\"Certificate1\"]}'"
sleep 2
docker exec cli bash -c "$CA1Peer0Cli peer chaincode invoke -C appchannel -n dpki -c '{\"Args\":[\"getCertificate\",\"serialNumber1\"]}'"
docker exec cli bash -c "$CA1Peer0Cli peer chaincode invoke -C appchannel -n dpki -c '{\"Args\":[\"getCertificateByUUID\",\"UUID1\"]}'"