FROM harbor.ymt.io/inf/golang:1.16.2 as builder
ENV GOPROXY=http://goproxy.ymt360.com
ENV GOSUMDB=off

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

ADD . .
RUN CGO_ENABLED=0 go build -o ymt-ladon -a -ldflags '-s' main.go &&\
    chmod +x ymt-ladon

FROM scratch
ADD http://source.ymt.io/optools/cacert.pem /etc/ssl/certs/cacert.pem
COPY --from=builder /app/passport /passport
CMD ["/passport"]
