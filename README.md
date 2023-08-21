# one-eyed-willy

[![codecov](https://codecov.io/gh/anibaldeboni/one-eyed-willy/graph/badge.svg?token=TGS4JFQ5WV)](https://codecov.io/gh/anibaldeboni/one-eyed-willy)

A web service written in Go to generate pdfs from html and merge multiple files

# What?

One-Eyed-Willy (OEW) is a simple experiment to create a web app that generates pdf files from a html. The app is named after the pirate of the Goonies movie. When I was a child I loved this movie and prefered pirates over ninjas ;-)

OEW makes use of the `chromedp` package to render the html and generate the pdf file. Currently we don't have a native Go implementation of such tool, from my research `go-weasyprint` is a promising alternative but unfortunately it hasn't reached production level yet and `wkthmltopdf` is outdated and undermaintained.
The pdf merge api is provided by `pdfcpu` package.

# Project structure

```
.
├── cmd
│   └── oew
│       └── main.go     // application entrypoint
├── docs                // swagger documentation, automatically generated with `make docs`
├── internal            // application internal packages, not reusable
│   ├── config          // app and echo configurations
│   ├── handler         // handlers are the controllers and usecases
│   └── router          // initialize echo with the application configs
├── pkg                 // public packages, could be reused by other projects
│   ├── logger          // a wrapper around zap logger
│   ├── pdf             // pdf related functions, currently generate from html and merge files
│   └── utils           // generic functions created to makes some tasks simpler
└── testdata            // files and data used in tests
```

# Environment variable

OEW will automaticaly load your application's env vars from `.env.*` files, where `*` is the values of `ENVIRONMENT` var i.e: if `ENVIRONMENT=production` it will look for a `.env.production` file.
If `ENVIRONTMENT` is empty it will load `.env.development` by default.

# Installation

```
git clone git@github.com:anibaldeboni/one-eyed-willy.git
make install_deps
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

You may run the application directly from the source using:

```
make run
```

Or inside a docker container:

```
make docker-build
make docker
```

It will start a server on port `8080`

# API

`POST http://localhost:8080/pdf/generate` generates a pdf file from a html base64 encoded string

`POST http://localhost:8080/pdf/merge` merges two or more pdf files

`POST http://localhost:8080/pdf/encrypt` encrypt a pdf file witha given password

# Frontend

OEW has a web frontend. You may access it in your browser at `http://localhost:8080/`
