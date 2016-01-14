package userAuth
import (
	"github.com/hoysoft/JexGO/sessions"
	"time"
	"encoding/json"
)

const expiryDelta = 10 * time.Second
const (
	keyToken = "oauth2_token"
	PathCallback = "/oauth/callback"
	PathError = "/oauth/error"
)

type IToken interface {
	Access() string
	Refresh() string
	Expired() bool
	ExpiryTime() time.Time
	Valid() bool
}

type token struct {
	Error string `json:"error"`
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken  string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType    string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	Expires_in   time.Duration  `json:"expires_in,omitempty"`
	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry       time.Time `json:"expiry,omitempty"`

	// raw optionally contains extra metadata from the server
	// when updating a token.
	raw          interface{}
}

func (this *token)Access() string {
	return this.AccessToken
}


func (t *token) Expired() bool {
	if t == nil {
		return true
	}
	return !t.Valid()
}

// ExpiryTime returns the expiry time of the user's access token.
func (t *token) ExpiryTime() time.Time {
	return t.Expiry
}

//是否有效
func (t *token) Valid() bool {
	return t!=nil && t.AccessToken != "" && !t.expired()
}

//是否过期
func (t *token)expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Add(-expiryDelta).Before(time.Now())
}

func (t *token)Refresh() string {
	return t.RefreshToken
}

func unmarshallToken(s sessions.Session) (t *token) {
	if s.Get(keyToken) == nil {
		return
	}
	data := s.Get(keyToken).([]byte)
	var tk token
	json.Unmarshal(data, &tk)
	return &tk
}