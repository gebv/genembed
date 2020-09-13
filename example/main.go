package main

import "log"

func main() {
	log.Println(string(getLocalFile("file1")))
	log.Println(string(getLocalFile("file2")))
	log.Println(string(getLocalFile("not exists")))
}

//go:generate ../embed.sh file1

//go:generate ../embed.sh file2
