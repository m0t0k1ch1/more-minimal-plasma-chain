#!/bin/sh

set -eux

# deploy root chain contract
docker-compose exec child plasma deploy --privkey 0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d

# deposit
docker-compose exec child plasma deposit make --amount 1000000000000000000 --privkey 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1

# post tx
docker-compose exec child plasma tx post --pos 1000000000 --address 0x22d491bde2303f2f43325b2108d26f1eaba1e32b --amount 500000000000000000 --privkey 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1

# create block
docker-compose exec child plasma block fix

# confirm tx
docker-compose exec child plasma txin confirm --pos 2000000000 --privkey 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1

# start exit (invalid)
docker-compose exec child plasma exit start --pos 1000000000 --privkey 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1

# challenge exit
docker-compose exec child plasma exit challenge --pos 1000000000 --vspos 2000000000 --privkey 0x6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c

# start exit (valid)
docker-compose exec child plasma exit start --pos 2000000000 --privkey 0x6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c

# increase time (2 weeks)
curl -s -X POST http://127.0.0.1:8545 --data '{"jsonrpc": "2.0", "method": "evm_increaseTime", "params": [1209600], "id": 0}'
curl -s -X POST http://127.0.0.1:8545 --data '{"jsonrpc": "2.0", "method": "evm_mine", "params": [], "id": 0}'

# process exit
docker-compose exec child plasma exit process --privkey 0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d
