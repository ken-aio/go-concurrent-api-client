# go-concurrent-api-client
Sample for conccurent api request using golang.  

# Library
Using below libraries.  
  
* go-floc
https://github.com/workanator/go-floc
  
* resty
https://github.com/go-resty/resty
  
# How to use
```
$ make run-mock
...
$ make deps
$ make run
```

# run mock api server
```
$ make run-mock
cd mock && dep ensure -v
Gopkg.lock was already in sync with imports and Gopkg.toml
(1/8) Wrote github.com/valyala/fasttemplate@master
(2/8) Wrote github.com/valyala/bytebufferpool@master
(3/8) Wrote github.com/labstack/gommon@0.2.6
(4/8) Wrote github.com/mattn/go-colorable@v0.0.9
(5/8) Wrote github.com/mattn/go-isatty@v0.0.3
(6/8) Wrote github.com/labstack/echo@3.3.5
(7/8) Wrote golang.org/x/crypto@master
(8/8) Wrote golang.org/x/sys@master
cd mock && gin -a 9998 -p 9999 -i --notifications run api-server.go
[gin] Listening on port 9999
[gin] Building...
[gin] Build finished
server start in port: 9998

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v3.3.5
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:9998
```

## Get mock api
See mock dir for detail.  
```
$ curl localhost:9999/api/titles
{
  "titles": [
    {
      "id": 1
    },
    {
      "id": 2
    }
  ]
}
```
```
$ curl localhost:9999/api/titles/1
{
  "id": 1,
  "name": "title1",
  "desc": "this is title1 detail"
}
```
```
$ curl localhost:9999/api/titles/1/episodes
{
  "episodes": [
    {
      "id": 1
    },
    {
      "id": 2
    }
  ]
}
```
```
$ curl localhost:9999/api/titles/1/episodes/1
{
  "id": 1,
  "name": "episode1",
  "desc": "this is episode1 detail"
}
```
