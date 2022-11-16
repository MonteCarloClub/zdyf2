package main

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math/big"
    "time"
)

func fromBase10(base10 string) *big.Int {
    i, ok := new(big.Int).SetString(base10, 10)
    if !ok {
        panic("bad number: " + base10)
    }
    return i
}

func rsa_test() {
    fmt.Println("rsa test ---------------------")
    fmt.Print("rsa gen: ")
    start := time.Now().UnixNano()
    for i := 0; i < 100; i += 1 {
        rng := rand.Reader

        message := []byte("message to be signed")

        var rsaPrivateKey = &rsa.PrivateKey{
            PublicKey: rsa.PublicKey{
                N: fromBase10("9353930466774385905609975137998169297361893554149986716853295022578535724979677252958524466350471210367835187480748268864277464700638583474144061408845077"),
                E: 65537,
            },
            D: fromBase10("7266398431328116344057699379749222532279343923819063639497049039389899328538543087657733766554155839834519529439851673014800261285757759040931985506583861"),
            Primes: []*big.Int{
                fromBase10("98920366548084643601728869055592650835572950932266967461790948584315647051443"),
                fromBase10("94560208308847015747498523884063394671606671904944666360068158221458669711639"),
            },
        }

        hashed := sha256.Sum256(message)

        rsa.SignPKCS1v15(rng, rsaPrivateKey, crypto.SHA256, hashed[:])
    }
    end := time.Now().UnixNano()
    fmt.Println(float64(end - start) / 1e9)

    fmt.Print("rsa gen & verify: ")
    start = time.Now().UnixNano()
    for i := 0; i < 100; i += 1 {
        rng := rand.Reader

        message := []byte("message to be signed")

        var rsaPrivateKey = &rsa.PrivateKey{
            PublicKey: rsa.PublicKey{
                N: fromBase10("9353930466774385905609975137998169297361893554149986716853295022578535724979677252958524466350471210367835187480748268864277464700638583474144061408845077"),
                E: 65537,
            },
            D: fromBase10("7266398431328116344057699379749222532279343923819063639497049039389899328538543087657733766554155839834519529439851673014800261285757759040931985506583861"),
            Primes: []*big.Int{
                fromBase10("98920366548084643601728869055592650835572950932266967461790948584315647051443"),
                fromBase10("94560208308847015747498523884063394671606671904944666360068158221458669711639"),
            },
        }

        hashed := sha256.Sum256(message)

        rsa.SignPKCS1v15(rng, rsaPrivateKey, crypto.SHA256, hashed[:])
        signature, _ := hex.DecodeString("ad2766728615cc7a746cc553916380ca7bfa4f8983b990913bc69eb0556539a350ff0f8fe65ddfd3ebe91fe1c299c2fac135bc8c61e26be44ee259f2f80c1530")

        rsa.VerifyPKCS1v15(&rsaPrivateKey.PublicKey, crypto.SHA256, hashed[:], signature)
    }
    end = time.Now().UnixNano()
    fmt.Println(float64(end - start) / 1e9)
}
