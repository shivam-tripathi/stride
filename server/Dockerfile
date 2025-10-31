ARG env

FROM public.ecr.aws/docker/library/golang:1.23-alpine3.20 AS build

RUN mkdir -p /usr/app
WORKDIR /usr/app

RUN apk add build-base
RUN go install github.com/google/wire/cmd/wire@latest

COPY . .


# Install make and build using Makefile
RUN apk add --no-cache make
RUN CGO_ENABLED=0 GOOS=linux make all
RUN cp server /go/bin/server

FROM public.ecr.aws/docker/library/golang:1.23-alpine3.20 AS deployment
ARG env

RUN mkdir -p /usr/app
WORKDIR /usr/app

COPY  --from=build /usr/app/server .

RUN echo "Building docker image for env: $env"

ENV AWS_REGION=us-east-1
ENV DEBUG=q:app:*
ENV ENV=${env}
ENV PORT=8080

EXPOSE 8080

CMD ./server