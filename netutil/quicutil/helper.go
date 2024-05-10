package quicutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"sync"
)

const protos = "quic"

var (
	sererTLSConfig  *tls.Config
	clientTLSConfig *tls.Config

	onceServer sync.Once
	onceClient sync.Once
)

func ServerTLSConfig() *tls.Config {
	onceServer.Do(func() {
		key, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			panic(err)
		}
		template := x509.Certificate{SerialNumber: big.NewInt(1)}
		certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
		if err != nil {
			panic(err)
		}
		keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
		certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

		tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			panic(err)
		}
		sererTLSConfig = &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			NextProtos:   []string{protos},
		}
	})

	return sererTLSConfig
}

func ClientTLSConfig() *tls.Config {
	onceClient.Do(func() {
		clientTLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{protos},
		}
	})

	return clientTLSConfig
}
