package main

import (
	"log"

	"github.com/gebv/genembed/example/somepkg"
)

func main() {
	log.Println("file1", string(EmbedFiles["file1"]))
	log.Println("file2", string(EmbedFiles["file2"]))
	log.Println("not exists", string(EmbedFiles["not exists"]))
	log.Println("file from pkg", string(somepkg.EmbedFiles["somefile"]))
}

//go:generate genembed EmbedFiles file1 file2
