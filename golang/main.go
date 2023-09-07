package main

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func setupDnsToRedis(redisClient *redis.Client) error {
	datas := [][2]string{
		{"google.com.vn", "192.168.0.1"},
		{"facebook.com", "192.168.0.2"},
		{"github.com", "192.168.0.3"},
		{"vnexpress.net", "192.168.0.4"},
	}

	for _, data := range datas {
		err := redisClient.Set(ctx, data[0], data[1], 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := setupDnsToRedis(rdb)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 10053,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	for {
		message := make([]byte, 1024)
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}
		go response(conn, remote, message[:rlen], rdb)
	}
}

func response(udpServer net.PacketConn, addr net.Addr, buf []byte, redisClient *redis.Client) {
	domain := string(buf)
	fmt.Printf("received: %s from %s\n", domain, addr)

	if ipAddress, err := redisClient.Get(ctx, domain).Result(); err != nil || err == redis.Nil {
		udpServer.WriteTo(nil, addr)
	} else {
		udpServer.WriteTo([]byte(ipAddress), addr)
	}
}
