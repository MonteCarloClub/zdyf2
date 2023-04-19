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

CleanBe() {
  cd CA && make clean && cd ..
  cd RA && make clean && cd ..
}

Build() {
  cd CA && make all && cd ..
  cd RA && make all && cd ..
  cd Test && make all && cd ..
  cd network && ./build.sh && cd ..
  sudo docker run -d --privileged=true --restart=always -p 6379:6379 -v ${PWD}/redis/conf/redis.conf:/etc/redis/redis.conf -v ${PWD}/redis/data:/data --name myredis redis redis-server /etc/redis/redis.conf --appendonly yes
  # docker run -d --privileged=true --restart=always -p 6379:6379 -v /home/user/zdyf2/redis/conf/redis.conf:/etc/redis/redis.conf -v /home/user/zdyf2/redis/data:/data --name myredis redis redis-server /etc/redis/redis.conf --appendonly yes
}

RunCA() {
  cd ./CA
  for i in $(seq $1 1 $2)
  do
    nohup ./$CA -port=$((9000+$i)) &
  done
  cd ..
}

RunRA() {
  cd ./RA
  for i in $(seq $1 1 $2)
  do
    nohup ./$RA -port=$((8000+$i)) -name=$(($i)) &
  done
  cd ..
}

RunAll() {
  cd network && ./start.sh && cd ..
  RunCA $1 $2
  RunRA $1 $2
}

RunBe(){
  RunCA $1 $2
  RunRA $1 $2
}

MakeBe(){
  cd CA && make all && cd ..
  cd RA && make all && cd ..
  cd Test && make all && cd ..
}


RunTest(){
  echo "申请及验证证书数量： $1" 
  cd ./Test
  ./$Test -n $1
  cd ..
  echo "总计生成证书数量：" 
  curl "http://10.176.40.47/dpki/GetCertificateNumber"
}

killBe(){
  pgrep abs_server_ | xargs kill -9
}

KillAll() {
  pgrep abs_server_ | xargs kill -9
  cd network && ./stop.sh && cd ..
}

if [ $1 == 'clean' ]; then Clean
elif [ $1 == 'build' ]; then Build
elif [ $1 == 'run' ]; then RunAll $2 $3
elif [ $1 == 'test' ]; then RunTest $2
elif [ $1 == 'runbe' ]; then RunBe $2 $3
elif [ $1 == 'killbe' ]; then killBe
elif [ $1 == 'makebe' ]; then MakeBe
elif [ $1 == 'cleanbe' ]; then CleanBe
elif [ $1 == 'stop' ]; then KillAll
else echo "unknown command"
fi