#FROM golang:1.11.5-stretch
FROM iotex/iotex-core:v0.5.1

# Install project
WORKDIR $GOPATH/src/github.com/Infinity-Stones/iotex_payout
COPY . .
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
	dep ensure && \
	go build && \
    chmod 755 iotex_payout && \
    cp $GOPATH/src/github.com/Infinity-Stones/iotex_payout/iotex_payout /usr/local/bin/iotex_payout && \
    ioctl config set endpoint api.iotex.one:80 --insecure
