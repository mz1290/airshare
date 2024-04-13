build: Dockerfile
	docker build -t airshare .

run:
	docker run \
	-it --rm \
	-p 8080:8080 \
	--name airshare airshare