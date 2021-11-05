#force rebuild docker build --no-cache -t tdewin/labtoken:latest .
FROM golang AS compiler
ENV CGO_ENABLED=0
ENV GOPROXY=direct
RUN go get -u github.com/tdewin/labtoken && go install github.com/tdewin/labtoken@latest && chmod 755 /go/bin/labtoken

FROM alpine
LABEL maintainer="@tdewin"
LABEL src="github.com/tdewin/labtoken"
WORKDIR /usr/sbin/
COPY --from=compiler /go/bin/labtoken /usr/sbin/labtoken
EXPOSE 8080
ENTRYPOINT /usr/sbin/labtoken