FROM golang:1.24.0-alpine AS base
WORKDIR /work
RUN apk update && apk add git
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

FROM base AS dev
VOLUME [ "./cmd" ]
VOLUME [ "./internal" ]
CMD ["go", "run", "./cmd/server/main.go"]

FROM base AS build
COPY cmd/ cmd/
COPY internal/ internal/
RUN go build -o /tmp/server ./cmd/server/main.go

#FROM scratch AS deploy
#COPY --from=base /tmp/server /usr/local/bin/server
#CMD ["server"]