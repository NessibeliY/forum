FROM golang:1.22.6-alpine
RUN apk add --no-cache build-base sqlite
WORKDIR /forum 
COPY . .
RUN go build -o forum .
CMD ["./forum"]