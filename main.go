package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// charset used for random string creation
const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

// configuration
var length = os.Getenv("SHORTURL_LENGTH")
var shortUrlHost = os.Getenv("SHORTURL_HOST")
var redisEnabled, err = strconv.ParseBool(os.Getenv("REDIS_ENABLED"))
var Urls map[string]string

// evaluate seededRand
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// type Url
type Url struct {
	LongUrl  string
	ShortUrl string
}

// shorten function > longUrl > shortUrl
func shortenUrl(u Url) string {
	// normalize shortUrl Host
	if shortUrlHost == "" {
		//default
		shortUrlHost = "http://localhost/"
	}
	// check for trailing slash
	if !strings.HasSuffix(shortUrlHost, "/") {
		shortUrlHost = shortUrlHost + "/"
	}
	// length of string creation
	if length == "" {
		// default
		length = "5"
	}
	iLength, err := strconv.Atoi(length)
	if err != nil {
		log.Fatal(err)
	}
	// define byte size
	b := make([]byte, int(iLength))

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	// while already existent, rerun
	if redisEnabled {
		for len(redisRead(string(b))) != 0 {
			shortenUrl(u)
		}
	} else {
		for len(Urls[string(b)]) != 0 {
			shortenUrl(u)
		}
	}
	return string(b)
}

// index / handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "gursht, your url shortener.")
}

// shorten url and write it to redis
func NewUrlHandler(w http.ResponseWriter, r *http.Request) {
	var u Url

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// short the url
	u.ShortUrl = shortenUrl(u)
	if redisEnabled {
		u = redisWrite(u)
	} else {
		Urls[u.ShortUrl] = u.LongUrl
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = io.WriteString(w, `{"LongUrl":"`+u.LongUrl+`","ShortUrl":"`+shortUrlHost+u.ShortUrl+`"}`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// get the shortened url
func GetUrlHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve vars
	vars := mux.Vars(r)
	// if redis contains this shortUrl
	var longUrl string
	if redisEnabled {
		longUrl = redisRead(vars["url"])
	} else {
		longUrl = Urls[vars["url"]]
	}
	if len(longUrl) > 0 {
		http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "https://"+r.Host, http.StatusTemporaryRedirect)
	}

}

// a very simple health check
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, `{"alive": true}`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {
	// check redis connection if enabled
	if redisEnabled {
		// redis check
		fmt.Println("Trying connection to redis")
		c := redisConn()
		c.Close()
	} else {
		Urls = make(map[string]string)
	}
	// router
	r := mux.NewRouter()
	r.HandleFunc("/healthz", HealthCheckHandler)
	r.HandleFunc("/{url}", GetUrlHandler)
	r.PathPrefix("/").Methods("GET").HandlerFunc(IndexHandler)
	r.PathPrefix("/").Methods("POST").HandlerFunc(NewUrlHandler)
	http.Handle("/", r)

	// listener
	port := os.Getenv("PORT")
	if port == "" {
		// default
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
