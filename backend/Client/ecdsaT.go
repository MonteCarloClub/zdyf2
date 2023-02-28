package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "fmt"
    "time"
)

func ecdsa_test() {
    fmt.Println("ecdsa test ---------------------")
    fmt.Print("ecdsa gen: ")
    start := time.Now().UnixNano()
    for i := 0; i < 100; i += 1 {
        privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

        msg := "hello, world"
        hash := sha256.Sum256([]byte(msg))

        ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
    }
    end := time.Now().UnixNano()
    fmt.Println(float64(end - start) / 1e9)

    fmt.Print("ecdsa gen & verify: ")
    start = time.Now().UnixNano()
    for i := 0; i < 100; i += 1 {
        privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

        msg := "hello, world"
        hash := sha256.Sum256([]byte(msg))

        sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
        if err != nil {
            panic(err)
        }

        ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
    }
    end = time.Now().UnixNano()
    fmt.Println(float64(end - start) / 1e9)

}
