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
cd cmd
go build -o goperf


~~~


## 吞吐量测试

### server
~~~
goperf -s -ip=127.0.0.1 -port=10001  -io -b=500

s:开启服务端
ip:tcp服务端绑定的IP
port:tcp服务端绑定的端口号
k:开启长连接，默认关闭
b:每次发送数据大小，默认500，单位：B
io：默认-并发测试，io-吞吐量测试
~~~

### client

####直连
~~~
goperf -ip=127.0.0.1 -port=10000  -io -b=500 -amount=1 -concurrency=1  -destips=127.0.0.1 -destport=10001

ip:tcp客户端绑定的IP
port:tcp客户端绑定的端口号
k:开启长连接，默认关闭
b:每次发送数据大小，默认500，单位：B
io：默认-并发测试，io-吞吐量测试
amount：测试发送的数据总量，单位GB
concurrency：并发数
socks：连接类型 默认-直连，socks-socks5代理
destips：目标ip，即目标tcp服务端ip地址集合，可为多个tcp节点，以","分割
destport：目标端口号：即目标tcp服务端端口号

~~~

####代理
~~~
goperf -ip=127.0.0.1 -port=10000  -io  -b=500 -amount=1 -concurrency=1 -socks -destips=127.0.0.1 -destport=10001 -proxyip=127.0.0.1 -proxyports=11000,12000 -multiport
         -user=admin -pwd=123456

ip:tcp客户端绑定的IP
port:tcp客户端绑定的端口号
k:开启长连接，默认关闭
b:每次发送数据大小，默认500，单位：B
io：默认-并发测试，io-吞吐量测试
amount：测试发送的数据总量，单位GB
concurrency：并发数
socks：连接类型 默认-直连，socks-socks5代理
destips：目标ip，即目标tcp服务端ip地址集合，可为多个tcp节点，以","分割
destport：目标端口号：即目标tcp服务端端口号
proxyip:socks5代理ip地址
proxyports:socks5代理端口号，可以开启多个端口号转发，以","分割
multiport：socks5代理端口号类型，默认-单端口，multiport-多端口
user:socks5代理用户名
pwd:socks5代理密码
~~~

## 并发测试

### server
~~~
goperf -s -ip=127.0.0.1 -port=10001  -b=500

s:开启服务端
ip:tcp服务端绑定的IP
port:tcp服务端绑定的端口号
k:开启长连接，默认关闭
b:每次发送数据大小，默认500，单位：B
io：默认-并发测试，io-吞吐量测试
~~~
### client

####直连
~~~
goperf -ip=127.0.0.1 -port=10000   -b=500 -amount=1 -concurrency=100  -destips=127.0.0.1 -destport=10001

ip:tcp客户端绑定的IP
port:tcp客户端绑定的端口号
k:开启长连接，默认关闭
b:每次发送数据大小，默认500，单位：B
io：默认-并发测试，io-吞吐量测试
amount：测试发送的数据总量，单位GB
concurrency：并发数
socks：连接类型 默认-直连，socks-socks5代理
destips：目标ip，即目标tcp服务端ip地址集合，可为多个tcp节点，以","分割
destport：目标端口号：即目标tcp服务端端口号

~~~

####代理
~~~
goperf -ip=127.0.0.1 -port=10000  -b=500 -amount=1 -concurrency=100 -socks -destips=127.0.0.1 -destport=10001 -proxyip=127.0.0.1 -proxyports=11000,12000 -multiport
         -user=admin -pwd=123456

ip:tcp客户端绑定的IP
port:tcp客户端绑定的端口号
k:开启长连接，默认关闭
b:每次发送数据大小，默认500，单位：B
io：默认-并发测试，io-吞吐量测试
amount：测试发送的数据总量，单位GB
concurrency：并发数
socks：连接类型 默认-直连，socks-socks5代理
destips：目标ip，即目标tcp服务端ip地址集合，可为多个tcp节点，以","分割
destport：目标端口号：即目标tcp服务端端口号
proxyip:socks5代理ip地址
proxyports:socks5代理端口号，可以开启多个端口号转发，以","分割
multiport：socks5代理端口号类型，默认-单端口，multiport-多端口
user:socks5代理用户名
pwd:socks5代理密码
~~~
###注意
-io,即测试吞吐量时，concurrency尽量不要太大，默认值即可

###效果

![并发测试](https://images.gitee.com/uploads/images/2020/0709/100931_ba66c857_671199.png "concurrency.png")
![吞吐量测试](https://images.gitee.com/uploads/images/2020/0709/100949_e130e667_671199.png "iops.png")

