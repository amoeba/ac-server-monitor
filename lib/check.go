package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

// check.go
//
// This method was first established by jovrtn <https://github.com/jovrtn> some
// time after the 2017 server shutdown and was initially used on asheron.net
// (now defunct). Emulator projects such as GDLE and ACE have added specialized
// support for the method and client launchers such as ThwargLauncher have
// replicated the approach. This code is essentially a 1:1 clone of the
// ThwargLauncher implementation.

const timeout = 30

// iotoba converts a uint8[] to a byte[]
func iatoba(input []uint8) []byte {
	buffer := new(bytes.Buffer)
	writeerr := binary.Write(buffer, binary.LittleEndian, input)

	if writeerr != nil {
		fmt.Println("binary.Write failed:", writeerr)
		panic(1)
	}

	return (buffer.Bytes())
}

// FakeLoginPacket() creates a byte[] suitable for sending to a server in order
// to check whether that server is up. The packet doesn't contain valid login
// credentials.
//
// It returns a byte[].
func FakeLoginPacket() []byte {
	raw := []uint8{
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

	return iatoba(raw)
}

// CheckResponseLength is used by Check to test whether the number of bytes
// returned matches a certain heuristical value
//
// It returns true or false depending on whether the number of bytes matches
// the expected value.
func CheckResponseLength(nbytes int) (bool, error) {
	if (nbytes == 52 || nbytes == 28) {
		return true, nil
	}

	return false, fmt.Errorf("In function CheckResponseLength, number of bytes read was neither 52 or 28 as expected.")
}

// Check checks whether or not a Server is up
//
// It returns true or false, depending on whether the server is up and may
// return an error if the checking process fails
func Check(srv Server) (bool, error) {
	connectionstring := fmt.Sprintf("%s:%s", srv.Host, srv.Port)
	conn, err := net.DialTimeout("udp", connectionstring, timeout*time.Second)

	if err != nil {
		return false, err
	}

	// Set up read and write deadlines so we mark as down and move on
	conn.SetReadDeadline(time.Now().Add(timeout * time.Second))
	conn.SetWriteDeadline(time.Now().Add(timeout * time.Second))

	// Send our fake login packet
	loginpacket := FakeLoginPacket()

	_, err = conn.Write(loginpacket)

	if err != nil {
		return false, err
	}

	readbuffer := make([]byte, 1024)

	var nbytes int
	nbytes, err = conn.Read(readbuffer)

	if err != nil {
		return false, err
	}

	_, err = CheckResponseLength(nbytes)

	if err != nil {
		return false, fmt.Errorf("Check of read response length failed, got: %X", readbuffer)
	}

	return true, nil
}

func CheckOne(host string, port string) {
	log.Printf("Checking %s:%s", host, port)

	s := Server{host, port}
	_, err := Check(s)

	if (err != nil) {
		log.Fatalf("Failed to check %s:%s.", host, port)
	}

	log.Printf("%s:%s is up!", host, port)
}
