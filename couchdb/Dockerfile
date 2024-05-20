FROM bitnami/couchdb:3.3.3 AS couchdb

COPY --from=bitnami/golang:1.22.3 /opt/bitnami/go /opt/bitnami/go

ENV PATH="/opt/bitnami/go/bin:${PATH}"
