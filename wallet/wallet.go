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
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ralpioxxcs/n-coin/utils"
)

const (
	signature     string = "e3fbceac0036f1231cce5c37cfc02d2c1dd69539dc0808b870370e09e356e317d095b5c814ee6cc71f1cb6e8469ba5558883e3306caecd23c3ab2df68311843b"
	privateKey    string = "30770201010420050f05a7c11b8b21486038edd756d1005483a60589883669a00300f099fdbe22a00a06082a8648ce3d030107a144034200046ae41d46087bdaca79457409321be11cda3deb0694b024e4ba63b2840325a36105381bda9cb7995eeee1644c58489a3caa7633aa7dd4d834f6137dc8a9d9f2cb"
	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
)

func Start() {

	// private-key 16진수(인코딩) 검증
	privByte, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)

	_, err = x509.ParseECPrivateKey(privByte)
	utils.HandleErr(err)

	sigBytes, err := hex.DecodeString(signature)
	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]

	var bigR, bigS = big.Int{}, big.Int{}
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)
	fmt.Println(bigR, bigS)

	/*
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

		keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)

		fmt.Printf("%x\n\n\n", keyAsBytes)

		utils.HandleErr(err)

		hashAsBytes, err := hex.DecodeString(hashedMessage)

		utils.HandleErr(err)

		// Sign require privateKey
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)

		signature := append(r.Bytes(), s.Bytes()...)

		fmt.Printf("%x\n\n", signature)

		utils.HandleErr(err)

		// Verify require publick key
		ok := ecdsa.Verify(&privateKey.PublicKey, hashAsBytes, r, s)

		fmt.Println(ok)
	*/
}
