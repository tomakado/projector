#!/bin/sh

docker build -f DockerfileTests -t projector-tests .. && \
	docker run projector-tests
