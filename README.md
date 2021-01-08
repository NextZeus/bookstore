# Go Environments

```bash
# $HOME/.bashrc 
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
export GO111MODULE=on
export GOPROXY=https://goproxy.cn

# source .bashrc 生效配置
```

# install protoc dependency
[installing-protoc](http://google.github.io/proto-lens/installing-protoc.html)
```bash

PROTOC_ZIP=protoc-3.14.0-osx-x86_64.zip curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*' rm -f $PROTOC_ZIP

```

# install protobuf and proto and protoc-gen-go
```bash
brew tap grpc/grpc
brew install protobuf
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go

# protoc command global 否则会抛出异常 "protoc": executable file not found in $PATH
cp $HOME/go/bin/protoc-gen-go /usr/local/bin/

```

# ETCD (run in docker)
[running a single node etcd in container](https://etcd.io/docs/v3.4.0/op-guide/container/)
[quay.io/repository](https://quay.io/repository/coreos/etcd?tag=latest&tab=tags)
```bash
# 删除已有container
docker container rm /etcd

docker run -d -p 2379:2379 -p 2380:2380 \
  --volume=etcd-data:/etcd-data \
  --name etcd quay.io/coreos/etcd:latest \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name node1 \
  --initial-advertise-peer-urls http://192.168.48.143:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://192.168.48.143:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster node1=http://192.168.48.143:2380
 
```

# Mysql & Redis 
- 都是用的docker启动 没有本地安装
- mysql root@123456

# Test

```bash
xiaodong@localhost ~> sudo ulimit -n 2000
xiaodong@localhost ~> wrk -t10 -c1000 -d40s --latency "http://localhost:8888/check?book=go-zero"
Running 40s test @ http://localhost:8888/check?book=go-zero
  10 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    35.90ms   13.45ms 270.31ms   86.09%
    Req/Sec   672.71    584.46     2.60k    76.36%
  Latency Distribution
     50%   34.47ms
     75%   39.82ms
     90%   46.54ms
     99%   89.17ms
  267633 requests in 40.07s, 33.95MB read
  Socket errors: connect 759, read 47, write 0, timeout 0
Requests/sec:   6678.38
Transfer/sec:    867.41KB


```

# 启动顺序
1. etcd
2. rpc services
3. bookstore api server

