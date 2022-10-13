package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	bencode "github.com/jackpal/bencode-go"
)

//go:embed index.html
var f embed.FS

var wsConn *websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var Id string = "abcdefghij0123456789"

// var Id string
var spisok = []string{"router.bittorrent.com:6881", "router.utorrent.com:6881", "dht.transmissionbt.com:6881", "dht.libtorrent.org:25401"}
var SpHash = []TorrAttr{}
var SpMeta = make(map[string]string)
var mut = &sync.RWMutex{}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	wsConn, _ = upgrader.Upgrade(w, r, nil)
	log.Println("Client Connected")
}

func main() {

	go mainCicle()
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", wsEndpoint)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFS(f, "index.html")
		//	tmpl, _ := template.ParseFiles("index.html")
		tmpl.Execute(w, SpHash)
	})

	mux.HandleFunc("/data.json", dataPage)

	http.ListenAndServe(":80", mux)

}

func mainCicle() {
	for {
		for _, v := range spisok {
			n := findNode(v)
			for _, v := range n {
				getHash(v.id, v.ip, 0)
			}
		}

	}
}

func dataPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	buf, _ := json.Marshal(SpHash)
	fmt.Fprint(w, string(buf))
}

func findNode(server string) []Noda {
	buf := make([]byte, 600)
	conn, _ := net.Dial("udp4", server)
	if conn != nil {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		conn.Write([]byte(fmt.Sprintf("d1:ad2:id20:%s6:target20:%se1:q9:find_node1:t2:aa1:y1:qe", Id, Id)))
		bufio.NewReader(conn).Read(buf)
		conn.Close()
		var n findNodes
		bencode.Unmarshal(bytes.NewReader(buf), &n)
		return getNodes(n.R.Nodes)
	}
	return []Noda{}
}

func handShake(addr, hash string) bool {
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
			return getName(conn, addr, hash)
		}
	}
	return false
}

func getName(conn net.Conn, addr, hash string) bool {
	oS := ""
	conn.Write([]byte("\x00\x00\x00\x1a\x14\x00d1:md11:ut_metadatai2eee\x00\x00\x00\x1b\x14\x02\x64\x38\x3a\x6d\x73\x67\x5f\x74\x79\x70\x65\x69\x30\x65\x35\x3a\x70\x69\x65\x63\x65\x69\x30\x65\x65"))
	p := make([]byte, 2000)
	n, _ := bufio.NewReader(conn).Read(p)
	if n != 0 {
		otv := string(p[:n])
		oS += otv
		for {
			n, _ = bufio.NewReader(conn).Read(p)
			if n != 0 {
				otv = string(p[:n])
				oS += otv
			} else {
				break
			}
		}

		if strings.Index(oS, "msg_typei1e5") != -1 {
			sp := strings.Split(oS, ":")
			for k, v := range sp {
				if strings.Contains(v, "name") {
					if k+1 < len(sp) {
						vr := sp[k+1]
						l := getLength(sp)
						tS := getTime()
						name := vr[:len(vr)-2]
						//	fmt.Println(tS + "\t" + "magnet:?xt=urn:btih:" + hash + "\t" + l + "\t" + name)

						mut.Lock()
						if wsConn != nil {
							buf, _ := json.Marshal(TorrAttr{Hash: hash, Time: tS, Weight: l, Name: name})
							wsConn.WriteMessage(1, buf)
						}
						SpHash = append(SpHash, TorrAttr{Hash: hash, Time: tS, Weight: l, Name: name})
						SpMeta[hash] = oS
						mut.Unlock()
						return true
					}
					return false
				}
			}
		}
	}
	conn.Close()
	return false
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
					if handShake(addr, hash) {
						break
					}
				}
			}
		}
	}
}
