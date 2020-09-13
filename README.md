# go-embed

Embed resource files into your code via `go generate` command.

This is an **example** of embedding files via `go generate ...` command in your exists code. No binaries used. Generating with a bash script [embed.sh](embed.sh)

**NOTE: Please be mindful. Check what files you embed. The content of these files will be available from your code.**

# Quickstart

1. Add a comment `go:generate ...` to your code for each file.
Follow example code ([see more](example))
```go
...
package main

import "log"

func main() {
	log.Println(string(getLocalFile("file1")))
	log.Println(string(getLocalFile("file2")))
	log.Println(string(getLocalFile("not exists file")))
}

//go:generate ../embed.sh file1

//go:generate ../embed.sh file2

```

2. Run the script to update the embedded files

```bash
find ./* -name '*_embeded.go' -print0 | xargs -0 rm
go generate ./...
find ./* -name '*_embeded.go.tmp' -print0 | xargs -0 rm
find ./* -name '*_embeded.go' -exec gofmt -w {} +
```

3. The original code has the expected content. Why? Because in the same package the file was created. Below is an example the generated file (by pattern `<packagename>_embeded.go`).

```go
// Code generated  DO NOT EDIT.
package main

// local files in the current package
var localFiles = map[string][]byte{
	// embemded files
	// @InsertAfterBreakpoint
	"file2": {0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x20, 0x66, 0x69, 0x6c, 0x65, 0x32, 0x0a},
	"file1": {0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x20, 0x66, 0x69, 0x6c, 0x65, 0x31, 0x0a},
}

// getLocalFile returns slice bytes of local file if eixsts
// Returns nil if not registred file
// TODO: nomemcopy?
func getLocalFile(name string) []byte {
	dat, exists := localFiles[name]
	if !exists {
		return nil
	}
	return dat
}

```
