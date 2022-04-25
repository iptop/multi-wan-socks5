package main

import (
	"fmt"
	"github.com/armon/go-socks5"
	"golang.org/x/net/context"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type InterfaceAddress struct {
	NAME string
	IPV4 string
}

func get_all_inter_face() []InterfaceAddress {
	var inters, _ = net.Interfaces()
	var arr = []InterfaceAddress{}
	for _, inter := range inters {
		if inter.MTU == -1 {
			continue
		}
		if inter.Flags&net.FlagUp != net.FlagUp {
			continue
		}
		var item InterfaceAddress
		item.NAME = inter.Name
		var ad, _ = inter.Addrs()
		if len(ad) < 2 {
			continue
		}
		var v4 = strings.Split(ad[1].String(), "/")[0]
		item.IPV4 = v4
		arr = append(arr, item)
	}
	return arr
}

func select_inter_face(inter_faces_list []InterfaceAddress) []string {
	var input string
	println("Input interface index For example, 1,2,3")
	_, _ = fmt.Scanf("%s", &input)
	var arr = strings.Split(input, ",")
	var rt = []string{}
	for _, str := range arr {
		str = strings.Trim(str, "")
		idx, _ := strconv.Atoi(str)
		rt = append(rt, inter_faces_list[idx].IPV4)
	}
	return rt
}

func s5(dial func(ctx context.Context, net_, addr string) (net.Conn, error)) {

	conf := &socks5.Config{
		Dial: dial,
	}
	server, _ := socks5.New(conf)

	// Create SOCKS5 proxy on localhost port 8000
	var _ = server.ListenAndServe("tcp", "127.0.0.1:8000")
}

func get_dial_func(local_addr_list []string) func(ctx context.Context, net_, addr string) (net.Conn, error) {

	var getLocalIp = func() string {
		rand.Seed(time.Now().UnixNano())
		var i = rand.Intn(len(local_addr_list))
		println(local_addr_list[i])
		return local_addr_list[i]
	}

	return func(ctx context.Context, net_, addr string) (net.Conn, error) {
		var rs = strings.Split(addr, ":")
		var localaddr net.TCPAddr
		var remoteaddr net.TCPAddr
		localaddr.IP = net.ParseIP(getLocalIp())
		remoteaddr.IP = net.ParseIP(rs[0])
		remoteaddr.Port, _ = strconv.Atoi(rs[1])
		return net.DialTCP(net_, &localaddr, &remoteaddr)
	}
}

func main() {
	var inter_faces_list = get_all_inter_face()
	for idx, inter := range inter_faces_list {
		println(idx, inter.NAME, inter.IPV4)
	}
	var local_addr_list = select_inter_face(inter_faces_list)
	if len(local_addr_list) == 0 {
		return
	}
	s5(get_dial_func(local_addr_list))
}
