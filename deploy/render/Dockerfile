# Wraps the application image with Litestream so SQLite survives Render's ephemeral disk.
# Build with: docker build --build-arg APP_IMAGE=<app image> -f deploy/render/Dockerfile .
ARG APP_IMAGE
FROM litestream/litestream:0.3.13 AS litestream

FROM ${APP_IMAGE}
COPY --from=litestream /usr/local/bin/litestream /usr/local/bin/litestream
COPY deploy/render/litestream.yml /etc/litestream.yml
COPY deploy/render/entrypoint.sh /render-entrypoint.sh
USER root
RUN chmod +x /render-entrypoint.sh
USER 1000:1000
ENTRYPOINT ["/render-entrypoint.sh"]
