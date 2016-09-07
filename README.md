# sample-golang-graceful-restart

Sample Golang net/http graceful restart with https://github.com/shogo82148/go-gracedown


### install dependencies

```
go get -u github.com/mattn/gom
gom install
```

### build & start server (terminal 1)

```
go build -o bin/server
./bin/server
2016/09/07 13:47:05 Server PID: 2625
```

### access to the url that wait 10 secs using curl (terminal 2)

```
curl http://localhost:8080/count/wait
```

### or using wrk (terminal 2)

```
wrk --timeout 15s http://localhost:8080/count/wait
```

### send SIGHUP (in this example) to the server process (terminal 3)

```
kill -SIGHUP 2625
```
