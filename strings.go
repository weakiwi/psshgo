package main

import (
	"bufio"
	"crypto/md5"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"strings"
)

func mustGetStringVar(c *cli.Context, key string) string {
	v := strings.TrimSpace(c.String(key))
	if v == "" {
		log.Fatalf("%s must be provided", key)
	}
	return v
}

func md5File(srcfile string) {
	file, err := os.Open(srcfile)
	if err != nil {
		log.Fatalf("md5Fiel os.Open error: %v", err)
		return
	}

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return
	}
	log.Printf("%x  %s\n", h.Sum(nil), srcfile)
}

func ComputeLine(path string) (num int) {
	f, err := os.Open(path)
	if nil != err {
		log.Fatalf("ComputeLine os.Open error: %v", err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		_, err := r.ReadString('\n')
		if io.EOF == err || nil != err {
			break
		}
		num += 1
	}
	return
}
