FROM golang:1.26 AS builder

WORKDIR /build

RUN apt-get update && \
    apt-get install -y make nodejs npm

RUN npm install -g @apidevtools/swagger-cli

COPY . .

RUN make install_tools

RUN make clean init

FROM registry.access.redhat.com/ubi9-minimal

WORKDIR /app

COPY --from=builder /build/generated/app/backend ./

EXPOSE 1323

CMD ["./backend"]