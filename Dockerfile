FROM golang:1.22.6-alpine
RUN apk add --no-cache build-base gcc
WORKDIR /forum
COPY . .
RUN apk add --no-cache sqlite
RUN go build -o forum main.go

CMD ["./forum"]
