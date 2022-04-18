# How to Containerize your go application : https://www.docker.com/blog/tag/go-env-series/
# Official Dockerfile reference: https://www.docker.com/blog/containerize-your-go-developer-environment-part-1/

#
# Build
#
FROM golang:1.18-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o

#
# Override config file with specified env
#
ARG env=stg
