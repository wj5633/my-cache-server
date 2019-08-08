package cache

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	net.Conn
	r *bufio.Reader
}

func newTCPClient(server string) *tcpClient {
	c, e := net.Dial("tcp", server+":12346")
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)
	return &tcpClient{c, r}
}

func (c *tcpClient) sendGet(key string) {
	_, _ = c.Write([]byte(fmt.Sprintf("G%d %s", len(key), key)))
}

func (c *tcpClient) sendSet(key, value string) {
	_, _ = c.Write([]byte(fmt.Sprintf("S%d %d %s%s", len(key), len(value), key, value)))
}

func (c *tcpClient) sendDel(key string) {
	_, _ = c.Write([]byte(fmt.Sprintf("D%d %s", len(key), key)))
}

func readLen(r *bufio.Reader) int {
	tmp, e := r.ReadString(' ')
	if e != nil {
		log.Println(e)
		return 0
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		log.Println(e)
		return 0
	}
	return l
}

func (c *tcpClient) recvResponse() (string, error) {
	valLen := readLen(c.r)
	if valLen == 0 {
		return "", nil
	}
	if valLen < 0 {
		err := make([]byte, -valLen)
		_, e := io.ReadFull(c.r, err)
		if e != nil {
			return "", e
		}
		return "", errors.New(string(err))
	}
	value := make([]byte, valLen)
	_, e := io.ReadFull(c.r, value)
	if e != nil {
		return "", e
	}
	return string(value), nil
}

func (c *tcpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		c.sendGet(cmd.Key)
		cmd.Value, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "set" {
		c.sendSet(cmd.Key, cmd.Value)
		_, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "del" {
		c.sendDel(cmd.Key)
		_, cmd.Error = c.recvResponse()
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

func (c *tcpClient) PipelineRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	for _, cmd := range cmds {
		if cmd.Name == "get" {
			c.sendGet(cmd.Key)
		}
		if cmd.Name == "set" {
			c.sendSet(cmd.Key, cmd.Value)
		}
		if cmd.Name == "del" {
			c.sendDel(cmd.Key)
		}
	}
	for _, cmd := range cmds {
		cmd.Value, cmd.Error = c.recvResponse()
	}
}
