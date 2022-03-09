package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"testing"
)

const (
	testKey    = "3077020101042093f488e1b1251e0383d93af7598cf1758f61973e23ae165890cd42b1aa60f248a00a06082a8648ce3d030107a144034200049cf705e8a504ccfe63958d5d2cb419041ba297db9be21ad597f0eabfa5b0f0971edbb4e87e34e69f2758b6ac124be69b9e297c416c7c9725bd02d0f2654d37bc"
	testPaylod = "0000f4b3f80d8bd561c4281c18da4ed780902fd583ec9a3ce1d5d9369bc3e628"
	testSig    = "12fa326c2553cee4473392d930426b2ad89431bd96287c158543a75051ee0790617e9e9b74a0cdffba5baa21727c2fc4ef3d049edc1b9c9f0c13e473d4f1bd35"
)

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)

	return w
}

func TestSign(t *testing.T) {
	s := Sign(testPaylod, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign should return a hex encoded string, got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{input: testPaylod, ok: true},
		{input: "4000f4b3f80d8bd561c4281c18da4ed780902fd583ec9a3ce1d5d9369bc3e628", ok: false},
	}

	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify test-signature and test-payload")
		}
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx")
	if err == nil {
		t.Error("restoreBigInts() should return error when payload is not hex")
	}
}
