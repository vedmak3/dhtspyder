package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/cristalhq/bencode"
)

func hashToText(m []uint8) (vr string) {
	for _, v := range m {
		el := fmt.Sprintf("%x", v)
		if len(el) == 1 {
			el = "0" + el
		}
		vr += el
	}
	return
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func priv(i interface{}) map[string]interface{} {
	return i.(map[string]interface{})
}
func getPort(i []uint8) int {
	vr := fmt.Sprintf("%x%x", i[0], i[1])
	zn, _ := strconv.ParseInt(vr, 16, 32)
	return int(zn)
}

func (o Obj) getHash() {
	conn, err := net.Dial("udp", o.addr)
	err = conn.SetDeadline(time.Now().Add(3 * time.Second))
	q1 := []byte("d1:ad2:id20:" + Id + "6:target20:")
	q1 = append(q1, o.nodaId...)
	q1 = append(q1, []byte("e1:q17:sample_infohashes1:t2:aa1:y1:qe")...)
	conn.Write(q1)
	p := make([]byte, 2000)
	_, err = bufio.NewReader(conn).Read(p)
	if p != nil {
		checkErr(err)
		var data interface{}
		bencode.Unmarshal(p, &data)
		if data != nil {
			r := priv(data)["r"]
			if r != nil {
				if priv(r)["samples"] != nil {
					m := priv(r)["samples"].([]uint8)
					num := len(m) / 20
					for i := 0; i < num; i++ {
						hash := hashToText(m[20*i : 20*(i+1)])
						fmt.Println("magnet:?xt=urn:btih:" + hash)
						getPeers(o.addr, hash)
					}

				}
			}
		}
	}
}
