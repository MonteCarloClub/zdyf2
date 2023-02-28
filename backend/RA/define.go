package main

type Certificate struct {
	Version        string `json:"version"`
	SerialNumber   string `json:"serialNumber"`
	Signature      string `json:"signatureName"`
	Issuer         string `json:"issuer"`
	IssuerCA       string `json:"issuerCA"`
	ValidityPeriod string `json:"validityPeriod"`
	ABSUID         string `json:"ABSUID"`
	ABSAttribute   string `json:"ABSAttribute"`
}

type CertificateResponse struct {
	CertificateContent Certificate `json:"certificate"`
	ABSSign            string      `json:"absSignature"`
}

type RevokeResponse struct {
	Status       string `json:"status"`
	Timestamp    string `json:"timestamp"`
	Tx           string `json:"tx"`
	SerialNumber string `json:"serialNumber"`
}
