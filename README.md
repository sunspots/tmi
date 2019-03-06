tmi
--------
**New breaking changes see [CHANGES.md](CHANGES.md)**

A simple library to interface with Twitch's IRC, TMI  
The library is intended to be used as a base for bots, clients, web-relays, overlays, statistics etc.

If you have issues, ideas, feature requests etc. Feel free to open an issue, or send a pull request!

## Features
* Buffered blocking reader
* TCP or Websocket
* TLS support (`secure` option for both TCP/Websocket)
* IRCv3 tags
* Automated PING/PONG
* Decent performance message parsing, >650k messages per second (i7-4600U) (see message_test.go)
* Concurrency-safe message sending for async message handling

## Installation

With Google's [Go](http://www.golang.org) installed on your machine:

    $ go get github.com/sunspots/tmi

## Usage
See examples and documentation

## Documentation

View godocs

    $ godoc github.com/sunspots/tmi

## Tests

    $ go test

Integration tests (connects to Twitch anonymously)

    $ go test -tags=integration

