package encodable

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vmihailenco/msgpack"

	"github.com/onflow/flow-go/crypto"
)

func toHex(bs []byte) string {
	return fmt.Sprintf("%x", bs)
}

func fromHex(b []byte) ([]byte, error) {
	var x string
	if err := json.Unmarshal(b, &x); err != nil {
		return nil, err
	}
	return hex.DecodeString(x)
}

// NetworkPubKey wraps a public key and allows it to be JSON encoded and decoded. It is not defined in the
// crypto package since the crypto package should not know about the different key types.
type NetworkPubKey struct {
	crypto.PublicKey
}

func (pub NetworkPubKey) MarshalJSON() ([]byte, error) {
	if pub.PublicKey == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(toHex(pub.Encode()))
}

func (pub *NetworkPubKey) UnmarshalJSON(b []byte) error {
	bz, err := fromHex(b)
	if err != nil {
		return err
	}

	if len(bz) == 0 {
		return nil
	}

	pub.PublicKey, err = crypto.DecodePublicKey(crypto.ECDSAP256, bz)
	return err
}

// NetworkPrivKey wraps a private key and allows it to be JSON encoded and decoded. It is not defined in the
// crypto package since the crypto package should not know about the different key types. More importantly, private
// keys should not be automatically encodable/serializable to prevent accidental secret sharing. The bootstrapping
// package is an exception, since it generates private keys that need to be serialized.
type NetworkPrivKey struct {
	crypto.PrivateKey
}

func (priv NetworkPrivKey) MarshalJSON() ([]byte, error) {
	if priv.PrivateKey == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(toHex(priv.Encode()))
}

func (priv *NetworkPrivKey) UnmarshalJSON(b []byte) error {
	bz, err := fromHex(b)
	if err != nil {
		return err
	}

	if len(bz) == 0 {
		return nil
	}
	priv.PrivateKey, err = crypto.DecodePrivateKey(crypto.ECDSAP256, bz)
	return err
}

// StakingPubKey wraps a public key and allows it to be JSON encoded and decoded. It is not defined in the
// crypto package since the crypto package should not know about the different key types.
type StakingPubKey struct {
	crypto.PublicKey
}

func (pub StakingPubKey) MarshalJSON() ([]byte, error) {
	if pub.PublicKey == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(toHex(pub.Encode()))
}

func (pub *StakingPubKey) UnmarshalJSON(b []byte) error {
	bz, err := fromHex(b)
	if err != nil {
		return err
	}

	if len(bz) == 0 {
		return nil
	}
	pub.PublicKey, err = crypto.DecodePublicKey(crypto.BLSBLS12381, bz)
	return err
}

// StakingPrivKey wraps a private key and allows it to be JSON encoded and decoded. It is not defined in the
// crypto package since the crypto package should not know about the different key types. More importantly, private
// keys should not be automatically encodable/serializable to prevent accidental secret sharing. The bootstrapping
// package is an exception, since it generates private keys that need to be serialized.
type StakingPrivKey struct {
	crypto.PrivateKey
}

func (priv StakingPrivKey) MarshalJSON() ([]byte, error) {
	if priv.PrivateKey == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(toHex(priv.Encode()))
}

func (priv *StakingPrivKey) UnmarshalJSON(b []byte) error {
	bz, err := fromHex(b)
	if err != nil {
		return err
	}

	if len(bz) == 0 {
		return nil
	}
	priv.PrivateKey, err = crypto.DecodePrivateKey(crypto.BLSBLS12381, bz)
	return err
}

// RandomBeaconPubKey wraps a public key and allows it to be JSON encoded and decoded. It is not defined in the
// crypto package since the crypto package should not know about the different key types.
type RandomBeaconPubKey struct {
	crypto.PublicKey
}

func (pub RandomBeaconPubKey) MarshalJSON() ([]byte, error) {
	if pub.PublicKey == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(toHex(pub.Encode()))
}

func (pub *RandomBeaconPubKey) UnmarshalJSON(b []byte) error {
	bz, err := fromHex(b)
	if err != nil {
		return err
	}

	if len(bz) == 0 {
		return nil
	}
	pub.PublicKey, err = crypto.DecodePublicKey(crypto.BLSBLS12381, bz)
	return err
}

func (pub RandomBeaconPubKey) MarshalMsgpack() ([]byte, error) {
	if pub.PublicKey == nil {
		return nil, fmt.Errorf("empty public key")
	}
	return msgpack.Marshal(toHex(pub.PublicKey.Encode()))
}

func (pub *RandomBeaconPubKey) UnmarshalMsgpack(b []byte) error {
	bz, err := fromHex(b)
	if err != nil {
		return err
	}

	if len(bz) == 0 {
		return nil
	}
	pub.PublicKey, err = crypto.DecodePublicKey(crypto.BLSBLS12381, bz)
	return err
}

func (pub *RandomBeaconPubKey) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, pub.PublicKey.Encode())
}

// RandomBeaconPrivKey wraps a private key and allows it to be JSON encoded and decoded. It is not defined in
// the crypto package since the crypto package should not know about the different key types. More importantly, private
// keys should not be automatically encodable/serializable to prevent accidental secret sharing. The bootstrapping
// package is an exception, since it generates private keys that need to be serialized.
type RandomBeaconPrivKey struct {
	crypto.PrivateKey
}

func (priv RandomBeaconPrivKey) MarshalJSON() ([]byte, error) {
	if priv.PrivateKey == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(toHex(priv.Encode()))
}

func (priv *RandomBeaconPrivKey) UnmarshalJSON(b []byte) error {
	bz, err := fromHex(b)
	if err != nil {
		return err
	}

	if len(bz) == 0 {
		return nil
	}
	priv.PrivateKey, err = crypto.DecodePrivateKey(crypto.BLSBLS12381, bz)
	return err
}
