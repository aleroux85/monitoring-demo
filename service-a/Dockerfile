FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN apk add --update git gcc linux-headers
RUN go get github.com/prometheus/client_golang/prometheus
RUN CGO_ENABLED=0 go build main.go

FROM scratch
COPY --from=builder /app/main /main
# COPY ./service-a/main /main
ENTRYPOINT ["/main"]