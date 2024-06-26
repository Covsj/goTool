package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"golang.org/x/crypto/blake2b"
)

func IsValidSignature(publicKey, msg, signature []byte) bool {
	if len(msg) > 256 {
		h := blake2b.Sum256(msg)
		msg = h[:]
	}
	var (
		sigs        [64]byte
		fixedPubKey [32]byte
		sig         = new(schnorrkel.Signature)
	)
	copy(fixedPubKey[:], publicKey[:])
	copy(sigs[:], signature[:])
	pubKey, err := schnorrkel.NewPublicKey(fixedPubKey)
	if err != nil {
		return false
	}
	if err := sig.Decode(sigs); err != nil {
		return false
	}
	res, err := pubKey.Verify(sig, schnorrkel.NewSigningContext([]byte("substrate"), msg))
	if err != nil {
		return false
	}
	return res
}
