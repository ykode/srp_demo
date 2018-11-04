package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
	"math/big"
	"os"
)

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
	N = new(big.Int).SetBytes(dhparam)
	g = big.NewInt(2)
	k = calculateHashBigInt(g, N)
)

const keyinfo = "SRP Demo Key Information"

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

func hkdfFromKey(salt []byte, ikm []byte, iteration int) [][]byte {
	hkdf := hkdf.New(sha256.New, ikm, salt, []byte(keyinfo))
	okm := make([]byte, 16)
	out := make([][]byte, iteration)

	for i := 0; i < iteration; i += 1 {
		io.ReadFull(hkdf, okm)
		out[i] = okm
	}

	return out
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username:")
	username, _ := reader.ReadString('\n')

	fmt.Print("Password:")
	password, _ := reader.ReadString('\n')

	a, _ := cryptrand(128)
	A := new(big.Int).Exp(g, a, N)

	fmt.Printf("a = 0x%x\nA (hex) : 0x%x\nA (base64): %s\n", a, A, base64.StdEncoding.EncodeToString(A.Bytes()))

	fmt.Print("Salt: ")
	saltStr, _ := reader.ReadString('\n')
	saltBytes, _ := base64.StdEncoding.DecodeString(saltStr)
	x := new(big.Int).SetBytes(calculateHash(saltBytes, []byte(username+":"+password)))
	v := new(big.Int).Exp(g, x, N)

	fmt.Printf("v = 0x%x\nv (base64): %s\n", v, base64.StdEncoding.EncodeToString(v.Bytes()))

	fmt.Print("B: ")
	BStr, _ := reader.ReadString('\n')
	BBytes, _ := base64.StdEncoding.DecodeString(BStr)
	B := new(big.Int).SetBytes(BBytes)

	u := calculateHashBigInt(A, B)

	S_c := new(big.Int).Exp(new(big.Int).Sub(B, new(big.Int).Mul(k, new(big.Int).Exp(g, x, N))),
		new(big.Int).Add(a, new(big.Int).Mul(u, x)), N)

	keyBytes := hkdfFromKey(u.Bytes(), S_c.Bytes(), 4)

	K_c := new(big.Int).SetBytes(keyBytes[0])

	fmt.Printf("S_c = 0x%x\n, K1_c =  0x%x\n", S_c, K_c)

	M1_c := calculateHashBigInt(K_c, new(big.Int).Exp(A, B, N))
	M1_base64 := base64.StdEncoding.EncodeToString(M1_c.Bytes())

	fmt.Printf("M1_c : %s\n", M1_base64)

	fmt.Print("M2_s: ")
	M2_base64, _ := reader.ReadString('\n')
	M2Bytes, _ := base64.StdEncoding.DecodeString(M2_base64)
	M2 := new(big.Int).SetBytes(M2Bytes)

	M2_c := calculateHashBigInt(K_c, new(big.Int).Exp(A, M1_c, N))

	if M2.Cmp(M2_c) == 0 {
		fmt.Println("Authorised")
	} else {
		fmt.Println("Wrong Username or password")
	}

}
