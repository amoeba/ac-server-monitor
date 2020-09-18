package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func isup(connection string) bool {
	fmt.Println("Checking ", connection)
	conn, connerror := net.DialTimeout("udp", connection, 5*time.Second)

	if connerror != nil {
		fmt.Println("Error:", connerror)
		return false
	}

	message := []uint8{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x93, 0x00,
		0xd0, 0x05, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x31, 0x38, 0x30, 0x32, 0x00, 0x00, 0x34, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x3e, 0xb8, 0xa8, 0x58, 0x1c, 0x00, 0x61, 0x63, 0x73, 0x65,
		0x72, 0x76, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65,
		0x72, 0x3a, 0x6a, 0x6a, 0x39, 0x68, 0x32, 0x36, 0x68, 0x63,
		0x73, 0x67, 0x67, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	loginpacket := iatoba(message)
	conn.Write(loginpacket)
	readbuffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, readerr := conn.Read(readbuffer)

	if readerr != nil {
		fmt.Println("Error", readerr)
		return false
	}

	if n == 52 {
		return true
	}

	return false
}

func iatoba(input []uint8) []byte {
	buffer := new(bytes.Buffer)
	writeerr := binary.Write(buffer, binary.LittleEndian, input)

	if writeerr != nil {
		fmt.Println("binary.Write failed:", writeerr)
		panic(1)
	}

	return (buffer.Bytes())
}
