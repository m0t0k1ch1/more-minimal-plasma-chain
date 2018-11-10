FROM node:10.13.0-stretch

RUN npm install -g ganache-cli@6.1.8

RUN mkdir -p /var/lib/ganache

ENTRYPOINT ["ganache-cli"]
CMD ["--host", "0.0.0.0", "--port", "8545", "--networkId", "5777", "--db", "/var/lib/ganache", "--deterministic"]
