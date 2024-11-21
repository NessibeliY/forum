build:
	docker build -t golang_image:1 .
run:
	docker run --rm -p 8080:8080 --name forum_container golang_image:1
