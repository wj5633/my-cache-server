package main

import (
	"wj5633/my-cache-server/chapter1/server/cache"
	"wj5633/my-cache-server/chapter1/server/http"
	"wj5633/my-cache-server/chapter2/server/tcp"
)

func main() {
	ca := cache.New("inmemory")
	go tcp.New(ca).Listen()
	http.New(ca).Listen("")
}
