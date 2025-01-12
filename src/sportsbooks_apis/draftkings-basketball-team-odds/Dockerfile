FROM golang:latest

WORKDIR /src

COPY get_data.go /src/

RUN go mod init get-data-script

RUN go get github.com/gocolly/colly/v2

RUN go mod tidy

CMD ["go", "run", "get_data.go"]
