# Somewhere

Somewhere is a go program that acts as an SPA server that further proxies request to another API.

This can be used to avoid CORS errors as the SPA can make requests to this server hosting it while the server itself handles foreign connections.

## Development

- Install go and configure workspace
- `go get -u github.com/cephalization/somewhere`
- `cd $GOPATH/github.com/cephalization/somewhere`
