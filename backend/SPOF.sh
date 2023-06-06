#!/bin/bash

Test=dpki_test

RunTest(){
  echo "SPOF测试： 申请100个证书" 
  cd ./Test
  ./$Test -m SPOF
  cd ..
}

RunTest