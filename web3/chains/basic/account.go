package basic

type AddressUtil interface {
	// EncodePublicKeyToAddress 将公钥编码为地址
	EncodePublicKeyToAddress(publicKey string) (string, error)
	// DecodeAddressToPublicKey 将地址解码为公钥
	DecodeAddressToPublicKey(address string) (string, error)
	// IsValidAddress 是否是合法地址
	IsValidAddress(address string) bool
}

type Account interface {
	PrivateKey() ([]byte, error)
	PrivateKeyHex() (string, error)
	PublicKey() []byte
	PublicKeyHex() string
	Address() string
	Sign(message []byte, password string) ([]byte, error)
	SignHex(messageHex string, password string) (*OptionalString, error)
}
