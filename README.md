[![Go Report Card](https://goreportcard.com/badge/github.com/sauravhiremath/ptransfer?style=for-the-badge)](https://goreportcard.com/report/github.com/sauravhiremath/ptransfer)
[![Lines of code](https://img.shields.io/tokei/lines/github/sauravhiremath/ptransfer?style=for-the-badge)](https://github.com/sauravhiremath/ptransfer)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sauravhiremath/ptransfer?style=for-the-badge)](https://github.com/golang/go)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/sauravhiremath/ptransfer?style=for-the-badge)](https://github.com/sauravhiremath/ptransfer)

![image](https://user-images.githubusercontent.com/28642011/95137895-a9c1e200-0786-11eb-84de-281ef9c05874.png)

# ptransfer

*File transfer over parallel TCP, requires port-forwarding your TCP ports.*

Heavily inspired from [magic-wormhole](https://github.com/warner/magic-wormhole) and [wormhole](https://github.com/schollz/wormhole) except it doesn't have the rendevous server, or the transit relay. It's not really anything like it, except that its file transfer over TCP with AES-128 encryption built over it. Here you can transfer a file using multiple TCP ports simultaneously. 

## Normal use

### Server computer 

Be sure to open up TCP ports 27001-27009 on your port forwarding. Also, get your public address:

```
$ curl icanhazip.com
X.Y.W.Z
```

Then get and run *ptransfer* with a 16 byte key for encryption/password:

```
$ go get github.com/sauravhiremath/ptransfer
$ ptransfer -file SOMEFILE -p ABCDEFGHIJKLMNOP
```

*ptransfer* automatically knows to run as a server when the `-file` flag is set.

### Client computer

```
$ go get github.com/sauravhiremath/ptransfer
$ ptransfer -server X.Y.W.Z -p ABCDEFGHIJKLMNOP
```

*ptransfer* automatically knows to run as a client when the `-server` flag is set.


## Building for use without flags

For people that don't have or don't want to build from source and don't want to use the command line, you can build it for them to have the flags set automatically! Build the ptransfer binary so that it always behaves as a client to a specified server, so that someone just needs to click on it.

```
go build -ldflags "-s -w -X main.serverAddress=X.Y.W.Z" -o client.exe
```

Likewise you could do the same for the server:

```
go build -ldflags "-s -w -X main.fileName=testfile" -o server.exe
```

# Encryption AES 128 Usage

TODO

# Development Setup

> Note: This setup uses go1.15.2

Installing the dependencies in vendor

```
go mod vendor
```

Build the project to get the executable

```
go build -o build/ptransfer ./cmd/ptransfer
```

Run the project directly without `build`

```
go run cmd/ptransfer/main.go
```

Creating Releases (adding new tags, following semVer)

```
git tag v0.1.0
git push --tags
```

# License

MIT
