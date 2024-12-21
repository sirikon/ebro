ARG GO_VERSION=0.0.0
ARG PYTHON_VERSION=0.0.0



FROM golang:${GO_VERSION} AS dist
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



FROM python:${PYTHON_VERSION} AS website

# Poetry
ARG POETRY_VERSION='0.0.0'
ENV POETRY_VERSION="${POETRY_VERSION}"
ENV POETRY_INSTALLER_COMMIT="d62875fc05fb20062175cd14d19a96dbefa48640"
ENV POETRY_HOME=/poetry
ENV PATH="/poetry/bin:${PATH}"
RUN curl -fsSL "https://raw.githubusercontent.com/python-poetry/install.python-poetry.org/${POETRY_INSTALLER_COMMIT}/install-poetry.py" | python -

WORKDIR /w
COPY ./meta/python/_/pyproject.toml ./meta/python/_/pyproject.toml
COPY ./meta/python/_/poetry.lock ./meta/python/_/poetry.lock
COPY ./meta/python/poetry.sh ./meta/python/poetry.sh
RUN ./meta/python/poetry.sh install
COPY ./meta/python/ensure-venv.sh ./meta/python/ensure-venv.sh
COPY ./meta/website/_/ ./meta/website/_/
COPY ./meta/website/build.sh ./meta/website/build.sh
COPY ./docs ./docs
RUN ./meta/website/build.sh



FROM scratch
COPY --from=dist /w/out /out
COPY --from=website /w/out /out
CMD ["false"]
