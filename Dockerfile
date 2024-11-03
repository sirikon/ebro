FROM golang:1.23.2
WORKDIR /wd
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY meta ./meta
COPY .git ./.git
RUN ./meta/release.sh
