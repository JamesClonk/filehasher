package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const maxGoroutines = 7

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <filelist>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	out := make(chan HashInfo, maxGoroutines*2)
	go computeHashes(out, os.Args[1])
	output(out)
}

type HashInfo struct {
	md5  []byte
	path string
}

func computeHashes(output chan HashInfo, filename string) {
	waiter := &sync.WaitGroup{}

	filelist, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer filelist.Close()

	reader := bufio.NewReader(filelist)
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			if runtime.NumGoroutine() > maxGoroutines {
				processFile(line, output, nil)
			} else {
				waiter.Add(1)
				go processFile(line, output, func() { waiter.Done() })
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
	}

	waiter.Wait()
	close(output)
}

func processFile(filename string, output chan HashInfo, done func()) {
	if done != nil {
		defer done()
	}

	fi, err := os.Stat(filename)
	if err != nil {
		log.Println(err)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	hash := md5.New()
	if size, err := io.Copy(hash, file); size != fi.Size() || err != nil {
		if err != nil {
			log.Println(err)
		} else {
			log.Println("could not read the whole file:", filename)
		}
		return
	}

	output <- HashInfo{hash.Sum(nil), filename}
}

func output(output <-chan HashInfo) {
	for hashinfo := range output {
		fmt.Printf("%x    %s\n", hashinfo.md5, hashinfo.path)
	}
}
