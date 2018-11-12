FROM m0t0k1ch1/mmpc-child-base

WORKDIR /root

RUN apk --update add bash

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./
RUN chmod +x /root/wait-for-it.sh

RUN mkdir -p /var/lib/mmpc /usr/local/etc/mmpc

COPY ./app/config.json /usr/local/etc/mmpc/
COPY ./cli/config.json ./

ENTRYPOINT ["./wait-for-it.sh", "root:8545", "--", "more-minimal-plasma-chain"]
CMD ["--conf", "/usr/local/etc/mmpc/config.json"]
