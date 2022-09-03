package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/cristalhq/bencode"
	bencode2 "github.com/jackpal/bencode-go"
)

var Id = "abcdefghij0123456789"

type Obj struct {
	addr   string
	nodaId []byte
	tcp    net.Conn
}

type FindNode struct {
	Ip []uint8 `bencode:"ip"`
	R  struct {
		Id    []uint8 `bencode:"id"`
		Nodes []uint8 `bencode:"nodes"`
	} `bencode:"r"`
	T []uint8 `bencode:"t"`
	Y []uint8 `bencode:"y"`
}

func main() {
	do0()
}

func do0() {
	conn, _ := net.Dial("udp", "router.bittorrent.com:6881")
	conn.Write([]byte(fmt.Sprintf("d1:ad2:id20:%s6:target20:%se1:q9:find_node1:t2:aa1:y1:qe", Id, Id)))
	//n, _ := bufio.NewReader(conn).Read(p)
	bf := make([]byte, 800)
	//bufio.NewReader(conn).Read(bf)
	var a FindNode
	conn.Close()
	//if n != 0 {
	bencode2.Unmarshal(conn, &a)
	var data interface{}
	bencode.Unmarshal(bf, &data)
	r := priv(data)["r"]
	nodes := priv(r)["nodes"].([]uint8)
	num := len(nodes) / 26
	for i := 0; i < num; i++ {
		noda := nodes[26*i : 26*(i+1)]
		id := noda[:20]
		ip := noda[20:24]
		port := getPort(noda[24:])
		addr := fmt.Sprintf("%v.%v.%v.%v:%v", ip[0], ip[1], ip[2], ip[3], port)
		//fmt.Println(addr)
		o := Obj{addr: addr, nodaId: id}
		//fmt.Println(hashToText(id))
		o.getHash()
	}
	//}
}

func handShake(addr, hash string) {
	q1 := make([]byte, 20)
	hex.Decode(q1, []byte(hash))
	q2 := make([]byte, 20)
	hex.Decode(q2, []byte(Id))
	conn, _ := net.DialTimeout("tcp4", addr, 5*time.Second)
	if conn != nil {
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
	fmt.Println("good connect")
	conn.Write([]byte("\x00\x00\x00\x1a\x14\x00d1:md11:ut_metadatai2eee\x00\x00\x00\x1b\x14\x02\x64\x38\x3a\x6d\x73\x67\x5f\x74\x79\x70\x65\x69\x30\x65\x35\x3a\x70\x69\x65\x63\x65\x69\x30\x65\x65"))
	p := make([]byte, 2000)
	n, _ := bufio.NewReader(conn).Read(p)
	if n != 0 {
		fmt.Println("++++++++++++++++++")
		fmt.Println(string(p))
		fmt.Println("++++++++++++++++++")
		hz := p[6:]
		fmt.Println(n)
		sp := bytes.Split(hz, []byte(":"))
		for k, v := range sp {
			if bytes.Contains(v, []byte("name")) {
				if len(sp) > k {
					fmt.Println(string(sp[k+1]))
					break
				}

			}
		}
		conn.Close()
		fmt.Println("-----------------")
	}

}

func getPeers(addr, hash string) {
	q1 := make([]byte, 20)
	hex.Decode(q1, []byte(hash))
	q2 := make([]byte, 20)
	hex.Decode(q2, []byte(Id))
	conn, err := net.Dial("udp", addr)
	checkErr(err)
	if conn != nil {
		q := []byte("d1:ad2:id20:")
		q = append(q, q2...)
		q = append(q, []byte("9:info_hash20:")...)
		q = append(q, q1...)
		q = append(q, []byte("e1:q9:get_peers1:t2:aa1:y1:qe")...)
		conn.Write(q)
		checkErr(err)
		p := make([]byte, 800)
		_, err = bufio.NewReader(conn).Read(p)
		conn.Close()
		var data interface{}
		bencode.Unmarshal(p, &data)
		if data != nil {
			r := priv(data)["r"]
			val := priv(r)["values"]
			if val != nil {
				valI := val.([]interface{})
				for _, v := range valI {
					el := v.([]uint8)
					port := el[4:]
					adr := fmt.Sprintf("%v.%v.%v.%v:%v", el[0], el[1], el[2], el[3], getPort(port))
					fmt.Println("---", adr)
					handShake(adr, hash)
				}
			}
		}
	} else {
		fmt.Println("no connect")
	}

}
