package cryptography

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/harbourrocks/harbour/pkg/context"
	"io/ioutil"
)

// GenerateX5C returns a x5c signature of the certificate
//  the certificate is expected to be PEM encoded
func GenerateX5C(ctx context.HRock, certificatePath string) (x5c []string, err error) {
	log := ctx.L

	log.WithField("certPath", certificatePath).Trace("Generating x5c signature")

	// read certificate file
	certBytes, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		log.WithError(err).WithField("certPath", certificatePath).Error("Failed to open certificate file")
		return
	}

	// generate x5c, we first have to convert from pem to der format, otherwise x509.ParseCertificates throws an error
	block, _ := pem.Decode(certBytes)
	certificates, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		log.WithError(err).WithField("certPath", certificatePath).Error("Failed to parse certificate file")
		return
	}

	// generate the x5c header which is simple a string array of the certificate chain
	// each certificate is represented by a base64 encoded string
	x5c = make([]string, len(certificates))
	for i, certificate := range certificates {
		x5c[i] = base64.StdEncoding.EncodeToString(certificate.Raw)
	}

	log.Trace("x5c generated")

	return
}
