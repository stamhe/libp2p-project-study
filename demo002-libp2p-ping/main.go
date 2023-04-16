package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	libp2p "github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	ping "github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func main() {
	// 默认监听所有可用的 ipv4, ipv6 网络接口
	// node, err := libp2p.New()

	// 指定监听的端口
	// node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/8080"))

	// 不指定监听的端口, 使用随机的本地端口, 禁用默认的 ping 协议
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"), libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	fmt.Println("Listen address: ", node.Addrs())

	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("libp2p node address: ", addrs[0])

	// 如果带了连接参数，就连接指定的地址，并发发送 5 个 ping 消息
	if len(os.Args) > 1 {
		addr, _ := multiaddr.NewMultiaddr(os.Args[1])
		peer, _ := peerstore.AddrInfoFromP2pAddr(addr)

		if err := node.Connect(context.Background(), *peer); err != nil {
			panic(err)
		}

		fmt.Println("send 5 ping msg to ", addr)
		ch := pingService.Ping(context.Background(), peer.ID)
		for i := 0; i < 5; i++ {
			res := <-ch
			fmt.Println("pinged", addr, "in", res.RTT)
		}
	} else {
		// wait for SIGINT or SIGTERM signal
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Received signal, shutting down...")
	}

	if err := node.Close(); err != nil {
		panic(err)
	}
}
