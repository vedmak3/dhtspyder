package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	bencode "github.com/jackpal/bencode-go"
)

var Id = "abcdefghij0123456789"

func main() {
	for {
		n := findNode()
		for _, v := range n {
			getHash(v.id, v.ip, 0)
		}
	}
}

func findNode() []Noda {
	buf := make([]byte, 600)
	conn, _ := net.Dial("udp4", "router.bittorrent.com:6881")
	defer conn.Close()
	if conn != nil {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		conn.Write([]byte(fmt.Sprintf("d1:ad2:id20:%s6:target20:%se1:q9:find_node1:t2:aa1:y1:qe", Id, Id)))
		bufio.NewReader(conn).Read(buf)
		conn.Close()
		var n FindNodes
		bencode.Unmarshal(bytes.NewReader(buf), &n)
		return getNodes(n.R.Nodes)
	}
	return []Noda{}
}

func handShake(addr, hash string) {
	q1 := make([]byte, 20)
	hex.Decode(q1, []byte(hash))
	q2 := make([]byte, 20)
	hex.Decode(q2, []byte(Id))
	conn, _ := net.DialTimeout("tcp4", addr, 5*time.Second)
	if conn != nil {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		q := []byte("\x13\x42\x69\x74\x54\x6f\x72\x72\x65\x6e\x74\x20\x70\x72\x6f\x74\x6f\x63\x6f\x6c\x00\x00\x00\x00\x00\x10\x00\x05")
		q = append(q, q1...)
		q = append(q, q2...)
		conn.Write(q)
		p := make([]byte, 400)
		n, _ := bufio.NewReader(conn).Read(p)
		if n != 0 {
			getName(conn, addr, hash)
		}
	}
}

func getName(conn net.Conn, addr, hash string) {
	conn.Write([]byte("\x00\x00\x00\x1a\x14\x00d1:md11:ut_metadatai2eee\x00\x00\x00\x1b\x14\x02\x64\x38\x3a\x6d\x73\x67\x5f\x74\x79\x70\x65\x69\x30\x65\x35\x3a\x70\x69\x65\x63\x65\x69\x30\x65\x65"))
	p := make([]byte, 2000)
	n, _ := bufio.NewReader(conn).Read(p)
	if n != 0 {
		otv := string(p)
		if strings.Index(otv, "msg_typei1e5") != -1 {
			sp := strings.Split(otv, ":")
			for k, v := range sp {
				if strings.Contains(v, "name") {
					vr := sp[k+1]
					fmt.Println("magnet:?xt=urn:btih:"+hash, vr[:len(vr)-2])
				}
			}
		}
	}
	conn.Close()
}

func getPeers(addr, hash string) {
	q1 := make([]byte, 20)
	hex.Decode(q1, []byte(hash))
	q2 := make([]byte, 20)
	hex.Decode(q2, []byte(Id))
	conn, _ := net.Dial("udp4", addr)
	defer conn.Close()
	if conn != nil {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		q := []byte("d1:ad2:id20:")
		q = append(q, q2...)
		q = append(q, []byte("9:info_hash20:")...)
		q = append(q, q1...)
		q = append(q, []byte("e1:q9:get_peers1:t2:aa1:y1:qe")...)
		conn.Write(q)
		p := make([]byte, 800)
		bufio.NewReader(conn).Read(p)
		var per getPeer
		bencode.Unmarshal(bytes.NewReader(p), &per)
		if len(per.R.Values) > 0 {
			for _, v := range per.R.Values {
				el := []byte(v)
				if len(el) == 6 {
					ip := el[:4]
					port := getPort(el[4:])
					addr := fmt.Sprintf("%v.%v.%v.%v:%v", ip[0], ip[1], ip[2], ip[3], port)
					handShake(addr, hash)
				}
			}
		}
	}
}
