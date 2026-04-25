package main

const max_clients = 10

var PORT = 8080
var proxy_socketID int

var sem = make(chan struct{}, max_clients)

func main() {
	var cache CacheStore = &Cache{}
	_ = cache
}
