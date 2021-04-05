# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build -tags $(env)
GOTEST=$(GOCMD) test

# Set and confirm environment `BRIDGE`, it should be one of devnet/testnet/mainnet
env=$(BRIDGE)
BaseDir=build/$(env)

bridge_http:
	mkdir -p $(BaseDir)/bridge_http
	@cp -i cmd/bridge_http/app_$(env).conf $(BaseDir)/bridge_http/app.conf
	@$(GOBUILD) -o $(BaseDir)/bridge_http/bridge_http cmd/bridge_http/main.go

ethereum_listen:
	mkdir -p $(BaseDir)/ethereum_listen/
	@cp -i conf/config_$(env).json $(BaseDir)/ethereum_listen/config.json
	@$(GOBUILD) -o $(BaseDir)/ethereum_listen/ethereum_listen cmd/ethereum_listen/*.go

poly_listen:
	mkdir -p $(BaseDir)/poly_listen
	@cp -i conf/config_$(env).json $(BaseDir)/poly_listen/config.json
	@$(GOBUILD) -o $(BaseDir)/poly_listen/poly_listen cmd/poly_listen/main.go

server:
	make bridge_http bridge_tools bridge_server

prepare-deploy-tool:
	@mkdir -p $(BaseDir)/deploy_tool/keystore
	@mkdir -p $(BaseDir)/deploy_tool/leveldb
	@cp cmd/deploy_tool/config_$(env).json $(BaseDir)/deploy_tool/config.json
	@cp -R keystore/$(env)/ $(BaseDir)/deploy_tool/keystore/

deploy-tool:
	@$(GOBUILD) -o $(BaseDir)/deploy_tool/deploy_tool cmd/deploy_tool/*.go

all:
	make bridge_http bridge_tools bridge_server ethereum_listen poly_listen
#
# compile-linux:
# 	GOOS=linux GOARCH=amd64 $(GOBUILD) -o build/$(ENV)/robot-linux cmd/main.go
#
# robot:
# 	@echo test case $(t)
# 	./build/$(ENV)/robot -config=build/$(ENV)/config.json -t=$(t)
#
# clean: