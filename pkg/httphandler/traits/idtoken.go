package traits

import (
	"github.com/harbourrocks/harbour/pkg/auth"
)

// IdTokenTrait returns an idToken to the caller
type IdTokenTrait interface {
	GetToken() *auth.IdToken
	SetToken(t auth.IdToken)
	GetTokenStr() string
	SetTokenStr(string)
	HttpTrait
}

// IdTokenModel holds the idToken
type IdTokenModel struct {
	idToken    *auth.IdToken
	idTokenStr string
}

func (m IdTokenModel) GetToken() *auth.IdToken {
	return m.idToken
}

func (m *IdTokenModel) SetToken(t auth.IdToken) {
	m.idToken = &t
}

func (m IdTokenModel) GetTokenStr() string {
	return m.idTokenStr
}

func (m *IdTokenModel) SetTokenStr(t string) {
	m.idTokenStr = t
}

func AddIdToken(trait IdTokenTrait) {
	r := trait.GetRequest()

	token, tokenStr, err := auth.HeaderAuth(r, trait.GetOidcConfig())
	if err != nil {
		return
	}

	idToken, err := auth.IdTokenFromToken(token)
	if err != nil {
		return
	}

	trait.SetToken(idToken)
	trait.SetTokenStr(tokenStr)
}
