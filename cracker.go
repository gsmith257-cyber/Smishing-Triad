package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

var commonSuffixes = []string{"1", "12", "123", "123456", "!", "$", "!!", "..", "@", "234", "456"}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: cracker.exe <wordlist> <hashes>")
		os.Exit(1)
	}

	wordlistPath := os.Args[1]
	hashesPath := os.Args[2]

	wordlistFile, err := os.Open(wordlistPath)
	if err != nil {
		panic(err)
	}
	defer wordlistFile.Close()

	hashesFile, err := os.Open(hashesPath)
	if err != nil {
		panic(err)
	}
	defer hashesFile.Close()

	targetHashes := make([]string, 0)
	scanner := bufio.NewScanner(hashesFile)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)
	for scanner.Scan() {
		targetHashes = append(targetHashes, scanner.Text())
	}

	sem := make(chan bool, 100)

	scanner = bufio.NewScanner(wordlistFile)
	scanner.Buffer(buf, 10*1024*1024)
	for scanner.Scan() {
		word := scanner.Text()

		for _, suffix := range commonSuffixes {
			fullWord := word + suffix

			sem <- true
			go func(fullWord string) {
				defer func() { <-sem }()

				buf := pool.Get().(*bytes.Buffer)
				defer pool.Put(buf)

				buf.Reset()
				buf.WriteString(fullWord + "wangduoyu666!.+-")

				hash := md5.Sum(buf.Bytes())
				hashString := hex.EncodeToString(hash[:])

				hash = md5.Sum([]byte(hashString))
				hashString = hex.EncodeToString(hash[:])

				hash = md5.Sum([]byte(hashString))
				hashString = hex.EncodeToString(hash[:])

				for _, targetHash := range targetHashes {
					if strings.EqualFold(hashString, targetHash) {
						fmt.Printf("Found match: %s -> %s\n", fullWord, hashString)
						return
					}
				}
			}(fullWord)
		}

		// Try the regular word without the prefix or suffix
		sem <- true
		go func(word string) {
			defer func() { <-sem }()

			buf := pool.Get().(*bytes.Buffer)
			defer pool.Put(buf)

			buf.Reset()
			buf.WriteString(word + "wangduoyu666!.+-")

			hash := md5.Sum(buf.Bytes())
			hashString := hex.EncodeToString(hash[:])

			hash = md5.Sum([]byte(hashString))
			hashString = hex.EncodeToString(hash[:])

			hash = md5.Sum([]byte(hashString))
			hashString = hex.EncodeToString(hash[:])

			for _, targetHash := range targetHashes {
				if strings.EqualFold(hashString, targetHash) {
					fmt.Printf("Found match: %s -> %s\n", word, hashString)
					return
				}
			}
		}(word)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}
