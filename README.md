# one-eyed-willy

A go restful api to generate pdf from html

# What?

One-Eyed-Willy (OEW) is a simple experiment to create a web app that generates pdf files from a html. The app is named after the pirate of the Goonies movie. When I was a child I loved this movie and prefered pirates over ninjas ;-)

OEW makes use of the `chormedp` package to render the html and generate the pdf file. Currently we don't have a native Go implementation of such tool, from my research `go-weasyprint` is a promising alternative but unfortunately it hasn't reached production level and `wkthmltopdf` is outdated and undermaintained.

# Installation

```
git clone git@github.com:anibaldeboni/one-eyed-willy.git
```

# Unit tests

```
make test
```

# Linter tests

```
make lint
```

# Building

```
make build
```

# Running

```
make run
```

It will start a server on port `8080`

`GET http://localhost:8080/docs` shows the swagger documentation

`POST http://localhost:8080/pdf` generates a pdf file from html base64 encoded string

```
{
    "html": "PGh0bWw+CjxoZWFkPgoJPHRpdGxlPk15IFBERiBGaWxlPC90aXRsZT4KPC9oZWFkPgo8Ym9keT4KCTxwPkhlbGxvIHRoZXJlISBJJ20gYSBwZGYgZmlsZSBnZW5lcmF0ZSBmcm9tIGEgaHRtbCB1c2luZyBnbyBhbmQgZ29wZGYgcGFja2FnZTwvcD4KPC9ib2R5Pgo8L2h0bWw+"
}
```
