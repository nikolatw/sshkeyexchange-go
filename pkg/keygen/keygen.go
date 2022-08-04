package keygen

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

type Keys struct {
	RSA     *rsa.PrivateKey
	Private []byte
	Public  []byte
}

func New() (*Keys, error) {
	keys := new(Keys)

	privateKey, publicKey, err := keys.newKeys()
	if err != nil {
		return nil, err
	}

	err = keys.marshal(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func NewWithPasscode(passcode string) (*Keys, error) {
	keys := new(Keys)

	privateKey, publicKey, err := keys.newKeys()
	if err != nil {
		return nil, err
	}

	privateKey, err = keys.addPasscode(privateKey, passcode)
	if err != nil {
		return nil, err
	}

	publicKey, err = keys.addPasscode(publicKey, passcode)
	if err != nil {
		return nil, err
	}

	err = keys.marshal(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (keys *Keys) marshal(privateKey *pem.Block, publicKey *pem.Block) error {
	encodedPrivateKey, err := keys.plain(privateKey)
	if err != nil {
		return err
	}
	keys.Private = encodedPrivateKey.Bytes()

	encodedPublicKey, err := keys.plain(publicKey)
	if err != nil {
		return err
	}
	keys.Public = encodedPublicKey.Bytes()

	return nil
}

func (keys *Keys) newKeys() (*pem.Block, *pem.Block, error) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

	keys.RSA = key

	privateKey := keys.makePrivateKey()
	publicKey, err := keys.makePublicKey()
	if err != nil {
		return nil, nil, err
	}
	return privateKey, publicKey, nil
}

func (keys *Keys) SSHPublicKey() ([]byte, error) {
	sshPub, err := ssh.NewPublicKey(&keys.RSA.PublicKey)
	if err != nil {
		return nil, err
	}

	sshPubBytes := sshPub.Marshal()
	parsed, err := ssh.ParsePublicKey(sshPubBytes)
	if err != nil {
		return nil, err
	}

	record := ssh.MarshalAuthorizedKey(parsed)
	return record, nil
}

func (keys *Keys) plain(block *pem.Block) (bytes.Buffer, error) {
	buf := bytes.Buffer{}

	err := keys.encodeToBuffer(&buf, block)
	if err != nil {
		return buf, err
	}

	return buf, nil
}

func (keys *Keys) makePrivateKey() *pem.Block {
	privateKey := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keys.RSA),
	}

	return privateKey
}

func (keys *Keys) makePublicKey() (*pem.Block, error) {
	asn1Bytes, err := asn1.Marshal(keys.RSA.PublicKey)

	publicKey := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	return publicKey, err
}

func (keys *Keys) encodeToBuffer(buf *bytes.Buffer, block *pem.Block) error {
	return pem.Encode(buf, block)
}

func (keys *Keys) addPasscode(block *pem.Block, passcode string) (*pem.Block, error) {
	block, err := x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(passcode), x509.PEMCipherAES256)
	if err != nil {
		return nil, err
	}

	return block, nil
}
