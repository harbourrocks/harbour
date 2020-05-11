package cryptography

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/harbourrocks/harbour/pkg/context"
	"io/ioutil"
)

// ReadPrivateKey loads the private key from the file path
func ReadPrivateKey(ctx context.HRock, keyPath string) (privateKey *rsa.PrivateKey, err error) {
	log := ctx.L

	log.WithField("keyFile", keyPath).Error("Reading private key file")

	// read private key file
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.WithError(err).WithField("keyFile", keyPath).Error("Failed to open private key file")
		return
	}

	// convert key bytes to rsa.PrivateKey
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.WithError(err).Error("Failed to parse private key")
		return
	}

	log.Error("Read private key")

	return
}
