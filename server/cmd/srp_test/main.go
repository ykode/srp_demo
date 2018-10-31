package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
	"math/big"
)

// Constant DHPARAM
var (
	dhparam = []byte{
		0xBD, 0xE5, 0xC8, 0x29, 0xE8, 0xD1, 0xFE, 0x9D, 0xD8, 0x51, 0xB3, 0xE7,
		0xC6, 0x3B, 0xA3, 0x58, 0xDD, 0xDE, 0x32, 0x9B, 0x98, 0x9A, 0x00, 0x49,
		0xAB, 0x00, 0x6A, 0xAD, 0xD8, 0x0A, 0xAC, 0xE8, 0xE3, 0xFF, 0xC2, 0x82,
		0xD8, 0x94, 0xB5, 0x72, 0x5F, 0x2D, 0x72, 0xD5, 0xD9, 0x87, 0x43, 0xFC,
		0xF1, 0xA9, 0xC0, 0x2C, 0x60, 0xB2, 0xED, 0xBD, 0xEA, 0x7B, 0x03, 0x28,
		0xD8, 0xD3, 0x65, 0x5E, 0xD9, 0xB1, 0x82, 0xBE, 0x6C, 0x5B, 0x03, 0xB5,
		0xC8, 0x4B, 0x75, 0x34, 0x40, 0x4D, 0x9A, 0x65, 0xD6, 0xE6, 0x49, 0xDF,
		0x5A, 0x28, 0xF5, 0x2A, 0xEF, 0x35, 0x3C, 0xA5, 0x4A, 0x45, 0x30, 0x14,
		0xFB, 0x37, 0xAE, 0x8F, 0x97, 0xC1, 0x92, 0x9B, 0x01, 0x2B, 0x16, 0xEA,
		0x21, 0xA0, 0x1A, 0xDD, 0xDF, 0xC4, 0xBA, 0x05, 0xBC, 0xC7, 0x4E, 0x8F,
		0x9A, 0x50, 0xE4, 0x22, 0x58, 0x0D, 0xFB, 0xCB,
	}

	salt = mustCryptRand(64)
)

const (
	FACTOR    = 2
	USERNAME  = "didiet"
	PASSWORD  = "somepass"
	HKDF_INFO = "Demo App Key Generation"
)

func calculateHash(salt []byte, payload []byte) []byte {
	mac := hmac.New(sha256.New, salt)
	mac.Write(payload)
	return mac.Sum(nil)
}

func calculateHashBigInt(salt *big.Int, payload *big.Int) *big.Int {
	return new(big.Int).SetBytes(calculateHash(salt.Bytes(), payload.Bytes()))
}

func cryptrand(length int) (*big.Int, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(b), nil
}

func mustCryptRand(length int) *big.Int {
	r, err := cryptrand(length)

	if err != nil {
		panic(err)
	}
	return r
}

func hkdfFromKey(ikm []byte, iteration int) []byte {
	hkdf := hkdf.New(sha256.New, ikm, salt.Bytes(), []byte(HKDF_INFO))
	okm := make([]byte, 16)

	for i := 0; i < iteration; i += 1 {
		io.ReadFull(hkdf, okm)
	}

	return okm
}

func hkdfFromBig(ikm *big.Int, iteration int) *big.Int {
	return new(big.Int).SetBytes(hkdfFromKey(ikm.Bytes(), iteration))
}

func main() {
	N := new(big.Int).SetBytes(dhparam)
	g := new(big.Int).SetInt64(FACTOR)
	k := calculateHashBigInt(N, g)

	fmt.Printf("H, N, g Values are well known, public in nature\nH = HMAC256\nN = 0x%x\ng = 0x%x\nk = 0x%x\n\n", N, g, k)
	fmt.Println("0. Registration, server save I, s, v")

	x := new(big.Int).SetBytes(calculateHash(salt.Bytes(), []byte(USERNAME+":"+PASSWORD)))
	v := new(big.Int).Exp(g, x, N)

	fmt.Printf("I (username) \t= %s\np (password) \t= %s\nx (private key, ephemeral) = 0x%x\nv = 0x%x\ns (salt) = 0x%x\n\n",
		USERNAME, PASSWORD, x, v, salt)

	fmt.Println("Prime Exchange")
	a := mustCryptRand(128)
	A := new(big.Int).Exp(g, a, N)
	fmt.Printf("Client -> Server\nI : %s\nA = 0x%x\n\n", USERNAME, A)

	b := mustCryptRand(128)
	B := new(big.Int).Add(new(big.Int).Mul(k, v), new(big.Int).Exp(g, b, N))
	fmt.Printf("Server -> Client\ns : 0x%x\nB = 0x%x\n\n", salt, B)

	fmt.Println("Clint and Server both calculate scrambling parameter U")
	u := calculateHashBigInt(A, B)
	fmt.Printf("Client & Server: 0x%x\n\n", u)

	x = new(big.Int).SetBytes(calculateHash(salt.Bytes(), []byte(USERNAME+":"+PASSWORD)))
	S_c := new(big.Int).Exp(new(big.Int).Sub(B, new(big.Int).Mul(k, new(big.Int).Exp(g, x, N))),
		new(big.Int).Add(a, new(big.Int).Mul(u, x)), N)
	K_c := hkdfFromBig(S_c, 1)

	fmt.Printf("Client Calculation For Session Keys\nS_c = 0x%x\nK_c = 0x%x\n\n", S_c, K_c)

	S_s := new(big.Int).Exp(new(big.Int).Mul(A, new(big.Int).Exp(v, u, N)), b, N)
	K_s := hkdfFromBig(S_s, 1)

	fmt.Printf("Server Calculation For Session Keys\nS_s = 0x%x\nK_s = 0x%x\n\n", S_s, K_s)

	M_1 := calculateHashBigInt(S_c, new(big.Int).Exp(A, B, N))
	M_2 := calculateHashBigInt(S_s, new(big.Int).Exp(A, M_1, N))

	fmt.Printf("Verifier Values\nM1 = 0x%x\nM2 = 0x%x\n\n", M_1, M_2)

}
