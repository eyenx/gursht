# gursht - a simple URL shortener written in GO

I always wanted to write something in Golang, and this is my first real mini project trying this it out.

# Still a WIP 

## getting started


### build & run it manually

```
go build 
./gursth
```

### Docker 

This mini URL Shortener app uses Redis to keep state. You can use it directly with the provided `deploy/docker-compose.yml` or with this oneliner:

```
docker run -n gursht_redis -d redis
docker run -p 3000:3000  -e REDIS_HOST=gursth_redis ghcr.io/eyenx/gursht:latest
```

Access http://localhost and read the How-To


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

This will save the created short url into redis by mapping it with the long url as value.

Now, accessing your host with the newly created short url, will redirect you to the long one:

```
curl http://localhost:3000/2eCRU
<a href="https://example.com/my/very/long/url">Moved Permanently</a>.
```


## Additional configuration

configuration is done by setting environment variables:

```
REDIS_PORT # set the redis port to use
REDIS_HOST # set the redis host to use
SHORTURL_HOST # set the external short url hostname 
SHORTURL_LENGTH # set the length of the created random short path (default: 5)
```

## TODO

* Make it possible to provide the short url on request 
* Already saved longurl shouldn't be mapped to a new shorturl (might not work with redis as backend)
* index.html should provide a Getting started
* GitHub Workflow to autobuild the image
* `deploy/` folder providing docker-compose & kubernetes Yaml files
* Helm Chart?
