FROM golang:1.11.2-alpine

RUN apk --update add gcc git musl-dev

RUN go get -u -v github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/m0t0k1ch1
WORKDIR /go/src/github.com/m0t0k1ch1
RUN git clone https://github.com/m0t0k1ch1/more-minimal-plasma-chain.git
WORKDIR /go/src/github.com/m0t0k1ch1/more-minimal-plasma-chain
RUN dep ensure -v
RUN go install -v ./...

ENTRYPOINT ["more-minimal-plasma-chain"]
