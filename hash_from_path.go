package main

import (
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
		fmt.Printf("usage: %s <path>\n", filepath.Base(os.Args[0]))
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

func computeHashes(output chan HashInfo, dirname string) {
	waiter := &sync.WaitGroup{}
	filepath.Walk(dirname, makeWalkFunc(output, waiter))
	waiter.Wait()
	close(output)
}

func makeWalkFunc(output chan HashInfo, waiter *sync.WaitGroup) func(string, os.FileInfo, error) error {
	return func(path string, fi os.FileInfo, err error) error {
		if err == nil && fi.Size() > 0 &&
			(fi.Mode()&os.ModeType == 0) {
			if runtime.NumGoroutine() > maxGoroutines {
				processFile(path, fi, output, nil)
			} else {
				waiter.Add(1)
				go processFile(path, fi, output, func() { waiter.Done() })
			}
		}
		return nil
	}
}

func processFile(filename string, fi os.FileInfo, output chan HashInfo, done func()) {
	if done != nil {
		defer done()
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
