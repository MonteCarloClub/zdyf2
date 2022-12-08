package main

type Certificate struct {
    Version string        `json:"version"`
    SerialNumber string   `json:"serialNumber"`
    Signature string      `json:"signatureName"`
    Issuer string         `json:"issuer"`
    ValidityPeriod string `json:"validityPeriod"`
    ABSUID string         `json:"ABSUID"`
}

type CertificateResponse struct {
    CertificateContent Certificate `json:"certificate"`
    ABSSign string                 `json:"absSignature"`
}