# more-minimal-plasma-chain

[![GoDoc](https://godoc.org/github.com/m0t0k1ch1/more-minimal-plasma-chain?status.svg)](https://godoc.org/github.com/m0t0k1ch1/more-minimal-plasma-chain)

a Plasma chain for https://github.com/kfichter/more-minimal-plasma

## Quickstart

Please install [Docker Compose](https://docs.docker.com/compose/install) in advance.

__NOTICE: The private keys used in the following process are generated by ganache-cli with `--deterministic` option.__

- Operator
  - address: 0x90f8bf6a479f320ead074411a4b0e7944ea8c9c1
  - privkey: 0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d
- Alice
  - address: 0xffcf8fdee72ac11b5c542428b35eef5769c409f0
  - privkey: 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1
- Bob
  - address: 0x22d491bde2303f2f43325b2108d26f1eaba1e32b
  - privkey: 0x6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c

### Build & Run

``` sh
$ cd _docker/mmpc
$ docker-compose up --build -d
```

### Deploy root chain contract

Operator deploys root chain contract.

``` sh
$ docker-compose exec child plasma deploy --privkey 0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d
```

### Deposit

Alice deposits 1 ETH. Block#1 is automatically created.

``` sh
$ docker-compose exec child plasma deposit make --amount 1000000000000000000 --privkey 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1
```

### Transfer

Alice sends 0.5 ETH to Bob.

``` sh
$ docker-compose exec child plasma tx post --pos 1000000000 --address 0x22d491bde2303f2f43325b2108d26f1eaba1e32b --amount 500000000000000000 --privkey 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1
```

Operator creates block#2.

``` sh
$ docker-compose exec child plasma block fix
```

Alice confirms the tx.

``` sh
$ docker-compose exec child plasma txin confirm --pos 2000000000 --privkey 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1
```

### Start exit (invalid)

Alice starts invalid 1 ETH exit.

``` sh
$ docker-compose exec child plasma exit start --pos 1000000000 --privkey 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1
```

### Challenge exit

Bob challenges the invalid exit.

``` sh
$ docker-compose exec child plasma exit challenge --pos 1000000000 --vspos 2000000000 --privkey 0x6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c
```

### Start exit (valid)

Bob starts valid 0.5 ETH exit.

``` sh
$ docker-compose exec child plasma exit start --pos 2000000000 --privkey 0x6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c
```

### Process exits

2 weeks have passed.

``` sh
$ curl -X POST http://127.0.0.1:8545 --data '{"jsonrpc": "2.0", "method": "evm_increaseTime", "params": [1209600], "id": 0}'
$ curl -X POST http://127.0.0.1:8545 --data '{"jsonrpc": "2.0", "method": "evm_mine", "params": [], "id": 0}'
```

Operator processes exits.

``` sh
$ docker-compose exec child plasma exit process --privkey 0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d
```
