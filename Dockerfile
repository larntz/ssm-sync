FROM golang:1.18.5-alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ssm-sync .
FROM alpine
COPY --from=builder /build/ssm-sync /app/
WORKDIR /app
CMD ["./ssm-sync"]
