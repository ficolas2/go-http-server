package tls

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	// "encoding/hex"
	// "strings"
)

const (
	ContentType_ChangeCipherSpec byte = 0x14
	ContentType_Alert                 = 0x15
	ContentType_Handshake             = 0x16
	ContentType_Application           = 0x17
	ContentType_Heartbeat             = 0x18
)

const (
	Version_SSL3_0 int = 0x0300
	Version_TLS1_0     = 0x0301
	Version_TLS1_1     = 0x0302
	Version_TLS1_2     = 0x0303
	Version_TLS1_3     = 0x0304
)

const (
	HandshakeType_HelloRequest          int = 0
	HandshakeType_ClientHello               = 1
	HandshakeType_ServerHello               = 2
	HandshakeType_HelloVerifyRequest        = 3
	HandshakeType_NewSessionTicket          = 4
	HandshakeType_EndOfEarlyData            = 5
	HandshakeType_HelloRetryRequest         = 6
	HandshakeType_EncryptedExtensions       = 8
	HandshakeType_Certificate               = 11
	HandshakeType_ServerKeyExchange         = 12
	HandshakeType_CertificateRequest        = 13
	HandshakeType_ServerHelloDone           = 14
	HandshakeType_CertificateVerify         = 15
	HandshakeType_ClientKeyExchange         = 16
	HandshakeType_Finished                  = 20
	HandshakeType_CertificateURL            = 21
	HandshakeType_CertificateStatus         = 22
	HandshakeType_SupplementalData          = 23
	HandshakeType_KeyUpdate                 = 24
	HandshakeType_CompressedCertificate     = 25
	HandshakeType_MessageHash               = 254
)

type RecordHeader struct {
	ContentType byte
	Version     uint16
	Length      uint16
}

type HandshakeHeader struct {
	HandshakeType byte
	Length        uint32
}

type ClientHello struct {
	ClientVersion uint16
	ClientRandom []byte 
	SessionId    []byte
	CipherSuites []uint16
	CompressionMethods []byte
	Extensions []byte
}

type ServerHello struct {
	ServerVersion uint16
	ServerRandom []byte
	SessionId []byte
	CipherSuite uint16
	CompressionMethod byte
	Extensions []byte
}

func HandleConnection(connection net.Conn) {
	buf := make([]byte, 8192)
	n, err := connection.Read(buf)
	if err != nil {
		log.Fatal(err)
		return
	}

	hexdump(buf[:n])
	// newBuf := "16 03 01 00 a5 01 00 00 a1 03 03 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f 00 00 20 cc a8 cc a9 c0 2f c0 30 c0 2b c0 2c c0 13 c0 09 c0 14 c0 0a 00 9c 00 9d 00 2f 00 35 c0 12 00 0a 01 00 00 58 00 00 00 18 00 16 00 00 13 65 78 61 6d 70 6c 65 2e 75 6c 66 68 65 69 6d 2e 6e 65 74 00 05 00 05 01 00 00 00 00 00 0a 00 0a 00 08 00 1d 00 17 00 18 00 19 00 0b 00 02 01 00 00 0d 00 12 00 10 04 01 04 03 05 01 05 03 06 01 06 03 02 01 02 03 ff 01 00 01 00 00 12 00 00"
	// newBuf = strings.ReplaceAll(newBuf, " ", "")
	//
	// buf, _ = hex.DecodeString(newBuf)

	recordHeader := readTlsRecordHeader(buf)
	if recordHeader.ContentType != ContentType_Handshake {
		log.Fatal("Not a handshake")
		return //TODO error
	}

	handshakeHeader := readHandshake(buf)
	if handshakeHeader.HandshakeType != HandshakeType_ClientHello {
		log.Fatal("Not a client hello")
		return //TODO error
	}
	clientHello := readClientHello(buf)

	handleClientHello(clientHello)

	// Write server hello
	
}

func readTlsRecordHeader(buf []byte) *RecordHeader {
	return &RecordHeader {
		ContentType: buf[0],
		Version: binary.BigEndian.Uint16(buf[1:3]),
		Length: binary.BigEndian.Uint16(buf[3:5]),
	}
}

func readHandshake(buf []byte) *HandshakeHeader {
	handshakeHeaderBuf := buf[5:9]
	handshakeHeader := HandshakeHeader{
		HandshakeType: handshakeHeaderBuf[0],
		Length: uint32(handshakeHeaderBuf[1]) << 16 | uint32(handshakeHeaderBuf[2]) << 8 | uint32(handshakeHeaderBuf[3]),
	}

	return &handshakeHeader
}

func readClientHello(buf []byte) *ClientHello {
	sessionIdLength := buf[43]
	sessionId := buf[44:44+sessionIdLength]

	start := int(44 + sessionIdLength)
	cipherSuiteLength := binary.BigEndian.Uint16(buf[start:start+2])
	cipherSuites := make([]uint16, cipherSuiteLength/2)
	for i := 0; i < int(cipherSuiteLength); i += 2 {
		cipherSuite := binary.BigEndian.Uint16(buf[start+2+i:start+4+i])
		cipherSuites[i/2] = cipherSuite
	}

	start += int(cipherSuiteLength + 2)
	compressionMethodLength := int(buf[start])
	compressionMethods := buf[start+1:start+1+compressionMethodLength]

	// TODO extension
	// start += compressionMethodLength + 1
	// extension_length := binary.BigEndian.Uint16(buf[start:start+2])


	clientHello := ClientHello{
		ClientVersion: binary.BigEndian.Uint16(buf[9:11]),
		ClientRandom: buf[11:11+32],
		SessionId: sessionId,
		CipherSuites: cipherSuites,
		CompressionMethods: compressionMethods,
		// Extensions: TODO,
	}

	return &clientHello
}

func handleClientHello(clientHello *ClientHello) *ServerHello {
	// Using TLS1.2, so we assert that the client version is at least 1.2
	if clientHello.ClientVersion > Version_TLS1_2 {
		log.Fatal(fmt.Sprintf("Unsupported version %04X", clientHello.ClientVersion))
	}

	random := make([]byte, 32)
	rand.Read(random)

	return &ServerHello{
		ServerVersion: Version_TLS1_2,
		ServerRandom: random,
		CipherSuite: 0x1301,
		CompressionMethod: 0,
		Extensions: []byte{},
	}

}
