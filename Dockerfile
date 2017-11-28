FROM golang:1.9 

LABEL maintainer Corey Brooks <corey.brooks@food.ee>

# Install delve debugger
RUN go get github.com/derekparker/delve/cmd/dlv

# Mount the app as a volume here
WORKDIR /go/src/github.com/c-Brooks/go-2pc

EXPOSE 8080
EXPOSE 50051