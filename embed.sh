#!/bin/bash

dst=${GOPACKAGE}_embeded.go
if [ ! -f "$dst" ]; then
    echo "// Code generated  DO NOT EDIT." >> $dst
    echo "package ${GOPACKAGE}" >> $dst
    echo "// local files in the current package" >> $dst
    echo "var localFiles = map[string][]byte{" >> $dst
    echo "// embemded files" >> $dst
    echo "// @InsertAfterBreakpoint" >> $dst
    echo "}" >> $dst
    echo "" >> $dst
    echo "// getLocalFile returns slice bytes of local file if eixsts" >> $dst
    echo "// Returns nil if not registred file" >> $dst
    echo "// TODO: nomemcopy?" >> $dst
    echo "func getLocalFile(name string) []byte {" >> $dst
    echo "dat, exists := localFiles[name]" >> $dst
    echo "if !exists {" >> $dst
    echo "return nil" >> $dst
    echo "}" >> $dst
    echo "return dat" >> $dst
    echo "}" >> $dst
fi

cat $1 | hexdump -ve '1/1 "0x%.2x,"' | xargs -J % echo "\"$1\": {" % "}," | sed -i'.tmp' -e '/@InsertAfterBreakpoint$/ r /dev/stdin' ${GOPACKAGE}_embeded.go
