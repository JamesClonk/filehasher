#!/usr/bin/env goplay

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <path>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	filepath.Walk(os.Args[1], makeWalkFunc())
}

func makeWalkFunc() func(string, os.FileInfo, error) error {
	return func(path string, fi os.FileInfo, err error) error {
		if err == nil && fi.Size() > 0 &&
			(fi.Mode()&os.ModeType == 0) {
			fmt.Println(path)
		}
		return nil
	}
}
