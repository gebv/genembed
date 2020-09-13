package main

import (
	"log"

	"github.com/gebv/go-embed/example/somepkg"
)

func main() {
	log.Println("file1", string(getLocalFile("file1")))
	log.Println("file2", string(getLocalFile("file2")))
	log.Println("not exists", string(getLocalFile("not exists")))
	log.Println("file from pkg", somepkg.Value())
}

//go:generate $EMBEDBIN file1

//go:generate $EMBEDBIN file2
