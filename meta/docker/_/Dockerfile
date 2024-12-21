ARG GO_VERSION=0.0.0

FROM golang:${GO_VERSION}
WORKDIR /w

# Copy .mod and .sum files first to donwload dependencies and cache de layer
# before including any source code
COPY src/go.mod ./src/go.mod
COPY src/go.sum ./src/go.sum
RUN cd src && go mod download

# Copy the rest
COPY src ./src
COPY meta/build-dist.sh ./meta/build-dist.sh
COPY .git ./.git

# Build
RUN ./meta/build-dist.sh
RUN mv ./out /out
