package main

type findNodes struct {
	Y string `bencode:"y"`
	T string `bencode:"t"`
	Q string `bencode:"q"`
	R struct {
		Id    string `bencode:"id"`
		Nodes string `bencode:"nodes"`
	} `bencode:"r"`
	Ip string `bencode:"ip"`
}

type Noda struct {
	id []byte
	ip string
}

type hashSamples struct {
	Y string `bencode:"y"`
	T string `bencode:"t"`
	R struct {
		Id       string `bencode:"id"`
		Interval int    `bencode:"interval"`
		Nodes    string `bencode:"nodes"`
		Num      int    `bencode:"num"`
		Samples  string `bencode:"samples"`
	} `bencode:"r"`
	Ip string `bencode:"ip"`
}

type getPeer struct {
	V string `bencode:"v"`
	Y string `bencode:"y"`
	T string `bencode:"t"`
	R struct {
		P      int      `bencode:"p"`
		Id     string   `bencode:"id"`
		Nodes  string   `bencode:"nodes"`
		Token  string   `bencode:"token"`
		Values []string `bencode:"values"`
	} `bencode:"r"`
	Ip string `bencode:"ip"`
}

type TorrAttr struct {
	Hash   string
	Time   string
	Weight string
	Name   string
}
