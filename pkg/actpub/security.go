package actpub

type AsPublicKey struct {
	KeyId        string `json:"id"`
	Type         string `json:"type,omitempty"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}
