FROM golang:1.25 AS builder

WORKDIR /build

RUN go install go.k6.io/xk6/cmd/xk6@latest

RUN git clone https://github.com/NoOneBoss/xk6-noonechaos.git

RUN xk6 build \
    --with github.com/NoOneBoss/xk6-noonechaos=./xk6-noonechaos

RUN mv k6 /k6


FROM alpine:3.19

COPY --from=builder /k6 /k6

CMD ["cp", "/k6", "output"]