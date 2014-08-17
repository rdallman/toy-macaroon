package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rdallman/gorocksdb"
	"github.com/rogpeppe/macaroon"
	"gopkg.in/mgo.v2/bson"
)

type Auth struct {
	db *gorocksdb.DB
}

type Project struct {
	id     bson.ObjectId
	userid bson.ObjectId
	shared map[bson.ObjectId]struct{}
}

const (
	r = 1 << iota
	rw
	pk // peek

	queue = iota
	message

	secret = "youwillnevergetthis"
)

var (
	datadir = "datadir"
	authdb  *Auth
)

func init() {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, datadir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	authdb = &Auth{db: db}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth", auth).Methods("GET")
	r.HandleFunc("/new", newMac).Methods("GET")
	http.ListenAndServe(":9000", r)
}

func newMac(w http.ResponseWriter, r *http.Request) {
	m, err := macaroon.New([]byte(secret), "1", "example.com")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	b, err := m.MarshalJSON()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}

// this should be a caveat verifier
func auth(w http.ResponseWriter, r *http.Request) {
	macHdr := r.Header.Get("Macaroon")
	authHdr := r.Header.Get("Authorization")

	switch {
	case macHdr != "":
		mac := new(macaroon.Macaroon)
		err := mac.UnmarshalBinary([]byte(macHdr))
		if err != nil {
			w.Write([]byte("error deserializing Macaroon"))
			return
		}
		// TODO check for timeout, add discharge to Header
		// TODO check for third party auth
	case authHdr != "":
	default:
		w.Write([]byte("no auth supplied"))
		return
	}
}

func (a *Auth) authenticate(r *http.Request) {
}
