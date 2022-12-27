FROM shammishailaj/gobuilder:0.0.6 AS builder

ARG CGO
ARG OS
ARG ARCH

RUN mkdir /go/src/scli

WORKDIR /go/src/scli

COPY . ./

RUN go env

#RUN echo "CGO = "$CGO
#RUN echo "OS = "$OS
#RUN echo "ARCH = "$ARCH
RUN CGO_ENABLED=$CGO GOOS=$OS GOARCH=$ARCH govvv build -mod vendor -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/scli ./cmd/scli/scli.go
RUN cp ./bin/scli /bin/

# See: https://stackoverflow.com/a/52979541/6670698
FROM scratch AS production
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/scli/configs /scli/configs/
COPY --from=builder /go/src/scli/web /scli/web/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/src/scli/bin/scli /scli/scli
WORKDIR /scli
ENV TZ Local
ENTRYPOINT ["/scli/scli","serve"]