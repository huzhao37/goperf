/**
 * @Author: hiram
 * @Date: 2020/7/9 9:20
 */
package main

import (
	"flag"
	"fmt"
	"gitee.com/gbat/goperf"
	"os"
	"os/signal"
	"syscall"
)

var (
	server = flag.Bool("s", false, "default:start on tcp client ,select true for start on tcp server.")
)

func main() {
	flag.Parse()
	if *server {
		goperf.Server()
	} else {
		goperf.Client()
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
