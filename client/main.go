package main

import (
	"flag"
	"fmt"
	"wj5633/my-cache-server/cache-benchmark/cache"
)

func main() {
	server := flag.String("h", "localhost", "cache server address")
	op := flag.String("c", "get", "command, could be get/set/del")
	key := flag.String("k", "", "key")
	value := flag.String("v", "", "value")
	flag.Parse()
	client := cache.New("tcp", *server)
	cmd := &cache.Cmd{*op, *key, *value, nil}
	client.Run(cmd)
	if cmd.Error != nil {
		fmt.Println("error:", cmd.Error)
	} else {
		fmt.Println(cmd.Value)
	}
}
