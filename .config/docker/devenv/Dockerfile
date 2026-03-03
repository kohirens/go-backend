ARG USER_NAME='app'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/go-backend'

FROM golang:1-alpine3.23 AS base

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

# Update OS
RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
 && apk --no-progress --purge --no-cache upgrade \
 && apk add \
        curl \
        openssl \
        tzdata \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/*

ENV SHELL=/bin/sh
ENV TZ=America/Detroit

# Add a non-root group and user.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=3 \
    CMD exit 0

EXPOSE 443/tcp

#---------- end base stage ------------

#---------- dev stage ------------

FROM base AS dev

# Install development apps
RUN apk add \
        bash \
        git \
        jq \
        openssh \
        sudo \
        zip \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/*

    # Add the user to sudo groups.
RUN adduser ${USER_NAME} wheel \
 && echo '%wheel ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/wheel

RUN curl -L -o gen-ssl-cert.sh https://raw.githubusercontent.com/b01/script-lib/refs/heads/main/generate-self-signed-cert.sh \
 && chmod 555 ./gen-ssl-cert.sh \
 && mv gen-ssl-cert.sh /usr/local/bin

USER ${USER_NAME}

RUN sudo gen-ssl-cert.sh --company="containers" \
    --sans="DNS:localhost,DNS:ip6-loopback,IP:127.0.0.1,IP:0:0:0:0:0:0:0:1" \
    --out-dir="${HOME}/pki" \
    --cert-ext="crt" \
    "dev.kohirens.local" \
    -v

# Make directories with current user permissions
RUN mkdir -p ~/bin ~/src ~/pkg/mod/cache

# Set the working directory in the container
WORKDIR "/home/${USER_NAME}/src/${REPO}"

COPY --chmod=644 --chown=${USER_NAME}:${USER_NAME} .config/docker/.gitconfig /home/${USER_NAME}
COPY --chmod=755 --chown=${USER_NAME}:${USER_NAME} .config/docker/devenv/start.sh /home/${USER_NAME}/bin
COPY --chmod=755 --chown=${USER_NAME}:${USER_NAME} . .
COPY --from=stdlib --chown=${USER_NAME}:${USER_NAME} / /home/${USER_NAME}/src/github.com/kohirens/stdlib/

ENTRYPOINT [ "start.sh" ]

CMD [ "" ]

ENV PATH="${PATH}:/home/${USER_NAME}/bin"
ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0

#---------- end dev stage ------------
