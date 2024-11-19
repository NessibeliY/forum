FROM golang:1.20.1-alpine
RUN apk add build-base 
WORKDIR /forum 
COPY . .
RUN go build -o forum .
CMD ["./forum"]