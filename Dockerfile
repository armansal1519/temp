# Start from golang:1.12-alpine base image
#FROM golang:1.15-alpine AS build
#ENV http_proxy  http://fodev.org:8118
#
#WORKDIR /app
#
#COPY . ./
##RUN mkdir /app/images
#
## Install dependencies
#RUN go mod download && \
#  # Build the app
#  GOOS=linux GOARCH=amd64 go build -o main && \
#  # Make the final output executable
#  chmod +x ./main
#
#FROM alpine:latest
#
## Install os packages
#RUN apk --no-cache add bash
#
#WORKDIR /app
#
#COPY --from=build /app/main .
#
#CMD ["./main"]
#
#EXPOSE 3000

FROM golang:1.15-alpine AS build

WORKDIR /app

COPY . ./

# Install dependencies
RUN go mod download && \
  # Build the app
  GOOS=linux GOARCH=amd64 go build -o main && \
  # Make the final output executable
  chmod +x ./main

FROM alpine:latest

# Install os packages
RUN apk --no-cache add bash

WORKDIR /app

COPY --from=build /app/main .

CMD ["./main"]

EXPOSE 3000