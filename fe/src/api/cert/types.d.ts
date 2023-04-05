declare namespace API {
  type RevokeParams = {
    no: string;
  };

  type RevokeResponse = string | 'Revoke OK.';

  type QueryParams = {
    no: string;
  };

  type ApplyParams = {
    uid: string;
  };

  type Cert = {
    ABSAttribute: string;
    ABSUID: string;
    issuer: string;
    issuerCA: string;
    serialNumber: string;
    signatureName: string;
    validityPeriod: string;
    version: string;
  };

  type QueryResponse = {
    absSignature: string;
    certificate: Cert;
  }

  type VerifyParams = string;

  type VerifyResponse = string;

  type CertOnChainResponse = string | Cert;

  type HistoryParams = {
    index: number;
    count: number;
  };

  type HistoryResponse = {
    certificates: string[];
  };
}
