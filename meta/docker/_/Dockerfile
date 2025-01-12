ARG GO_VERSION=0.0.0
ARG PYTHON_VERSION=0.0.0


FROM golang:${GO_VERSION} AS go-base
WORKDIR /w
# Copy .mod and .sum files first to download dependencies and cache de layer
# before including any source code
COPY src/go.mod ./src/go.mod
COPY src/go.sum ./src/go.sum
RUN cd src && go mod download
# Copy the rest
COPY src ./src


FROM python:${PYTHON_VERSION} AS python-base
WORKDIR /w
# Poetry
ARG POETRY_VERSION='0.0.0'
ENV POETRY_VERSION="${POETRY_VERSION}"
ENV POETRY_INSTALLER_COMMIT="d62875fc05fb20062175cd14d19a96dbefa48640"
ENV POETRY_HOME=/poetry
ENV PATH="/poetry/bin:${PATH}"
RUN curl -fsSL "https://raw.githubusercontent.com/python-poetry/install.python-poetry.org/${POETRY_INSTALLER_COMMIT}/install-poetry.py" | python -
# Python deps
COPY ./meta/python/_/pyproject.toml ./meta/python/_/pyproject.toml
COPY ./meta/python/_/poetry.lock ./meta/python/_/poetry.lock
COPY ./meta/python/ensure-venv.sh ./meta/python/ensure-venv.sh
COPY ./meta/python/poetry.sh ./meta/python/poetry.sh
RUN ./meta/python/poetry.sh install


FROM go-base AS unit-test
COPY meta/test-unit.sh ./meta/test-unit.sh
RUN ./meta/test-unit.sh


FROM go-base AS e2e-build
COPY meta/build.sh ./meta/build.sh
RUN ./meta/build.sh


FROM python-base AS e2e-test
COPY ./meta/test-e2e.sh ./meta/test-e2e.sh
COPY ./tests ./tests
COPY ./playground ./playground
COPY ./docs ./docs
COPY --from=e2e-build /w/out/ebro ./out/ebro
COPY .git ./.git
RUN ./meta/test-e2e.sh


FROM go-base AS dist
# Next line is useless. It's to force the e2e-test
# stage to run
COPY --from=e2e-test /w/out/ebro /tested-ebro
COPY meta/build-dist.sh ./meta/build-dist.sh
COPY scripts ./scripts
COPY .git ./.git
RUN ./meta/build-dist.sh
RUN mv ./out /out


FROM python-base AS website
COPY ./meta/python/ensure-venv.sh ./meta/python/ensure-venv.sh
COPY ./meta/website/_/ ./meta/website/_/
COPY ./meta/website/build.sh ./meta/website/build.sh
COPY ./docs ./docs
COPY .git ./.git
RUN ./meta/website/build.sh
RUN mv ./out /out


FROM scratch AS all
COPY --from=dist /out /out
COPY --from=website /out /out
CMD ["false"]
