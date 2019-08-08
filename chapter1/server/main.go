package main

import (
	"wj5633/my-cache-server/chapter1/server/cache"
	"wj5633/my-cache-server/chapter1/server/http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen("")
}