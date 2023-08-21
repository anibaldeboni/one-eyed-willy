FROM golang:1.21 as build

ENV APP_DIR /src
WORKDIR $APP_DIR

COPY go.sum go.mod $APP_DIR
COPY . $APP_DIR

RUN go mod download
RUN CGO_ENABLED=0 go build -v -o $APP_DIR/bin $APP_DIR/cmd/oew/main.go

# Copy binary to app image.
FROM alpine as app
LABEL maintainer="Aníbal Deboni Neto <anibaldeboni@gmail.com>"
RUN echo @latest https://dl-cdn.alpinelinux.org/alpine/latest-stable/main >> /etc/apk/repositories && \
    echo @latest https://dl-cdn.alpinelinux.org/alpine/latest-stable/community >> /etc/apk/repositories && \
    echo @14.20.1 https://dl-cdn.alpinelinux.org/alpine/v3.18/main >> /etc/apk/repositories

ENV TZ=America/Sao_Paulo

# Install dependencies
RUN apk --no-cache --update add \
    tzdata \
    chromium@latest

WORKDIR /app
COPY --from=build /src/bin /app
COPY --from=build /src/.env.production /src/.env.development /app
EXPOSE 8080
CMD ["./bin"]
