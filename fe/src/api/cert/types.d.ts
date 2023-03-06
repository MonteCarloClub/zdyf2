declare namespace API {
  type RevokeParams = {
    no: string;
  };

  type RevokeResponse = string | 'Revoke OK.';

  type InfoParams = {
    uid: string;
  };

  type ApplyParams = {
    uid: string;
    attribute: string;
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
}
