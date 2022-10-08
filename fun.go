package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
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
		if h.R.Num > 0 {
			b := []byte(h.R.Samples)
			for i := 0; i < len(b)/20; i++ {
				hash := hashToText(b[i*20 : (i+1)*20])
				mut.RLock()
				_, fl := SpMeta[hash]
				mut.RUnlock()
				if !fl {
					go getPeers(addr, hash)
				}

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

func getTime() string {
	t := time.Now()
	hh := t.Hour()
	mm := t.Minute()
	return fmt.Sprintf("%v:%v", hh, mm)
}

func getLength(sp []string) string {
	var sum uint64
	for _, v := range sp {
		if strings.Contains(v, "lengthi") {
			a := strings.Index(v[7:], "e")
			if a > 0 {
				if a+7 < len(v) {
					vr, _ := strconv.ParseUint(v[7:a+7], 10, 64)
					sum += vr
				}
			}
		}
	}
	zn := float64(sum) / 1073741824
	return fmt.Sprintf("%.2f", zn)
}
