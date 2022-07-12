# How to Containerize your go application?
#
# https://www.docker.com/blog/containerize-your-go-developer-environment-part-1/
#
# Official Dockerfile reference:
#
# https://docs.docker.com/engine/reference/builder/
#

# Note: Remeber to execute "make producer" before building this image.

FROM golang:1.18-alpine
WORKDIR /apps
COPY . .
CMD producer
