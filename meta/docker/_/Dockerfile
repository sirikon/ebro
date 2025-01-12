FROM debian:bookworm-slim
WORKDIR /w
RUN apt-get update && apt-get install -y curl

# Mise
ENV MISE_VERSION=v2025.1.5
ENV MISE_INSTALL_PATH=/usr/local/bin/mise
RUN curl https://mise.run | sh
RUN echo 'eval "$(mise activate bash)"' >> ~/.bashrc
SHELL [ "bash", "-lc" ]

COPY .tool-versions .
RUN mise install python
RUN mise install
RUN rm .tool-versions
ENTRYPOINT [ "bash", "-lc" ]
