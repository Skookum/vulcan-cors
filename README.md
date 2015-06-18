# CORS Middleware for Vulcan Proxy
This is an attempt at creating a simple [CORS](http://www.w3.org/TR/cors/) middleware for [Vulcan Proxy](http://vulcanproxy.com)

## Caveats
I hate to start with the negative, but:
* I am pretty new to Go, so there's that
* I am pretty new to Vulcan, so there's that too
* I am scratching an itch, so if my itch didn't touch part of the CORS spec, I didn't scratch it. (A good example is `Access-Control-Allow-Credentials`)

## Install
```
go get github.com/skookum/vulcan-cors
```

## Usage
This presumes you have built new `vulcand` and `vctl` binaries per [the instructions](http://vulcanproxy.com/middlewares.html#example-auth-middleware)
1) Create a YAML file of your allowed hosts and methods:
```
http://google.com:
  - GET
  - POST
http://balls.com:
  - "*"
"*":
  - GET

```
(Notice that to allow anyting use `"*"`. The quotes are necessary. Probably another caveat.)

2) Add the middleware
```
vctl cors upsert -f someFrontend -corsFile=yourYaml.yml
```