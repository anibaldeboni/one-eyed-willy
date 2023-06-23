# one-eyed-willy
A go restful api to generate pdf from html

# What?
This is a simple experiment to create a web app that generates pdf files from a html. The app is named after the pirate of the Goonies movie. When I was a child I loved this movie and prefered pirates over ninjas ;-)

# Installation
```
git clone git@github.com:anibaldeboni/one-eyed-willy.git
make install_deps
```

# Running
```
make run
```
It will start a server on port `8080`

`POST http://localhost:8080/pdf` generates a pdf file from html base64 encoded string

```
{
    "html": "PGh0bWw+CjxoZWFkPgoJPHRpdGxlPk15IFBERiBGaWxlPC90aXRsZT4KPC9oZWFkPgo8Ym9keT4KCTxwPkhlbGxvIHRoZXJlISBJJ20gYSBwZGYgZmlsZSBnZW5lcmF0ZSBmcm9tIGEgaHRtbCB1c2luZyBnbyBhbmQgZ29wZGYgcGFja2FnZTwvcD4KPC9ib2R5Pgo8L2h0bWw+"
}
```