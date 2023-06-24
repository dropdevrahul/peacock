![test event parameter](https://github.com/dropdevrahul/peacock/actions/workflows/test.yml/badge.svg?event=push) ![test event parameter](https://github.com/dropdevrahul/peacock/actions/workflows/lint.yml/badge.svg?event=push)

## Gocache
A simple go based Key:Value in-memory cache which uses a custom tcp protocol for the server

## Running the project

```
go run main.go
```

or

```
go build

./gocache

```

* For help about run:
```
./gocache -h  

```

* For changing default host and port use the `-host` and `-port` option

Default port is `9999`

```
./gocache -host 0.0.0.0 -port 1265

```

## Running using Docker

* using the docker-compose

```
docker-compose up -d --build
```

* Image can be generated using the Dockerfile

```
 docker build -t gocache .
```


## Configuring LRU limits

The cache only supports Least Recently Used scheme for removing items once the -max-size limit is reached for cache. It is the maximum number of items that can be stored in the cache.

```
./gocache -maz-size 1000
```

## Client
There is a [client](https://github.com/dropdevrahul/gocacheclient) library to call the server over tcp
