package main

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/big"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/proxy"
)

var (
	fProxyPorts  = flag.String("proxyports", "10000", "Proxy ports array for sending traffic to(10000,10001,10002).")
	fDestPort    = flag.Int("destport", 10001, "Dest port.")
	fType        = flag.String("type", "direct", "Proxy type of the target, either 'direct' or 'socks'.")
	fAmount      = flag.Int("amount", 1, "Amount of traffic to send, in GB.")
	fConcurrency = flag.Int("concurrency", 1, "Number of concurrect connections for benchmark.")

	fProxyIp      = flag.String("proxyip", "10.2.155.242", "Proxy IP to listen on.")
	fDestIps      = flag.String("destips", "127.0.0.1", "Dest IPs to Access.")
	testType      = flag.Int("test", 0, "test type,ex:0 is maxcon,1 is IOPS.")
	proxyPortType = flag.Int("porttype", 0, "test type,ex:0 is single proxy port,1 is multi proxy port.")
	keepAlive     = flag.Bool("k", true, "is keep alive")

	flocalIp   = flag.String("ip", "127.0.0.1", "local IP to bind.")
	flocalPort = flag.String("port", "9980", "local PORT to bind.")

	buffer = flag.Int64("b", 500, "send bytes,unit: B.")

	user = flag.String("user", "admin", "socks5 proxy user.")
	pwd  = flag.String("pwd", "123456", "socks5 proxy pwd.")
)

func makeConnection() ([]net.Conn, error) {
	destIps := strings.Split(*fDestIps, ",")
	conns := make([]net.Conn, 0, 3)
	switch *fType {
	case "direct":
		lAddr := &net.TCPAddr{}
		var err error
		//本地地址  ipaddr是本地外网IP
		if *testType == 0 {
			lAddr = nil
		} else {
			lAddr, err = net.ResolveTCPAddr("tcp4", *flocalIp+":"+*flocalPort)
			if err != nil {
				panic(err)
			}
		}

		//destIp:="127.0.0.1"
		for _, item := range destIps {
			destIp := item
			conn, err := net.DialTCP("tcp4", lAddr, &net.TCPAddr{
				IP:   net.ParseIP(destIp),
				Port: *fDestPort,
			})
			if err != nil {
				panic(err)
			}
			conns = append(conns, conn)
		}
		return conns, nil
	case "socks":
		var err error
		proxyPort := 10000
		proxyPorts := strings.Split(*fProxyPorts, ",")
		if *proxyPortType == 0 {
			proxyPort, err = strconv.Atoi(proxyPorts[0])
			if err != nil {
				panic(err)
			}
		}
		if *proxyPortType == 1 {
			min, err := strconv.Atoi(proxyPorts[0])
			if err != nil {
				panic(err)
			}
			max, err := strconv.Atoi(proxyPorts[len(proxyPorts)-1])
			if err != nil {
				panic(err)
			}
			proxyPort = int(RangeRand(int64(min), int64(max)))
		}
		//port:=int(RangeRand(11000,11004))
		dialer, err := proxy.SOCKS5("tcp4", fmt.Sprintf("%s:%d", *fProxyIp, proxyPort), &proxy.Auth{User: *user, Password: *pwd}, proxy.Direct)
		if err != nil {
			return nil, err
		}
		for _, item := range destIps {
			destIp := item
			conn, err := dialer.Dial("tcp4", fmt.Sprintf("%s:%d", destIp, *fDestPort))
			if err != nil {
				panic(err)
			}
			conns = append(conns, conn)
		}
		return conns, nil
	default:
		return nil, errors.New("Unknown proxy type: " + *fType)
	}
}

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

func main() {
	flag.Parse()

	if *testType == 0 {
		Concurrency()
	}
	if *testType == 1 {
		IOPS()
	}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

//并发
func Concurrency() {
	var wg sync.WaitGroup
	var threeWg sync.WaitGroup
	startTime := time.Now().Unix()
	for i := 0; i < *fConcurrency; i++ {
		wg.Add(1)
		go func() {
			buf := make([]byte, *buffer)
			rand.Read(buf)

			conns, err := makeConnection()
			if err != nil {
				panic(err)
			}
			for _, item := range conns {
				conn := item
				threeWg.Add(1)
				go func() {
					var connWg sync.WaitGroup
					connWg.Add(2)
					go func() {

						_, err := conn.Write(buf)
						if err != nil {
							panic(err)
						}

						connWg.Done()
					}()
					go func() {
						totalBytes := *buffer
						for {
							var count uint64
							if err := binary.Read(conn, binary.BigEndian, &count); err != nil {
								panic(err)
							}
							if count >= uint64(totalBytes) {
								break
							}
						}
						connWg.Done()
					}()
					connWg.Wait()
					if !*keepAlive {
						conn.Close()
					}
					threeWg.Done()
				}()
			}
			threeWg.Wait()
			wg.Done()
		}()
	}
	wg.Wait()

	endTime := time.Now().Unix()
	elapsed := endTime - startTime
	if elapsed == 0 {
		fmt.Println("Finished in 0 second. Too fast for benchmark.")
		return
	}
	speed := int64(*fConcurrency) / elapsed
	fmt.Println("send:", *buffer, "B of data sent through", *fConcurrency, "connections in", elapsed, "seconds, with speed", speed, "op/s.")
}

//吞吐量
func IOPS() {

	const BufSize = 128 * 1024
	var wg sync.WaitGroup

	startTime := time.Now().Unix()
	for i := 0; i < *fConcurrency; i++ {
		wg.Add(1)
		go func() {
			buf := make([]byte, BufSize)
			rand.Read(buf)
			conns, err := makeConnection()
			if err != nil {
				panic(err)
			}
			var connWg sync.WaitGroup
			connWg.Add(2)
			go func() {
				totalBytes := int64(*fAmount) * 1024 * 1024 * 1024
				for totalBytes > 0 {
					_, err := conns[0].Write(buf)
					if err != nil {
						panic(err)
					}
					totalBytes -= BufSize
				}
				connWg.Done()
			}()
			go func() {
				totalBytes := int64(*fAmount) * 1024 * 1024 * 1024
				for {
					var count uint64
					if err := binary.Read(conns[0], binary.BigEndian, &count); err != nil {
						panic(err)
					}
					if count >= uint64(totalBytes) {
						break
					}
				}
				connWg.Done()
			}()
			connWg.Wait()
			if !*keepAlive {
				conns[0].Close()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	endTime := time.Now().Unix()
	elapsed := endTime - startTime
	if elapsed == 0 {
		fmt.Println("Finished in 0 second. Too fast for benchmark.")
		return
	}
	dataAmount := uint64(*fConcurrency) * uint64(*fAmount)

	speed := dataAmount * 1024 / uint64(elapsed)
	fmt.Println("send:", dataAmount, "GB of data sent through", *fConcurrency, "connections in", elapsed, "seconds, with speed", speed, "MB/s.")
}
