# goperf

## Description
tcp并发，吞吐量测试工具（可测试socks5代理）

## 准备工作
~~~
# Download this repo
mkdir -p ./git/src/gitee.com/gbat
cd ./git/src/gitee.com/gbat
git clone https://gitee.com/gbat/goperf.git
cd goperf

# Build benchmark tools
cd receiver
go build

cd send
go build

~~~


## 吞吐量测试

### receiver
~~~
recerver -ip=127.0.0.1 -port=10001 -k=false -test=1 -b=500

ip:tcp服务端绑定的IP
port:tcp服务端绑定的端口号
k:开启长连接，默认开启
b:每次发送数据大小，默认500，单位：B
test：0-并发测试，1-吞吐量测试
~~~

### send

####直连
~~~
send -ip=127.0.0.1 -port=10000 -k=false -test=1 -b=500 -amount=1 -concurrency=1 -type=direct -destips=127.0.0.1 -destport=10001

ip:tcp客户端绑定的IP
port:tcp客户端绑定的端口号
k:开启长连接，默认开启
b:每次发送数据大小，默认500，单位：B
test：0-并发测试，1-吞吐量测试
amount：测试发送的数据总量，单位GB
concurrency：并发数
type：连接类型 direct-直连，socks-socks5代理
destips：目标ip，即目标tcp服务端ip地址集合，可为多个tcp节点，以","分割
destport：目标端口号：即目标tcp服务端端口号

~~~

####代理
~~~
send -ip=127.0.0.1 -port=10000 -k=false -test=1 -b=500 -amount=1 -concurrency=1 -type=socks -destips=127.0.0.1 -destport=10001 -proxyip=127.0.0.1 -proxyports=11000,12000 -porttype=1
         -user=admin -pwd=123456

ip:tcp客户端绑定的IP
port:tcp客户端绑定的端口号
k:开启长连接，默认开启
b:每次发送数据大小，默认500，单位：B
test：0-并发测试，1-吞吐量测试
amount：测试发送的数据总量，单位GB
concurrency：并发数
type：连接类型 direct-直连，socks-socks5代理
destips：目标ip，即目标tcp服务端ip地址集合，可为多个tcp节点，以","分割
destport：目标端口号：即目标tcp服务端端口号
proxyip:socks5代理ip地址
proxyports:socks5代理端口号，可以开启多个端口号转发，以","分割
porttype：socks5代理端口号类型，0-单端口，1-多端口
user:socks5代理用户名
pwd:socks5代理密码
~~~

## 并发测试

### receiver
~~~
recerver -ip=127.0.0.1 -port=10001 -k=false -test=0 -b=500

ip:tcp服务端绑定的IP
port:tcp服务端绑定的端口号
k:开启长连接，默认开启
b:每次发送数据大小，默认500，单位：B
test：0-并发测试，1-吞吐量测试
~~~
### send

####直连
~~~
send -ip=127.0.0.1 -port=10000 -k=false -test=0 -b=500 -amount=1 -concurrency=100 -type=direct -destips=127.0.0.1 -destport=10001

ip:tcp客户端绑定的IP
port:tcp客户端绑定的端口号
k:开启长连接，默认开启
b:每次发送数据大小，默认500，单位：B
test：0-并发测试，1-吞吐量测试
amount：测试发送的数据总量，单位GB
concurrency：并发数
type：连接类型 direct-直连，socks-socks5代理
destips：目标ip，即目标tcp服务端ip地址集合，可为多个tcp节点，以","分割
destport：目标端口号：即目标tcp服务端端口号

~~~

####代理
~~~
send -ip=127.0.0.1 -port=10000 -k=false -test=0 -b=500 -amount=1 -concurrency=100 -type=socks -destips=127.0.0.1 -destport=10001 -proxyip=127.0.0.1 -proxyports=11000,12000 -porttype=1
         -user=admin -pwd=123456

ip:tcp客户端绑定的IP
port:tcp客户端绑定的端口号
k:开启长连接，默认开启
b:每次发送数据大小，默认500，单位：B
test：0-并发测试，1-吞吐量测试
amount：测试发送的数据总量，单位GB
concurrency：并发数
type：连接类型 direct-直连，socks-socks5代理
destips：目标ip，即目标tcp服务端ip地址集合，可为多个tcp节点，以","分割
destport：目标端口号：即目标tcp服务端端口号
proxyip:socks5代理ip地址
proxyports:socks5代理端口号，可以开启多个端口号转发，以","分割
porttype：socks5代理端口号类型，0-单端口，1-多端口
user:socks5代理用户名
pwd:socks5代理密码
~~~
###注意
test=1，即测试吞吐量时，concurrency尽量不要太大，默认值即可

###效果

![image](https://gitee.com/gbat/goperf/blob/master/static/concurrency.png)

![image](https://gitee.com/gbat/goperf/blob/master/static/iops.png)
