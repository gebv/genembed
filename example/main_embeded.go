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
