#!/bin/bash

CA=abs_server_ca
RA=abs_server_ra

Clean() {
  cd CA && make clean && cd ..
  cd RA && make clean && cd ..
  cd Client && make clean && cd ..
}

Build() {
  cd CA && make all && cd ..
  cd RA && make all && cd ..
  cd Client && make all && cd ..
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
  RunCA $1
  RunRA $2
}

KillAll() {
  pgrep abs_server_ | xargs kill -9
}

if [ $1 == 'clean' ]; then Clean
elif [ $1 == 'build' ]; then Build
elif [ $1 == 'run' ]; then RunAll $2 $3
elif [ $1 == 'stop' ]; then KillAll
else echo "unknown command"
fi