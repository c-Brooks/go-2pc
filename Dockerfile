FROM golang:1.9 

LABEL maintainer Corey Brooks <corey.brooks@food.ee>

# Mount the app as a volume here
WORKDIR /go/src/github.com/c-Brooks/go-2pc

EXPOSE 8080
EXPOSE 50051