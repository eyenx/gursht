# gursht - a simple URL shortener written in GO

![](https://github.com/eyenx/gursht/workflows/lint-and-test/badge.svg)
![](https://github.com/eyenx/gursht/workflows/release/badge.svg)
![](https://github.com/eyenx/gursht/workflows/qa/badge.svg)

**Still a WIP**

I always wanted to write something in Golang, and this is my first real mini project trying it out.

My main motivation behind this was getting a url shortener for my mutt configuration. [z3bra](http://z3bra.org) uses a url shortener too, and I wanted to implement this too. The main idea behind it, is getting rid of mutt's `+` line wraps when a link gets too long for the term.

Find a demo of the app running on https://s.uff.li. Feel free to use it.

## Getting started


### build & run it manually

```
git clone https://github.com/eyenx/gursht && cd gursht
go build
./gursht
```

### Docker

You can use it directly with the provided `deploy/docker-compose.yml` or with this oneliner:

#### In memory state

```
docker run -n gursht -p 3000:3000 ghc.rio/eyenx/gursht:latest
```

#### Example with Redis
```
docker run -n gursht_redis -d redis
docker run -p 3000:3000 -e REDIS_ENABLED=true -e REDIS_HOST=gursht_redis ghcr.io/eyenx/gursht:latest
```

Access http://localhost:3000 and read the howto.


### Kubernetes

tbd

### Helm

tbd

## Usage

a simple request looks like this:

```
curl -d '{"LongUrl":"https://example.com/my/very/long/url"}' localhost:3000 -H "Content-Type: application/json"
{"LongUrl":"https://example.com/my/very/long/url","ShortUrl":"http://localhost/2eCRU"}%
```

This will save the created short url by mapping it with the long url as value.

Now, accessing your host with the newly created short url, will redirect you to the long one:

```
curl http://localhost:3000/2eCRU
<a href="https://example.com/my/very/long/url">Moved Permanently</a>.
```


## Additional configuration

configuration is done by setting environment variables:

```
REDIS_ENABLED # set to true if you wanna use REDIS, default is inmemory go map
REDIS_PORT # set the redis port to use
REDIS_HOST # set the redis host to use
SHORTURL_HOST # set the external short url hostname
SHORTURL_LENGTH # set the length of the created random short path (default: 5)
```

## TODO

* Make it possible to provide the short url on request
* Already saved longurl shouldn't be mapped to a new shorturl (might not work with redis as backend)
* index.html should provide a Getting started
* `deploy/` folder providing docker-compose & kubernetes Yaml files
* Helm Chart?
