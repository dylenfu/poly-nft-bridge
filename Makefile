# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build -tags $(env)
GOTEST=$(GOCMD) test

# Set and confirm environment `BRIDGE`, it should be one of devnet/testnet/mainnet
env=$(BRIDGE)
BaseDir=build/$(env)

# command: `make bridge_http env=testnet` or `make bridge_http env=mainnet`
bridge_http:
	mkdir -p $(BaseDir)/bridge_http
	@cp -r conf/app.conf $(BaseDir)/bridge_http/
	@$(GOBUILD) -o $(BaseDir)/bridge_http/bridge_http cmd/bridge_http/main.go

# command `make bridge_tools env=testnet/mainnet/local`
bridge_tools:
	mkdir -p $(BaseDir)/bridge_tools
	# todo: copy config files to build/.../bridge_tools
	@$(GOBUILD) -o $(BaseDir)/bridge_tools/bridge_tools cmd/bridge_tools/*.go

bridge_server:
	mkdir -p $(BaseDir)/bridge_server
	# todo: copy config files to build/.../bridge_server
	@$(GOBUILD) -o $(BaseDir)/bridge_server/bridge_server cmd/bridge_server/main.go

ethereum_listen:
	mkdir -p $(BaseDir)/ethereum_listen
	# todo: copy config files to build/.../ethereum_listen
	@$(GOBUILD) -o $(BaseDir)/ethereum_listen/ethereum_listen cmd/ethereum_listen/*.go

neo_listen:
	mkdir -p $(BaseDir)/neo_listen
	# todo: copy config files to build/.../neo_listen
	@$(GOBUILD) -o $(BaseDir)/neo_listen/neo_listen cmd/neo_listen/main.go

ontology_listen:
	mkdir -p $(BaseDir)/ontology_listen
	# todo: copy config files to build/.../ontology_listen
	@$(GOBUILD) -o $(BaseDir)/ontology_listen/ontology_listen cmd/ontology_listen/main.go

poly_listen:
	mkdir -p $(BaseDir)/poly_listen
	# todo: copy config files to build/.../poly_listen
	@$(GOBUILD) -o $(BaseDir)/poly_listen/poly_listen cmd/poly_listen/main.go

chainfee_listen:
	mkdir -p $(BaseDir)/chainfee_listen
	# todo: copy config files to build/.../chainfee_listen
	@$(GOBUILD) -o $(BaseDir)/chainfee_listen/chainfee_listen cmd/chainfee_listen/*.go

crosschain_effect:
	mkdir -p $(BaseDir)/crosschain_effect
	# todo: copy config files to build/.../crosschain_effect
	@$(GOBUILD) -o $(BaseDir)/crosschain_effect/crosschain_effect cmd/crosschain_effect/*.go

coinprice_listen:
	mkdir -p $(BaseDir)/coinprice_listen
	# todo: copy config files to build/.../coinprice_listen
	@$(GOBUILD) -o $(BaseDir)/coinprice_listen/coinprice_listen cmd/coinprice_listen/*.go

server:
	make bridge_http bridge_tools bridge_server

deploy-tool:
	@mkdir -p $(BaseDir)/deploy_tool/keystore
	@mkdir -p $(BaseDir)/deploy_tool/leveldb
	@cp cmd/deploy_tool/config_$(env).json $(BaseDir)/deploy_tool/config.json
	@cp -R keystore/$(env)/ $(BaseDir)/deploy_tool/keystore/
	@$(GOBUILD) -o $(BaseDir)/deploy_tool/deploy_tool cmd/deploy_tool/*.go

all:
	make bridge_http bridge_tools bridge_server ethereum_listen neo_listen ontology_listen poly_listen chainfee_listen crosschain_effect coinprice_listen

clean:
	rm -rf build/*
#
# compile-linux:
# 	GOOS=linux GOARCH=amd64 $(GOBUILD) -o build/$(ENV)/robot-linux cmd/main.go
#
# robot:
# 	@echo test case $(t)
# 	./build/$(ENV)/robot -config=build/$(ENV)/config.json -t=$(t)
#
# clean: