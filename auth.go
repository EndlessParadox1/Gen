package gen

import (
	"encoding/base64"
	"net/http"
	"strconv"
)

const AuthUserKey = "user"

const AuthProxyUserKey = "proxy_user"

type Accounts map[string]string

type authPair struct {
	user  string
	value string
}

// BasicAuthWithRealm If the realm is empty, "Authorization Required" will be used by default.
func BasicAuthWithRealm(accounts Accounts, realm string) HandlerFunc {
	if realm == "" {
		realm = "Authorization Required"
	}
	realm = "Basic realm=" + strconv.Quote(realm)
	pairs := processAccounts(accounts)
	return func(c *Context) {
		user, found := searchCredential(pairs, c.Request.Header.Get("Authorization"))
		if !found {
			c.SetHeader("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Set(AuthUserKey, user)
		}
	}
}

// BasicAuth returns a Basic HTTP Authorization middleware.
func BasicAuth(accounts Accounts) HandlerFunc {
	return BasicAuthWithRealm(accounts, "")
}

// BasicAuthForProxyWithRealm If the realm is empty, "Proxy Authorization Required" will be used by default.
func BasicAuthForProxyWithRealm(accounts Accounts, realm string) HandlerFunc {
	if realm == "" {
		realm = "Proxy Authorization Required"
	}
	realm = "Basic realm=" + strconv.Quote(realm)
	pairs := processAccounts(accounts)
	return func(c *Context) {
		user, found := searchCredential(pairs, c.Request.Header.Get("Proxy-Authorization"))
		if !found {
			c.SetHeader("Proxy-Authenticate", realm)
			c.AbortWithStatus(http.StatusProxyAuthRequired)
		} else {
			c.Set(AuthProxyUserKey, user)
		}
	}
}

// BasicAuthForProxy returns a Basic HTTP Proxy Authorization middleware.
func BasicAuthForProxy(accounts Accounts) HandlerFunc {
	return BasicAuthForProxyWithRealm(accounts, "")
}

func processAccounts(accounts Accounts) []authPair {
	var pairs []authPair
	for user, password := range accounts {
		base := user + ":" + password
		value := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
		pairs = append(pairs, authPair{user, value})
	}
	return pairs
}

func searchCredential(a []authPair, authValue string) (user string, found bool) {
	if len(authValue) == 0 {
		return
	}
	for _, pair := range a {
		if pair.value == authValue {
			return pair.user, true
		}
	}
	return
}
