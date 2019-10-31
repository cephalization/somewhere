# Somewhere

Somewhere is a go program that acts as a Single Page Application server that further proxies requests to another API.

This can be used to avoid CORS errors as the SPA can make requests to this server hosting it while the server itself handles foreign connections.

## Development

- Install go and configure workspace
- `go get -u github.com/cephalization/somewhere`
- `cd $GOPATH/github.com/cephalization/somewhere`

## Usage
```
Usage: somewhere [arguments] directory

  -host string
        host to run server on (default "0.0.0.0")
  -phost string
        target host to proxy api requests to (default "0.0.0.0")
  -port string
        port to run server on (default "8080")
  -pport string
        target port to proxy api requests to (default "8081")
  -prefix string
        route prefix that will be proxied. All other routes will be served the SPA (default "api")
  -pscheme string
        target host scheme to proxy api requests to (ex. 'https://') (default "http://")
 ```
