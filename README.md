# one-eyed-willy

A go restful api to generate pdf from html

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

```
make run
```

It will start a server on port `8080`

`GET http://localhost:8080/docs` shows the swagger documentation

`POST http://localhost:8080/pdf` generates a pdf file from a html base64 encoded string

`POST http://localhost:8080/pdf/merge` merges two or more pdf files
