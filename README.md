# Gonews
[![Build Status](https://travis-ci.org/drgarcia1986/gonews.svg)](https://travis-ci.org/drgarcia1986/gonews)

A Golang CLI to browser news from [Hacker News](https://news.ycombinator.com/)  
Also is a Golang implementation of the awesome [pynews-cli](https://github.com/mazulo/pynews_cli).

## Installing
Get this project and install via `go get -u`.
```
$ go get -u github.com/drgarcia1986/gonews
```
Or you can download a binary on releases.

## Usage example
To get the last 10 `new` stories just call `gonews` on command line and press `Enter` on the news wanted.

You can determine the number max of news to return with parameter `--limit` and you can also choose between `new` or `top` stories
with parameter `--type`.
```
$ gonews --limit 20 --type top
```
For more information about usage, use command `gonews --help`
```
$ gonews --help
Usage of gonews:
  -limit int
        Number of Stories to get (default 10)
  -type string
        Stories Type ('new' or 'top') (default "top")
```

## Observation
This is a toy project to help me to pratice Golang, if you can help me with this, getting in touch :smile:.
