FROM golang:1.12.6-alpine3.10 AS golang

RUN apk add --no-cache git
RUN go get github.com/golang/dep && go install github.com/golang/dep/cmd/dep
RUN mkdir -p /builds/go/src/github.com/betorvs/sensu-go-slack-bot/
ENV GOPATH /builds/go
COPY . /builds/go/src/github.com/betorvs/sensu-go-slack-bot/
ENV CGO_ENABLED 0
RUN cd /builds/go/src/github.com/betorvs/sensu-go-slack-bot/ && dep ensure -v && go build

FROM alpine:3.10
WORKDIR /
VOLUME /tmp
RUN apk add --no-cache ca-certificates
COPY --from=golang /builds/go/src/github.com/betorvs/sensu-go-slack-bot/sensu-go-slack-bot /
RUN update-ca-certificates

EXPOSE 9090
RUN chmod +x /sensu-go-slack-bot
CMD ["/sensu-go-slack-bot"]