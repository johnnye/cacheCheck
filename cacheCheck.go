package cacheCheck

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"strings"
)

/*
	TODO: What we need to do here
	is to create our own response writer

	This response writer will wait for responses
	and write them to the cache automaticaly.


*/

var cacheHit bool = false

type Middleware struct {
	hit bool
	c   redis.Conn
	ks  string
}

func NewMiddleware(c redis.Conn, keyspace string) *Middleware {
	return &Middleware{false, c, keyspace}
}

func (l *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	key := strings.Join([]string{l.ks, req.URL.String()}, "")
	exists, err := redis.String(l.c.Do("GET", key))

	if err == nil && len(exists) > 0 {
		log.Printf("Cache hit for %v\n", key)
		cacheHit = true
		//TODO: Set The expire cache header
		w.Write([]byte(exists))
		return
	} else {
		cacheHit = false
		log.Printf("This Was a miss: %v", key)
		next(w, req)
	}
}

/*
	TODO: Set the cache body
	TODO: Seet the cache headers
*/

func SetCache(key string, val string, c redis.Conn) {
	reply, err := redis.String(c.Do("SET", key, val))
	log.Printf("reply: %v error: %v", reply, err)
}

func SetExpire(key string, ttl int, c redis.Conn) {
	reply, err := redis.Int(c.Do("EXPIRE", key, ttl))
	log.Printf("reply: %v error: %v", reply, err)
}

func RemoveCache(key string, c redis.Conn) {
	reply, err := redis.String(c.Do("DEL", key))
	log.Printf("Cache Deleted %v For: %v Error:", reply, key, err)
}
