package internal

import (
	"net/http"
	"time"
)

//NewCookie gives us a configured cookie to set for the client
func NewCookie(stateCookieName, state string) *http.Cookie {
	//TODO: suitable expiration time
	// expire in 1 month;
	// expiration := time.Now().Add(30 * 24 * time.Hour)
	//or expire in 24 hours
	//expiration := time.Now().Add(24 * time.Hour)
	expiration := time.Now().Add(time.Hour)
	//httpOnly security flag to secure our cookie from XSS
	//HttpOnly: not accessible from Javascript
	cookie := &http.Cookie{
		Name:  stateCookieName,
		Value: state,
		//https only
		// Secure:   true,
		HttpOnly: true,
		//MaxAge: ____,
		// SameSite: http.SameSiteLaxMode,
		Expires: expiration,
	}
	return cookie

}
