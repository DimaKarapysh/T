package bootstrap

import "T/internal/token"

func NewToken(accessToken *token.AccessToken) token.Token {
	return accessToken.Token
}
