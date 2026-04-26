package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

const max_clients = 100

var clientSem = Init(max_clients)

func handleClient(clientFD int) {
	clientSem.Wait()
	defer clientSem.Signal()
	defer syscall.Close(clientFD)

	buf := make([]byte, 1024)

	n, err := syscall.Read(clientFD, buf)
	if err != nil {
		fmt.Println("read error ", err)
		return
	}

	fmt.Println(string(buf[:n]))
	syscall.Write(clientFD, []byte("hello from the server"))
}

func RunMultiThreadedServer(port int) {
	proxyFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("socket creation failed ", err)
	}

	err = syscall.SetsockoptInt(proxyFd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	if err != nil {
		fmt.Println("setup error ", err)
	}

	sockAddr := &syscall.SockaddrInet4{Port: port}
	copy(sockAddr.Addr[:], net.ParseIP("0.0.0.0").To4())

	if err = syscall.Bind(proxyFd, sockAddr); err != nil {
		fmt.Println("binding failed ", err)
		os.Exit(1)
	}

	fmt.Println("binding to port ", port)

	if err := syscall.Listen(proxyFd, max_clients); err != nil {
		fmt.Println("listen failed ", err)
		os.Exit(1)
	}

	fmt.Println("server is listening...")

	for {
		clientFD, clientAddr, err := syscall.Accept(proxyFd)
		if err != nil {
			fmt.Println("accepting error ", err)
			continue
		}

		if addr, ok := clientAddr.(*syscall.SockaddrInet4); ok {
			fmt.Println("client connected from ", addr.Addr[:], addr.Port)
		}

		go handleClient(clientFD)
	}
}
