package main

// TODO 加入字段 时间戳、哈希，删除字段Issuer
type Certificate struct {
	Version        string `json:"version"`
	ABSUID         string `json:"ABSUID"`
	SerialNumber   string `json:"serialNumber"`
	Signature      string `json:"signatureName"`
	Issuer         string `json:"issuer"`
	IssuerCA       string `json:"issuerCA"`
	IssueTime      string `json:"IssueTime"`
	ValidityPeriod string `json:"validityPeriod"`
	ABSAttribute   string `json:"ABSAttribute"`
}

type CertificateResponse struct {
	CertificateContent Certificate `json:"certificate"`
	Hash               string      `json:"hash"`
	ABSSign            string      `json:"absSignature"`
}

type RevokeResponse struct {
	Status       string `json:"status"`
	Timestamp    string `json:"timestamp"`
	Tx           string `json:"tx"`
	SerialNumber string `json:"serialNumber"`
}
