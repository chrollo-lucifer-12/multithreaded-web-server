package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
)

const max_clients = 100

var clientSem = Init(max_clients)
var cache CacheStore = &Cache{}

func fetchFromOrigin(url string) string {
	fmt.Println("fetching from origin ", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return "Error fetching URL"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "Error reading response body"
	}

	return string(body)
}

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

	requested_url := strings.TrimSpace(string(buf[:n]))

	found_cache := cache.find(requested_url)

	if found_cache != nil {
		fmt.Println("cache hit ", found_cache.head.data)
		syscall.Write(clientFD, []byte(found_cache.head.data))
	} else {
		fmt.Println("cache miss")
		response := fetchFromOrigin(requested_url)
		syscall.Write(clientFD, []byte(response))
		cache.add_cache_element(response, len(response), requested_url)
	}

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
