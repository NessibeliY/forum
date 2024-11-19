build:
	docker build -t custom_golang:1 .
run:
	docker run --rm -p 8080:8080 custom_golang:1