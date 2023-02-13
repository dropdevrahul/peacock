From golang:1.20.0-buster

WORKDIR /usr/src/app


COPY go.mod go.sum ./

COPY . .

RUN go build -o /usr/local/bin/app


CMD ["app"]
