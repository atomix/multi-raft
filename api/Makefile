# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

build:
	go build ./...

generate: go docs

go:
	@find . -name '*.pb.go' -delete
	docker run -it \
		-v `pwd`:/build \
		atomix/codegen:go-latest \
	    --proto-path . --go-path . --import-path github.com/atomix/multi-raft-storage/api
	docker run -it \
		-v `pwd`:/build \
		atomix/codegen:latest \
		protoc -I=.:/go/src/github.com/gogo/protobuf \
			--go_out=Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,import_path=github.com/atomix/multi-raft-storage/api/atomix/multiraft/v1:. \
			atomix/multiraft/v1/descriptor.proto

docs:
	@find . -name '*.md' -delete
	docker run -it \
		-v `pwd`:/build \
		atomix/codegen:docs-latest \
		--proto-path . --docs-path . --docs-format markdown
