FROM golang:1.25.1 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o /app/user .

FROM scratch

ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/user /app/user
COPY ./etc /app/etc

CMD ["./user", "-f", "etc/user.yaml"]
