package secure
import (
	"crypto/sha512"
	"encoding/base64"
	"crypto/rsa"
	"crypto/rand"
	"errors"
	"crypto/x509"
	"encoding/pem"
)

func DoUserPasswdHash(input string) string{

	hasher := sha512.New()
	hasher.Write([]byte(input))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}


func GenerateRSAKeyPair() (privateKey *rsa.PrivateKey, err error) {

	privateKey, err = rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		return nil, errors.New("Problem while generating RSA key pair")
	}

	err = privateKey.Validate()

	if err != nil {
		return nil, errors.New("Problem while validating RSA key pair")
	}

	return privateKey, err
}

/*
private key returned in DER format
public key returned in PEM format
*/
func MarshalToX509KeyPair(priv *rsa.PrivateKey) (priv_der []byte,pub_pem []byte, err error){

	priv_der = x509.MarshalPKCS1PrivateKey(priv)
	pub_der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)

	if err != nil {
		return nil,nil, errors.New("Problem while marshalling public key")
	}

	pub_blk := pem.Block{
		Type: "PUBLIC KEY",
		Headers: nil,
		Bytes: pub_der,
	}

	pub_pem = pem.EncodeToMemory(&pub_blk)

	return priv_der, pub_pem, nil
}