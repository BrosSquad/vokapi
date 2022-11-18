FROM golang:1.19 as build

ARG VERSION="1.0.0"

WORKDIR /app

COPY . .

RUN apt update \
    && apt upgrade -y \
    && apt install gcc g++ make -y

RUN mkdir -p build \
    && make build VERSION=${VERSION} ENV=production

FROM alpine:3 as production

WORKDIR /app

COPY --from=build /app/bin/vokapi /bin/vokapi

RUN apk add tini \
    && chmod +x /bin/vokapi \
    && mkdir -p /var/lib/vokapi

ENTRYPOINT [ "/bin/tini" ]
CMD [ "vokapi", "--db-path=/var/lib/vokapi/data", "server", "--host=0.0.0.0", "--port=1389" ]
