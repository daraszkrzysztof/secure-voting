package logic
import (
	"github.com/daraszkrzysztof/secure-voting/secure"
)

type Admin struct {
	Id         string
	PasswdHash string
}

func (adm *Admin)authenticate(admPasswd string) bool {
	return adm.PasswdHash == secure.DoUserPasswdHash(admPasswd)
}