// Code generated  DO NOT EDIT.
package somepkg

// local files in the current package
var localFiles = map[string][]byte{
	// embemded files
	// @InsertAfterBreakpoint
	"somefile": {0x73, 0x6f, 0x6d, 0x65, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x0a},
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
