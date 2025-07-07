package main

import (
	"encoding/binary"
	"net"
	"net/url"
	"strconv"
)

type peer struct {
	addresss string
	port     int
}

type hanshake struct {
}

func getPeers() string {
	tr := getTrackerDetails("sample.torrent")
	baseurl, _ := url.Parse(tr.url)

	params := url.Values{
		"info_hash":  []string{string(tr.infohash)},
		"peer_id":    []string{"-GO0001-bhavyakhatri"},
		"port":       []string{"6881"},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(tr.len)},
	}

	baseurl.RawQuery = params.Encode()

	return baseurl.String()
}

func parsePeers(s []byte) []string {
	var ips []string
	for i := 0; i < len(s); i += 6 {
		temp := strconv.Itoa(int(s[i])) + "." + strconv.Itoa(int(s[i+1])) + "." + strconv.Itoa(int(s[i+2])) + "." + strconv.Itoa(int(s[i+3])) + ":" + strconv.Itoa(int(binary.BigEndian.Uint16(s[i+4:i+6])))
		ips = append(ips, temp)
	}

	return ips
}

func tcp_conn(peer string) {
	conn, _ := net.Dial("tcp", peer)
	defer conn.Close()
}
