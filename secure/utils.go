package secure
import (
	"crypto/sha512"
	"encoding/base64"
)

func DoUserPasswdHash(input string) string{

	hasher := sha512.New()
	hasher.Write([]byte(input))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}