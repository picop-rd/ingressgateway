FROM golang:1.21-bullseye AS builder

WORKDIR /go/src/github.com/picop-rd/ingressgateway/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /ingressgateway ./cmd/ingressgateway/main.go


FROM scratch

COPY --from=builder /ingressgateway /bin/ingressgateway
ENTRYPOINT [ "/bin/ingressgateway" ]

