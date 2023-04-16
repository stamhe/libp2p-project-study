package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	libp2p "github.com/libp2p/go-libp2p"
)

func main() {
	// 默认监听所有可用的 ipv4, ipv6 网络接口，端口号随机
	// node, err := libp2p.New()

	// 指定监听的端口
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/8080"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Listen address: ", node.Addrs())

	// wait for SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	if err := node.Close(); err != nil {
		panic(err)
	}
}
