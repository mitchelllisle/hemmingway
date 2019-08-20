package hemmingway

import (
	"crypto/rand"
	"golang.org/x/crypto/nacl/box"
	"io"
)

//sender := Crypto{}
//sender.GenerateKeyPair()
//
//receiver := Crypto{}
//receiver.GenerateKeyPair()
//
//encryptedMessage := sender.Encrypt(sender.PrivateKey, receiver.PublicKey, "encryption test")
//
//message := receiver.Decrypt(sender.PublicKey, receiver.PrivateKey, encryptedMessage)
//
//fmt.Print(string(message))
type Crypto struct {
	PublicKey *[32]byte
	PrivateKey *[32]byte
	PeerPublicKey *[32]byte
}

func (c *Crypto) GenerateKeyPair() {
	PublicKey, PrivateKey, err := box.GenerateKey(rand.Reader)
	FailOnError(err, "Generation of Key Pair failed")

	c.PublicKey = PublicKey
	c.PrivateKey = PrivateKey
}

func (c *Crypto) Encrypt(senderPrivateKey *[32]byte, recipientPublicKey *[32]byte, payload string) []byte {
	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.

	var nonce [24]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	FailOnError(err, "")

	msg := []byte(payload)
	// This encrypts msg and appends the result to the nonce.
	encrypted := box.Seal(nonce[:], msg, &nonce, recipientPublicKey, senderPrivateKey)

	return encrypted
}

func (c *Crypto) Decrypt(senderPublicKey *[32]byte, recipientPrivateKey *[32]byte, EncryptedMessage []byte) []byte {
	//The recipient can decrypt the message using their private key and the
	//sender's public key. When you decrypt, you must use the same nonce you
	//used to encrypt the message. One way to achieve this is to store the
	//nonce alongside the encrypted message. Above, we stored the nonce in the
	//first 24 bytes of the encrypted text.
	var decryptNonce [24]byte
	copy(decryptNonce[:], EncryptedMessage[:24])
	decrypted, ok := box.Open(nil, EncryptedMessage[24:], &decryptNonce, senderPublicKey, recipientPrivateKey)
	if !ok {
		panic("decryption error")
	}
	return decrypted
}
