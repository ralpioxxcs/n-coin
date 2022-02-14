package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	enc := gob.NewEncoder(&aBuffer)
	HandleErr(enc.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(dec.Decode(i))
}

func Hash(i interface{}) string {
	s := fmt.Sprint(i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s string, sep string, i int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]

}

func ToJson(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}
