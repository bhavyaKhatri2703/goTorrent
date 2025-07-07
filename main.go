package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/zeebo/bencode"
)

type tracker struct {
	url        string
	len        int
	infohash   []byte
	pieces     []byte
	piecesHash [][]byte
}

func decode(s []byte) (interface{}, int) {
	switch firstChar := s[0]; {
	case firstChar == 'i':
		return decodeInteger(s)
	case firstChar >= '0' && firstChar <= '9':
		return decodeStrings(s)
	case firstChar == 'l':
		return decodeLists(s)
	case firstChar == 'd':
		return decodeDict(s)
	}
	return nil, 0
}
func decodeInteger(s []byte) (int, int) {
	temp := ""
	end := 1
	for i := 1; s[i] != 'e'; i++ {
		temp += string(s[i])
		end++
	}
	num, _ := strconv.Atoi(temp)
	return num, end + 1
}
func decodeStrings(s []byte) ([]byte, int) {
	var found int
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			found = i
			break
		}
	}
	num, _ := strconv.Atoi(string(s[:found]))
	ans := s[found+1 : found+1+num]
	return ans, num + 1 + len(s[:found])
}
func decodeLists(s []byte) ([]interface{}, int) {
	s = s[1:]
	var arr []interface{}
	totalConsumed := 1
	for {
		val, consumed := decode(s)
		totalConsumed += consumed
		s = s[consumed:]
		arr = append(arr, val)
		if s[0] == 'e' {
			totalConsumed += 1
			break
		}
	}
	return arr, totalConsumed
}
func decodeDict(s []byte) (map[string]interface{}, int) {
	s = s[1:]
	mp := make(map[string]interface{})
	totalConsumed := 1
	for {
		key, consumed := decodeStrings(s)
		totalConsumed += consumed
		s = s[consumed:]
		val, consumed := decode(s)
		totalConsumed += consumed
		s = s[consumed:]
		mp[string(key)] = val
		if s[0] == 'e' {
			totalConsumed += 1
			break
		}
	}
	return mp, totalConsumed
}
func readfile(path string) []byte {
	data, err := os.ReadFile(path)
	if err == nil {
		return data
	}
	return []byte{}
}

func calcInfoHash(infoDict map[string]interface{}) []byte {

	enc, err := bencode.EncodeBytes(infoDict)
	if err != nil {

	}

	sha := sha1.New()
	sha.Write(enc)
	hash := sha.Sum(nil)

	return hash
}

// func calcPiecesHash

func getTrackerDetails(torrfile string) tracker {
	data := readfile(torrfile)
	ans, _ := decode(data)

	dict := ans.(map[string]interface{})
	length := dict["info"].(map[string]interface{})["piece length"]
	pieces := dict["info"].(map[string]interface{})["pieces"].([]byte)

	var piecesHash [][]byte
	for i := 0; i < len(pieces); i += 20 {
		if i+20 > len(pieces) {
			piecesHash = append(piecesHash, pieces[i:])
			break
		}
		piecesHash = append(piecesHash, pieces[i:i+20])
	}

	tr := tracker{url: string(dict["announce"].([]byte)), len: length.(int), infohash: calcInfoHash(dict["info"].(map[string]interface{})), pieces: pieces,
		piecesHash: piecesHash,
	}

	return tr
}
func main() {
	// num, end := decodeInteger("i11ei3232e")
	// ans, e := decodeStrings("11:abcdabcdabcli11e10:helicoptere10:helicopteri11ee")
	// ans, end := decodeLists("li11e10:helicoptere")
	// input := `d8:announce41:http://bttracker.debian.org:6969/announce7:comment35:"Debian CD from cdimage.debian.org"13:creation datei1573903810e9:httpseedsl145:https://cdimage.debian.org/cdimage/release/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.iso145:https://cdimage.debian.org/cdimage/archive/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.isoe4:infod6:lengthi351272960e4:name31:debian-10.2.0-amd64-netinst.iso12:piece lengthi262144e6:pieces26800:�����PS�^�� (binary blob of the hashes of each piece)ee`
	// ans, end := decode(input)
	// fmt.Println(ans, end)
	// data := readfile("sample.torrent")
	// ans, _ := decode(data)

	// dict := ans.(map[string]interface{})
	// fmt.Println(dict["info"].(map[string]interface{})["pieces"])
	url := getPeers()
	fmt.Println(url)
	resp, er := http.Get(url)
	body, _ := io.ReadAll(resp.Body)
	if er == nil {
		fmt.Println((body))
	}

	bodyDecoded, _ := decode(body)
	peers := parsePeers(bodyDecoded.(map[string]interface{})["peers"].([]byte))
	tcp_conn(peers[0])
	// fmt.Println(binary.BigEndian.Uint16(bodyDecoded.(map[string]interface{})["peers"].([]byte)[4:6]))

}
