FROM golang:alpine as builder
WORKDIR /src
COPY . .
ENV ZONEINFO /opt/zoneinfo.zip
RUN go build -v -o ./bin/gba ./cli/*.go

FROM alpine
COPY --from=builder /src/bin/gba /gba
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT [ "/gba" ]
