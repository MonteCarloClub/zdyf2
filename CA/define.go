package main

import "math/big"

type LagPoint struct {
    X *big.Int `json:"x"`
    Y *big.Int `json:"y"`
}

type PublicKey struct {
    G, P, Y *big.Int
}

type PrivateKey struct {
    PublicKey
    X *big.Int
}

type ABSSignature struct {
    C []*big.Int          `json:"c"`
    D []*big.Int          `json:"d"`
    R []*big.Int          `json:"r"`
    LagPoints []*LagPoint `json:"lagpoints"`
}
