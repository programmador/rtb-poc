## Installing dependencies

```
go get -u github.com/labstack/echo/...
go get github.com/bsm/openrtb
go get xojoc.pw/useragent
```


## Running

```
go run server.go
```


## Testing

```
curl -X POST http://localhost:1323 -d @bidrequest.sample.json --header "Content-Type: application/json"
```
