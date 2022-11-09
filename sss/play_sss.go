package main

import (
	"encoding/base64"
	"fmt"
	"github.com/dsprenkels/sss-go"
	"log"
	"strings"
)

func main() {
	// Make a new slice of secret data [42, ..., 42]

	str_data := "这是重要的东西1234是生生世世生生世世生生世世生生世世生生世世生生世世生生世世生生世世生生世世生生世世生生世世"

	chunks := make([][]byte, 0)
	for i, x := range []byte(str_data) {
		i0 := i / 64
		i1 := i % 64
		if len(chunks) == i0 {
			chunks = append(chunks, make([]byte, 64))
		}
		chunks[i0][i1] = x
	}
	count, threshould := 5, 3

	keys := make([]string, count)
	for i, chunk := range chunks {
		fmt.Println(i, chunk)
		// Create 5 shares; allow 4 to restore the original data
		shares, err := sss.CreateShares(chunk, count, threshould)
		if err != nil {
			log.Fatalln(err)
		}
		for i, share := range shares {
			if keys[i] != "" {
				keys[i] += "|"
			}
			keys[i] += base64.StdEncoding.EncodeToString(share)
		}
	}
	for i, key := range keys {
		fmt.Println(i, key)
	}

	restored_str := ""

	newKeys := []string{keys[1], keys[4], keys[2]}

	fmt.Println("restore with partial keys")
	newChunkedKeys := make([][]string, 0)
	for i, key := range newKeys {
		fmt.Println(i, key)
		keyChunks := strings.Split(key, "|")
		if len(newChunkedKeys) == 0 {
			newChunkedKeys = make([][]string, len(keyChunks))
			for j, _ := range newChunkedKeys {
				newChunkedKeys[j] = make([]string, len(newKeys))
			}
		} else if len(newChunkedKeys) != len(keyChunks) {
			log.Fatalln("key chunks not always the same")
		}
		for j, chunk := range keyChunks {
			newChunkedKeys[j][i] = chunk
		}
	}

	for _, chunk := range newChunkedKeys {
		new_shares := make([][]byte, len(chunk))
		for j, k := range chunk {
			new_shares[j], _ = base64.StdEncoding.DecodeString(k)
		}
		restoredChunk, err := sss.CombineShares(new_shares)
		if err != nil {
			log.Fatalln(err)
		}
		restored_str += string(restoredChunk)
	}
	// Try to restore the original secret
	fmt.Println(str_data)
	fmt.Println(restored_str)
	fmt.Println(len(restored_str), len(str_data))
}
