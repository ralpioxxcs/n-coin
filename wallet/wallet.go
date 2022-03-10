// 1. we hash the msg
//
// 2. generate key pair
// 		Keypair (privateKey, publicKey) (save privateKeyq to a file)
//
// 3. sign the hash
// 		("hashed_message" + privateKey) -> "signature"
//
// 4. verify
//		("hashed_message" + "signature" + publicKey) -> true / false

package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/fs"
	"math/big"
	"os"

	"github.com/ralpioxxcs/n-coin/utils"
)

type fileLayer interface {
	hasWalletFile() bool
	writeFile(name string, date []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type layer struct{}

func (layer) hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

var files fileLayer = layer{}

const (
	fileName string = "nomadcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string // pubKey
}

var w *wallet

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = files.writeFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

// restoreKey
func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := files.readFile(fileName)
	utils.HandleErr(err)

	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)

	return // named return!
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

// sign something (privateKey)
func Sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleErr(err)

	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	utils.HandleErr(err)

	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)

	return &bigA, &bigB, nil
}

// publicKey 검증
func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)

	// restore publicKey
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)

	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if files.hasWalletFile() {
			// yes 	-> restore from file
			w.privateKey = restoreKey()
		} else {
			// no	-> create prv key, save to file
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}

func Start() {

	/*
		// private-key 16진수(인코딩) 검증
		privByte, err := hex.DecodeString(privateKey)
		utils.HandleErr(err)

		private, err := x509.ParseECPrivateKey(privByte)
		utils.HandleErr(err)

		bytes, err := hex.DecodeString(signature)
		firstHalfBytes := bytes[:len(bytes)/2]
		secondHalfBytes := bytes[len(bytes)/2:]

		var bigA, bigB = big.Int{}, big.Int{}
		bigA.SetBytes(firstHalfBytes)
		bigB.SetBytes(secondHalfBytes)

		hashBytes, err := hex.DecodeString(hashedMessage)
		utils.HandleErr(err)

		ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigA, &bigB)

		fmt.Println(ok)
	*/

	/*
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		utils.HandleErr(err)
		keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
		utils.HandleErr(err)
		fmt.Printf("privateKey(bytes):%x\n\n\n", keyAsBytes)

		hashAsBytes, err := hex.DecodeString(hashedMessage)
		utils.HandleErr(err)

		// Sign require privateKey
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
		utils.HandleErr(err)
		signature := append(r.Bytes(), s.Bytes()...)
		fmt.Printf("signature(bytes):%x\n\n", signature)

		// Verify require publick key
		ecdsa.Verify(&privateKey.PublicKey, hashAsBytes, r, s)
	*/
}
