FROM golang:1.12.7 as builder

LABEL description="Webtoon service image builder"
LABEL maintainer="Youngjoon Lee taxihighway@gmail.com"

ARG BASEDIR=/opt/webtoon-service
COPY . $BASEDIR
WORKDIR $BASEDIR

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -a -installsuffix cgo \
    -o webtoon \
    cmd/webtoon.go


FROM alpine:3.10

ARG BASEDIR=/opt/webtoon-service
WORKDIR $BASEDIR
COPY --from=builder $BASEDIR/webtoon .

EXPOSE 8080
CMD ["./webtoon"]