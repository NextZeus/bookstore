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

## 操作命令
```bash
mkdir bookstore
cd bookstore

#初始化go.mod
go mod init bookstore 

#创建api目录
mkdir api 

#生成api文件
goctl api -o bookstore.api 

#使用goctl生成API Gateway代码
goctl api go -api bookstore.api -dir .

#在 api 目录下启动API Gateway服务，默认侦听在8888端口
go run bookstore.go -f etc/bookstore-api.yaml

# 测试API Gateway服务
curl -i "http://localhost:8888/check?book=go-zero"


# rpc service

# 在bookstore目录下创建rpc目录
mkdir rpc

# 在rpc/add目录下编写add.proto文件
#可以通过命令生成proto文件模板
goctl rpc template -o add.proto

#用goctl生成rpc代码，在rpc/add目录下执行命令
goctl rpc proto -src add.proto -dir .

#运行add rpc service
go run add.go -f etc/add.yaml

```

## API Gateway调用rpc service
- 修改配置bookstore-api.yaml
- 通过etcd自动去发现可用的add/check服务
```yaml
Add:
  Etcd:
    Hosts:
      - localhost:2379
    Key: add.rpc
Check:
  Etcd:
    Hosts:
      - localhost:2379
    Key: check.rpc
```
### 增加服务依赖
- internal/config/config.go
```go
type Config struct {
    rest.RestConf
    Add   zrpc.RpcClientConf     // 手动代码
    Check zrpc.RpcClientConf     // 手动代码
}
```
- internal/svc/servicecontext.go
```go
type ServiceContext struct {
    Config  config.Config
    Adder   adder.Adder          // 手动代码
    Checker checker.Checker      // 手动代码
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:  c,
        Adder:   adder.NewAdder(zrpc.MustNewClient(c.Add)),         // 手动代码
        Checker: checker.NewChecker(zrpc.MustNewClient(c.Check)),   // 手动代码
    }
}

```

## 问题
#### grpc 我用的是v1.29.1版本，升级到v1.35.0后抛出异常
[github issue](https://github.com/etcd-io/etcd/issues/12124)
```bash
localhost   ~/go/bookstore/rpc   master ●✚  go run add.go -f etc/add.yaml
# go.etcd.io/etcd/clientv3/balancer/picker
../../pkg/mod/go.etcd.io/etcd@v0.0.0-20200402134248-51bdeb39e698/clientv3/balancer/picker/err.go:25:9: cannot use &errPicker literal (type *errPicker) as type Picker in return argument:
	*errPicker does not implement Picker (wrong type for Pick method)
		have Pick(context.Context, balancer.PickInfo) (balancer.SubConn, func(balancer.DoneInfo), error)
		want Pick(balancer.PickInfo) (balancer.PickResult, error)
../../pkg/mod/go.etcd.io/etcd@v0.0.0-20200402134248-51bdeb39e698/clientv3/balancer/picker/roundrobin_balanced.go:33:9: cannot use &rrBalanced literal (type *rrBalanced) as type Picker in return argument:
	*rrBalanced does not implement Picker (wrong type for Pick method)
		have Pick(context.Context, balancer.PickInfo) (balancer.SubConn, func(balancer.DoneInfo), error)
		want Pick(balancer.PickInfo) (balancer.PickResult, error)
# github.com/tal-tech/go-zero/zrpc/internal/balancer/p2c
../../pkg/mod/github.com/tal-tech/go-zero@v1.1.2/zrpc/internal/balancer/p2c/p2c.go:41:32: not enough arguments in call to base.NewBalancerBuilder
	have (string, *p2cPickerBuilder)
	want (string, base.PickerBuilder, base.Config)
../../pkg/mod/github.com/tal-tech/go-zero@v1.1.2/zrpc/internal/balancer/p2c/p2c.go:58:9: cannot use &p2cPicker literal (type *p2cPicker) as type balancer.Picker in return argument:
	*p2cPicker does not implement balancer.Picker (wrong type for Pick method)
		have Pick(context.Context, balancer.PickInfo) (balancer.SubConn, func(balancer.DoneInfo), error)
		want Pick(balancer.PickInfo) (balancer.PickResult, error)


```

#### etcd timeout超时

```bash

localhost  ~/go/bookstore/rpc   master ●✚  go run add.go -f etc/add.yaml
Starting rpc server at 127.0.0.1:8080...
{"level":"warn","ts":"2021-02-04T16:04:10.279+0800","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-5dd23b6b-d28b-49ab-a5b7-619860b6ff98/127.0.0.1:2379","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = latest balancer error: all SubConns are in TransientFailure, latest connection error: connection error: desc = \"transport: Error while dialing dial tcp 192.168.48.143:2379: i/o timeout\""}

#原因：第一次运行成功时网络IP是 192.168.48.143, 这次切换了网络，IP变了 ，就超时了。切换到之前的网络就又好了

```