package main

import (
	"fmt"
	"illustrated-tls/fakerand"
	tls "illustrated-tls/tlscopy"
	"net"
	"os"
	"time"
	//"github.com/syncsynchalt/illustrated-tls/fakerand"
	//tls "github.com/syncsynchalt/illustrated-tls/tlscopy"
)

var fakeRandData = []byte{
	0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
	0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
	0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
	0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
	0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
	0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
	0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf,
	0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf,
	0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf,
	0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xeb, 0xec, 0xed, 0xee, 0xef,
	0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff,
}

// KeyWriter is an io.Writer meant to print the NSS key log to stdout
type keyWriter struct {
	hasWritten bool
}

func (kw *keyWriter) Write(b []byte) (n int, err error) {
	if !kw.hasWritten {
		os.Stdout.Write([]byte("# key log data follows:\n"))
		kw.hasWritten = true
	}
	return os.Stdout.Write(b)
}

// a server that starts a TLS connection on port 8443, reads "ping", and responds "pong".
func main() {

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		panic(err)
	}

	rand := fakerand.New(fakeRandData)
	ln, err := tls.Listen("tcp", ":8443", &tls.Config{
		Rand: rand,
		Time: func() time.Time { return time.Unix(1538708249, 0) },
		CipherSuites: []uint16{
			// for the purpose of education we avoid AEAD cipher suites
			0xc013, // ECDHE-RSA-AES128-SHA
			0xc009, // ECDHE-ECDSA-AES128-SHA
			0xc014, // ECDHE-RSA-AES256-SHA
			0xc00a, // ECDHE-ECDSA-AES256-SHA
			0x002f, // RSA-AES128-SHA
			0x0035, // RSA-AES256-SHA
			0xc012, // ECDHE-RSA-3DES-EDE-SHA
			0x000a, // RSA-3DES-EDE-SHA
		},
		Certificates: []tls.Certificate{cert},
		KeyLogWriter: &keyWriter{},
	})
	if err != nil {
		panic(err)
	}

	count := 0
	for {
		fmt.Println(fmt.Sprintf("Server is listening...%d", count))
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
		count++
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 读取客户端数据
	rdata := make([]byte, 1024)
	n, err := conn.Read(rdata)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	fmt.Println("Server read data:", string(rdata[:n]))

	// 响应客户端
	wdata := []byte("pong")
	n, err = conn.Write(wdata)
	if n != len(wdata) {
		fmt.Printf("Incorrect write of %d (expected %d)\n", n, len(wdata))
		return
	}
	fmt.Println("Server wrote data:", string(wdata[:n]))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}

	// 尝试再次读取（这次应该不会成功）
	n, err = conn.Read(rdata)
	if n != 0 || err == nil {
		fmt.Println("Unexpected success on second read")
	} else {
		fmt.Println("Expected error on second read:", err)
	}
}
