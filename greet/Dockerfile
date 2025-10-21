FROM golang:1.25.1 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o /app/greet greet.go


FROM scratch

ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/greet /app/greet
COPY ./etc /app/etc

EXPOSE 8888

CMD ["./greet", "-f", "etc/greet-api.yaml"]
