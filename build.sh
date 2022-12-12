#!/bin/bash

CA=abs_server_ca
RA=abs_server_ra
Test=dpki_test
Clean() {
  cd CA && make clean && cd ..
  cd RA && make clean && cd ..
  cd Test && make clean && cd ..
  cd network && ./clean.sh && cd ..
  cd redis/data && rm -rf * && docker rm -f myredis && cd ../..

}

Build() {
  cd CA && make all && cd ..
  cd RA && make all && cd ..
  cd Test && make all && cd ..
  cd network && ./build.sh && cd ..
  docker run -d --privileged=true --restart=always -p 6379:6379 -v ${PWD}/redis/conf/redis.conf:/etc/redis/redis.conf -v ${PWD}/redis/data:/data --name myredis redis redis-server /etc/redis/redis.conf --appendonly yes
}

RunCA() {
  cd ./CA
  for i in $(seq 1 1 $1)
  do
    nohup ./$CA -port=$((9000+$i)) &
  done
  cd ..
}

RunRA() {
  cd ./RA
  for i in $(seq 1 1 $1)
  do
    nohup ./$RA -port=$((8000+$i)) -name=$(($i)) &
  done
  cd ..
}

RunAll() {
  cd network && ./start.sh && cd ..
  RunCA $1
  RunRA $2
}

RunTest(){
  echo "申请及验证证书数量： 1000" 
  cd ./Test
  ./$Test -n 1000
  cd ..
  echo "总计生成证书数量：" 
  curl "http://10.176.40.46/dpki/GetCertificateNumber"
}

KillAll() {
  pgrep abs_server_ | xargs kill -9
  cd network && ./stop.sh && cd ..
}

if [ $1 == 'clean' ]; then Clean
elif [ $1 == 'build' ]; then Build
elif [ $1 == 'run' ]; then RunAll $2 $3
elif [ $1 == 'test' ]; then RunTest
elif [ $1 == 'stop' ]; then KillAll
else echo "unknown command"
fi