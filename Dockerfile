FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o crm.shop.com ./cmd/server/

FROM scratch

COPY ./config /config

COPY --from=builder /app/crm.shop.com /

ENTRYPOINT [ "./crm.shop.com" , "./config/dev.yaml" ]

