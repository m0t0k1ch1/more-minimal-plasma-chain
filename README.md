# more-minimal-plasma-chain

[![GoDoc](https://godoc.org/github.com/m0t0k1ch1/more-minimal-plasma-chain?status.svg)](https://godoc.org/github.com/m0t0k1ch1/more-minimal-plasma-chain)

a Plasma chain for https://github.com/kfichter/more-minimal-plasma

## Quickstart

Please install [Docker Compose](https://docs.docker.com/compose/install) in advance.

``` sh
$ git clone git@github.com:m0t0k1ch1/more-minimal-plasma-chain.git
$ cd more-minimal-plasma-chain/_docker
$ docker-compose build
$ docker-compose up -d ganache
$ docker-compose up -d childchain
$ docker-compose exec childchain plasma deploy --privkey 0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d
```
