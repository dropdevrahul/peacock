# A simple go based key:value in memory cache which can be set over tcp using a custom tcp library

## Running the project

```
go run main.go
```

or

```
go build

./gocache

./gocache -h  

./gocache -host 0.0.0.0 -port 1265
```

### Client
There is a [client](https://github.com/dropdevrahul/gocacheclient) library to call the server over tcp
