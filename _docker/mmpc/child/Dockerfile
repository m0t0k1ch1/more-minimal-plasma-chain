FROM golang:1.11.2-alpine

WORKDIR /root

RUN apk --update add bash

# Install wait-for-it.sh
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./
RUN chmod +x /root/wait-for-it.sh

# Install
RUN mkdir -p /usr/local/etc/mmpc /var/lib/mmpc
COPY ./app/more-minimal-plasma-chain /usr/local/bin/
COPY ./app/config.json /usr/local/etc/mmpc/

# Install CLI
COPY ./cli/plasma /usr/local/bin/
COPY ./cli/config.json ./

ENTRYPOINT ["./wait-for-it.sh", "root:8545", "--", "more-minimal-plasma-chain"]
CMD ["--conf", "/usr/local/etc/mmpc/config.json"]
