package userAuth

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/hoysoft/JexGO/render"
	"github.com/hoysoft/JexGO/sessions"
	"log"
	"net/http"
)

// These are the default configuration values for this package. They
// can be set at anytime, probably during the initial setup of Martini.
var (
// RedirectUrl should be the relative URL for your login route
	RedirectUrl string = "/login"



// SessionKey is the key containing the unique ID in your session
	SessionKey string = "AUTHUNIQUEID"
)


// User defines all the functions necessary to work with the user's authentication.
// The caller should implement these functions for whatever system of authentication
// they choose to use
type IUser interface {
	// Return whether this user is logged in or not
	IsAuthenticated() bool

	// Set any flags or extra data that should be available
	Login()

	// Clear any sensitive data out of the user
	Logout()

	// Return the unique identifier of this user object
	UniqueId() interface{}


	// Populate this user object with values
	GetById(id interface{}) error

	GetByToken(tokenAccess string) error

	IsAdmin() bool
}


// SessionUser will try to read a unique user ID out of the session. Then it tries
// to populate an anonymous user object from the database based on that ID. If this
// is successful, the valid user is mapped into the context. Otherwise the anonymous
// user is mapped into the contact.
// The newUser() function should provide a valid 0value structure for the caller's
// user type.
func SessionUser(newUser func() IUser) martini.Handler {
	return func(s sessions.Session, c martini.Context, l *log.Logger) {
				userId := s.Get(SessionKey)
				user := newUser()

				if userId != nil {
					err := user.GetById(userId)

					if err != nil {
                        s.Clear()
						user = newUser()
						l.Printf("Login Error0: %v\n", err)
					} else {
						user.Login()
					}
				}

		c.MapTo(user, (*IUser)(nil))
	}
}

//func SessionOAuthUser(newUser func() IUser) martini.Handler {
//	return func(s sessions.Session, c martini.Context, l *log.Logger) {
//		user := newUser()
//		token :=  unmarshallToken(s)
//		if token.Valid() {
//			err := user.GetByToken(token.Access())
//			if err != nil {
//				l.Printf("Login Error: %v\n", err)
//			} else {
//				user.Login()
//			}
//		}
//
//		c.MapTo(user, (*IUser)(nil))
//	}
//}

//oauth登陆成功的回调
func oAuthUserLoginCallback(s sessions.Session, user IUser, l *log.Logger, w http.ResponseWriter, req *http.Request) {
	token :=  unmarshallToken(s)
	if token.Valid() {
		err := user.GetByToken(token.Access())
		if err != nil {
			 l.Printf("Login Error1: %v\n", err,token.Access())

		} else {
			err := AuthenticateSession(s, user)
			if err != nil {
				http.Error(w,err.Error(),500)
			}
//			params := req.URL.Query()
//			redirect := params.Get(RedirectParam)
//			http.Redirect(w, req, redirect, 302)
		}
	}else{

		http.Redirect(w, req, RedirectUrl, 302)
	}
}





// AuthenticateSession will mark the session and user object as authenticated. Then
// the Login() user function will be called. This function should be called after
// you have validated a user.
func AuthenticateSession(s sessions.Session, user IUser) error {
	user.Login()
	return UpdateUser(s, user)
}

// Logout will clear out the session and call the Logout() user function.
func Logout(s sessions.Session, user IUser) {
	user.Logout()
	s.Delete(SessionKey)
}

// LoginRequired verifies that the current user is authenticated. Any routes that
// require a login should have this handler placed in the flow. If the user is not
// authenticated, they will be redirected to /login with the "next" get parameter
// set to the attempted URL.
func LoginRequired(s sessions.Session,r render.Render, user IUser, req *http.Request) {
	if user.IsAuthenticated() == false {
		s.Delete(SessionKey)
		path := fmt.Sprintf("%s?%s=%s", RedirectUrl, RedirectParam, req.URL.Path)
		r.Redirect(path, 302)
	}
}
//
//func LoginRequiredOAuthUser(r render.Render, user IUser, token IToken,req *http.Request) {
//	if user.IsAuthenticated() == false || token.Expired() {
//		user.Logout()
//		path := fmt.Sprintf("%s?%s=%s", RedirectUrl, RedirectParam, req.URL.Path)
//		r.Redirect(path, 302)
//	}
//}

// UpdateUser updates the User object stored in the session. This is useful incase a change
// is made to the user model that needs to persist across requests.
func UpdateUser(s sessions.Session, user IUser) error {
    s.Set(SessionKey, user.UniqueId())
	return nil
}


