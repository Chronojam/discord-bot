package main

import (
	"os"
	"io/ioutil"
	"encoding/binary"
	"bytes"
)

func main() {
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Size of Opec File
	size := int16(len(b))

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, size)
	if err != nil {
		panic(err)
	}
	buf.Write(b)

	err = ioutil.WriteFile(os.Args[2], buf.Bytes(), 0755)
	if err != nil {
		panic(err)
	}
}