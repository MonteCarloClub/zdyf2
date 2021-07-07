package main

import (
    "crypto/sha256"
    "flag"
    "math/big"
)

// This is the 1024-bit MODP group from RFC 5114, section 2.1:
const primeHex = "B10B8F96A080E01DDE92DE5EAE5D54EC52C99FBCFB06A3C69A6A9DCA52D23B616073E28675A23D189838EF1E2EE652C013ECB4AEA906112324975C3CD49B83BFACCBDD7D90C4BD7098488E9C219A73724EFFD6FAE5644738FAA31A4FF55BCCC0A151AF5F0DC8B4BD45BF37DF365C1A65E68CFDA76D4DA708DF1FB2BC2E4A4371"
const generatorHex = "A4D1CBD5C3FD34126765A442EFB99905F8104DD258AC507FD6406CFF14266D31266FEA1E5C41564B777E690F5504F213160217B4B01B886A5E91547F9E2749F4D7FBD7D3B9A92EE1909D0D2263F80A76A6A24C087A091F531DBF0A0169B6A28AD662A4D18E73AFA32D779D5918D08BC8858F4DCEF97C2A24855E6EEB22B3B2E5"

func fromHex(hex string) *big.Int {
    n, ok := new(big.Int).SetString(hex, 16)
    if !ok {
        panic("failed to parse hex number")
    }
    return n
}

var (
    numT *int
    numN *int
    numG *big.Int
    numP *big.Int
    priv []*PrivateKey
)

func init() {
    numT = flag.Int("t", 75, "t of (t, n).")
    numN = flag.Int("n", 100, "n of (t, n).")
    flag.Parse()
    numG = fromHex(generatorHex)
    numP = fromHex(primeHex)

    // p（P）阶循环群 G，g（G）为生成元，s_i（X）为私钥
    for i := 1; i <= *numN; i += 1 {
        priv = append(priv, &PrivateKey{
            PublicKey: PublicKey{
                G: fromHex(generatorHex),
                P: fromHex(primeHex),
            },
            X: big.NewInt(int64(i)),
        })
    }

    for _, key := range priv {
        key.Y = new(big.Int).Exp(key.G, key.X, key.P)
    }
}

func Generate(m string) *ABSSignature {
    M := []byte(m)

    // (t, n) 门限
    var R []*big.Int
    var T []*big.Int
    // 属于属性的 1~t
    for i := 1; i <= *numT; i += 1 {
        T = append(T, big.NewInt(10 + int64(i)))
    }
    for i, t := range T {
        R = append(R, new(big.Int).Exp(priv[i].G, t, priv[i].P))
    }

    // 不属于属性的 t+1～n
    var C []*big.Int
    var D []*big.Int
    for i := 1; i <= *numN - *numT; i += 1 {
        C = append(C, big.NewInt(30 + int64(i)))
        D = append(D, big.NewInt(40 + int64(i)))
    }
    for i, c := range C {
        R = append(R, new(big.Int).Mul(new(big.Int).Exp(numG, D[i], numP), new(big.Int).Exp(priv[i + *numT].Y, c, numP)))
    }

    buf := M
    for _, r := range R {
        buf = append(buf, r.Bytes()...)
    }
    result := sha256.Sum256(buf)
    resultTemp := result[:]
    // fmt.Printf("%x", result)

    lagPoints := []*LagPoint {
        {
            X: big.NewInt(0),
            Y: new(big.Int).SetBytes(resultTemp),
        },
    }
    for i := 1; i <= *numN - *numT; i += 1 {
        lagPoints = append(lagPoints, &LagPoint{
            X: big.NewInt(int64(*numT + i)),
            Y: C[i - 1],
        })
    }

    var CTemp []*big.Int
    var DTemp []*big.Int
    for i := 1; i <= *numT; i += 1 {
        cTemp := LagRange(lagPoints, big.NewInt(int64(i)))
        CTemp = append(CTemp, cTemp)
        dTemp := new(big.Int).Sub(T[i - 1], new(big.Int).Mul(cTemp, priv[i - 1].X))
        DTemp = append(DTemp, dTemp)
    }

    C = append(CTemp, C...)
    D = append(DTemp, D...)
    //for i, c := range C {
    //   c.Mod(c, numP)
    //   D[i] = new(big.Int).Mod(D[i], numP)
    //   R[i] = new(big.Int).Mod(R[i], numP)
    //}
    //fmt.Println(C)
    //fmt.Println(D)
    //fmt.Println(R)
    return &ABSSignature{
        C: C,
        D: D,
        R: R,
        LagPoints: lagPoints,
    }
}

func Verify(signature *ABSSignature) bool {
    C := signature.C
    D := signature.D
    R := signature.R
    lagPoints := signature.LagPoints

    // 验证
    for i, cT := range C {
        res := cT.Cmp(LagRange(lagPoints, big.NewInt(int64(i+1))))
        if res != 0 {
            return false
        }
    }

    for i, rT := range R {
        res := rT.Mod(rT, numP).Cmp(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(numG, D[i], numP), new(big.Int).Exp(priv[i].Y, C[i], numP)), numP))
        if res != 0 {
            return false
        }
    }

    return true
    // test
    //rt1 := new(big.Int).Exp(numG, big.NewInt(12), numP)
    //st := fromHex("40")
    //yt := new(big.Int).Exp(numG, st, numP)
    //ct := C[0]
    //tt := big.NewInt(12)
    //dt := new(big.Int).Sub(tt, new(big.Int).Mul(ct, st))
    //rt2 := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(numG, dt, numP), new(big.Int).Exp(yt, ct, numP)), numP)
    //fmt.Println(rt1)
    //fmt.Println(rt2)
    //
    //fmt.Println(rt1.Cmp(rt2))
}