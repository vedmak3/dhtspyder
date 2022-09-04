package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"

	bencode "github.com/jackpal/bencode-go"
)

func hashToText(m []byte) (vr string) {
	for _, v := range m {
		el := fmt.Sprintf("%x", v)
		if len(el) == 1 {
			el = "0" + el
		}
		vr += el
	}
	return
}

func getPort(i []uint8) int {
	vr := fmt.Sprintf("%x%x", i[0], i[1])
	zn, _ := strconv.ParseInt(vr, 16, 32)
	return int(zn)
}

func recurs(n int) string {
	st := ""
	for i := 0; i < n; i++ {
		st += "+"
	}
	return st
}

func getHash(id []byte, addr string, rec int) {
	conn, _ := net.Dial("udp4", addr)
	defer conn.Close()
	if conn != nil {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		q1 := []byte("d1:ad2:id20:" + Id + "6:target20:")
		q1 = append(q1, id...)
		q1 = append(q1, []byte("e1:q17:sample_infohashes1:t2:aa1:y1:qe")...)
		conn.Write(q1)
		buf := make([]byte, 2000)
		bufio.NewReader(conn).Read(buf)
		var h hashSamples
		bencode.Unmarshal(bytes.NewReader(buf), &h)
		/*if rec <= 1 {
			ns := getNodes(h.R.Nodes)
			for _, v := range ns {
				getHash(v.id, v.ip, rec+1)
			}
		}*/
		if h.R.Num > 0 {
			b := []byte(h.R.Samples)
			for i := 0; i < len(b)/20; i++ {
				hash := hashToText(b[i*20 : (i+1)*20])
				go getPeers(addr, hash)
			}
		}
	}
}

func getNodes(nods string) []Noda {
	ns := []Noda{}
	nodes := []byte(nods)
	num := len(nodes) / 26
	for i := 0; i < num; i++ {
		noda := nodes[26*i : 26*(i+1)]
		id := noda[:20]
		ip := noda[20:24]
		port := getPort(noda[24:])
		addr := fmt.Sprintf("%v.%v.%v.%v:%v", ip[0], ip[1], ip[2], ip[3], port)
		ns = append(ns, Noda{id: id, ip: addr})
	}
	return ns
}
