#!/usr/bin/env bash

source ./proto_dir.cfg

for ((i = 0; i < ${#all_proto[*]}; i++)); do
  proto=${all_proto[$i]}
  protoc -I ../../../  -I ./ --go_out=. --go-grpc_out=. $proto
  echo "protoc --go_out=plugins=grpc:." $proto
done
echo "proto file generate success"


j=0
for file in $(find ./crazy_server -name   "*.go"); do # Not recommended, will break on whitespace
    filelist[j]=$file
    j=`expr $j + 1`
done


for ((i = 0; i < ${#filelist[*]}; i++)); do
  proto=${filelist[$i]}
  cp $proto  ${proto#*./crazy_server/pkg/proto/}
done
rm crazy_server -rf
#find ./ -type f_packet_detail.sql -path "*.pb.go"|xargs sed -i 's/\".\/sdk_ws\"/\"crazy_server\/pkg\/proto\/sdk_ws\"/g'




