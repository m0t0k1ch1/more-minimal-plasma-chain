FROM golang:1.11.2-alpine

# Install tools
RUN apk --update add gcc git jq musl-dev

# Install dep
RUN go get -u -v github.com/golang/dep/cmd/dep

# Install
RUN mkdir -p $GOPATH/src/github.com/m0t0k1ch1
WORKDIR $GOPATH/src/github.com/m0t0k1ch1/
RUN git clone https://github.com/m0t0k1ch1/more-minimal-plasma-chain.git
WORKDIR $GOPATH/src/github.com/m0t0k1ch1/more-minimal-plasma-chain/
RUN dep ensure -v
RUN go install -v ./...

# Copy configs
RUN mkdir -p /usr/local/etc/mmpc
COPY ./app/config.json /usr/local/etc/mmpc/
COPY ./cli/config.json $HOME/

# Make data directory
RUN mkdir -p /var/lib/mmpc

WORKDIR $HOME/

ENTRYPOINT ["more-minimal-plasma-chain"]
CMD ["--conf", "/usr/local/etc/mmpc/config.json"]
