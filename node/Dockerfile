# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

FROM goreleaser/goreleaser-cross:v1.18.2-v1.9.0 AS build

RUN mkdir /build
WORKDIR /build

COPY ./go.mod /build
COPY ./go.sum /build
COPY ./cmd /build/cmd
COPY ./pkg /build/pkg

RUN go build -mod=readonly -trimpath -o /build/dist/bin/atomix-multi-raft-node ./cmd/atomix-multi-raft-node

FROM alpine:3.15

RUN apk add libc6-compat

RUN addgroup -S atomix && adduser -S -G atomix atomix

USER atomix

COPY --from=build /build/dist/bin/atomix-multi-raft-node /usr/local/bin/atomix-multi-raft-node

ENTRYPOINT ["atomix-multi-raft-node"]
