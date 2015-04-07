tmi
--------

A simple library to interface with Twitch's IRC, TMI  
The library is intended to be used as a base for bots, clients, web-relays, overlays, statistics etc.
On top of basic IRC functionality used by TMI, the parsing handles action(/me) messages, parsing of IRCv3 tags, connection timeout handling, automatic ping/pong etc.

THIS IS AN EARLY DRAFT, MANY THINGS ARE SUBJECT TO CHANGE!

The library doesn't explicitly handle commands (it's all passed as strings),
and it doesn't implement any events. Instead it is built on a blocking (buffered) function call, which returns an event or an error. This way you can add on any event system you really want (I have previously used [emission](https://github.com/chuckpreslar/emission) since it's close to the patterns I'm used to), or just chuck it into a loop.

The current version is heavily influenced by [irc-event](https://github.com/thoj/go-ircevent) and it's forks. I will be refactoring and rewriting a lot of it since I'm not entirely convinced with the approach, but it's what I needed to make it work for now.


## Installation

With Google's [Go](http://www.golang.org) installed on your machine:

    $ go get github.com/sunspotseu/tmi

## Usage
See examples

## Documentation

View godocs

    $ godoc github.com/sunspotseu/tmi

## TO DO
- Better handling of errors on the different loops,
cleaner approach to the whole looping insides.
- Lots of other stuff.
- Structure middleware usage.

## License

> The MIT License (MIT)

> Copyright (c) 2015 Sunspots

> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:

> The above copyright notice and this permission notice shall be included in
> all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
> THE SOFTWARE.
