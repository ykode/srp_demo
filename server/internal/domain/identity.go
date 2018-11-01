package domain

import (
	"errors"
	"fmt"
	"math/big"
)

type BigIntBytes []byte

func (b BigIntBytes) String() string {
	big := new(big.Int).SetBytes([]byte(b))

	return fmt.Sprintf("0x%x", big)
}

type Identity struct {
	userName string
	salt     BigIntBytes
	verifier BigIntBytes
}

func (id Identity) String() string {
	return fmt.Sprintf("userName: %s; salt: %s; verifier: %s", id.userName, id.salt, id.verifier)
}

func NewIdentity(username string, salt, verifier []byte) (*Identity, error) {

	if len(username) == 0 {
		return nil, errors.New("Username cannot be empty")
	}

	if len(username) < 8 {
		return nil, errors.New("Username should be more than 8 characters.")
	}

	if len(salt) == 0 {
		return nil, errors.New("Salt cannot be empty.")
	}

	if len(verifier) == 0 {
		return nil, errors.New("Verifier cannot be empty.")
	}

	return &Identity{
		userName: username,
		salt:     salt,
		verifier: verifier,
	}, nil
}

func (id *Identity) ChangeUserName(newUserName string) error {
	if len(newUserName) < 8 {
		return errors.New("Username should be more than 8 characters.")
	}

	id.userName = newUserName
	return nil
}

func (id Identity) UserName() string {
	return id.userName
}

func (id Identity) Salt() []byte {
	return id.salt
}

func (id Identity) Verifier() []byte {
	return id.verifier
}
