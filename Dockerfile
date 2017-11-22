FROM golang:1.8.3-alpine

WORKDIR /go/src/github.com/kristiehoward/go2flow
COPY . /go/src/github.com/kristiehoward/go2flow

RUN CGO_ENABLED=0 go install -v -ldflags="-s" && \
	GOOS=linux GOARCH=s390x  go install -v -ldflags="-s" && \
	GOOS=linux GOARCH=ppc64le  go install -v -ldflags="-s" && \
	GOOS=darwin GOARCH=amd64 go install -v -ldflags="-s" && \
	GOOS=windows GOARCH=amd64 go install -v -ldflags="-s"


FROM alpine:3.6

RUN apk --no-cache add ca-certificates

COPY --from=0 /go/bin/go2flow /go2flow-Linux-x86_64
COPY --from=0 /go/bin/linux_s390x/go2flow /go2flow-Linux-s390x
COPY --from=0 /go/bin/linux_ppc64le/go2flow /go2flow-Linux-ppc64le
COPY --from=0 /go/bin/darwin_amd64/go2flow /go2flow-Darwin-x86_64
COPY --from=0 /go/bin/windows_amd64/go2flow.exe /go2flow-Windows-x86_64

RUN ln -s /go2flow-Linux-x86_64 /usr/local/bin/go2flow

ENTRYPOINT ["go2flow"]