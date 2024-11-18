FROM golang:1.23.3
WORKDIR /wd

# Copy .mod and .sum files first to donwload dependencies and cache de layer
# before including any source code
COPY src/go.mod ./src/go.mod
COPY src/go.sum ./src/go.sum
RUN cd src && go mod download

# Copy the rest and run release
COPY src ./src
COPY meta/release.sh ./meta/release.sh
COPY .git ./.git
RUN ./meta/release.sh

RUN mv ./dist /dist
