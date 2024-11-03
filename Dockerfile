FROM golang:1.23.2
WORKDIR /wd
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY meta/release.sh ./meta/release.sh
COPY .git ./.git
RUN ./meta/release.sh
RUN mv ./dist /dist
