package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rogpeppe/macaroon"
)

// scheme:
//  root    bool
//  else    rw|r
//
// where:
//  rw -> read & write
//  r  -> read only

// e.g.
//  queues    rw
//  messages  r

var (
	auth string
)

func init() {
	// TODO config
	auth = "127.0.0.1:9000"
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hi", Auth(hi))
	http.ListenAndServe(":9999", r)
}

type MacrHandler func(http.ResponseWriter, *http.Request, *macaroon.Macaroon)

func hi(w http.ResponseWriter, r *http.Request, m *macaroon.Macaroon) {
	w.Write([]byte(fmt.Sprintf("hi, your macaroon sig is %s", string(m.Signature()))))
}

func Auth(h MacrHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This is fast, if it's not timed out
		macHdr := r.Header.Get("Macaroon")

		// These will have to talk to another authority to authorize.
		// TODO configurable timeout
		authHdr := r.Header.Get("Authorization")
		// keyHdr := r.Header.Get("Keystone")
		// etc..

		switch {
		case macHdr != "":
			mac := new(macaroon.Macaroon)
			err := mac.UnmarshalJSON([]byte(macHdr))
			if err != nil {
				w.Write([]byte("error deserializing Macaroon"))
				return
			}
			// TODO check for timeout, add discharge to Header
			// TODO check for third party auth
			h(w, r, mac)
		case authHdr != "":
		default:
			w.Write([]byte("no auth supplied"))
			return
		}
	}
}

// should we hatch a scheme in each client that will prefer
// to use macaroons, as well?
//
// i.e. the following heirarchy:
//
//  * client has macaroon for auth
//  * client has oauth token for auth
//  * client has keystone token for auth
//
//  for oauth/keystone:
//    * a macaroon will be discharged on the first request
//    * upon invalidation of macaroon, client must supply

// for clients this complicates things, because they should share
// an http client that uses the same macaroon on each request and
// on macaroon timeout a lot of things could go wrong; emphasis _each_ client.

// alternatively, could we just hand out a macaroon with a 3rd party
// caveat to auth service saying "make sure this user has access to this resource" and
// then the auth service will discharge macaroons that allow short term access
// to the resource?
