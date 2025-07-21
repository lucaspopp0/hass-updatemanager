ARG BUILD_FROM
ARG GO_BASE

FROM golang:1.24.5 AS go

FROM $BUILD_FROM AS homeassistant

COPY --from=go /usr/local/go /usr/local/go
ENV PATH "/usr/local/go/bin:$PATH"

RUN mkdir -p /update-manager/server

COPY ./server /update-manager/server

WORKDIR /update-manager/server
RUN go get ./...
RUN CG_ENABLED=0 go build -o server .

COPY run.sh /
RUN chmod a+x /run.sh

# Persistent data directory
WORKDIR /data

# Expose port
EXPOSE 8000
CMD [ "/run.sh" ]

LABEL org.opencontainers.image.source=https://github.com/lucaspopp0/hass-update-manager
LABEL org.opencontainers.image.description="Add-on for home assistant"
LABEL org.opencontainers.image.licenses=MIT
