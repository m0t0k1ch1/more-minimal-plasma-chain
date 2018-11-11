FROM golang:1.11.2-alpine

# Install tools
RUN apk --update add gcc git musl-dev

# Install dep
RUN go get -u -v github.com/golang/dep/cmd/dep

# Download
RUN mkdir -p $GOPATH/src/github.com/m0t0k1ch1
WORKDIR $GOPATH/src/github.com/m0t0k1ch1/
RUN git clone https://github.com/m0t0k1ch1/more-minimal-plasma-chain.git

# Download dependencies
WORKDIR $GOPATH/src/github.com/m0t0k1ch1/more-minimal-plasma-chain/
RUN dep ensure -v

ENTRYPOINT ["go", "install"]
CMD ["-v", "./..."]
