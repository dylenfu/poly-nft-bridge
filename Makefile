# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build -tags $(env)
GOTEST=$(GOCMD) test

# Set and confirm environment `BRIDGE`, it should be one of devnet/testnet/mainnet
env=$(BRIDGE)
BaseDir=build/$(env)

prepare:
	@mkdir -p $(BaseDir)/bridge_http
	@mkdir -p $(BaseDir)/ethereum_listen/
	@mkdir -p $(BaseDir)/poly_listen
	@mkdir -p $(BaseDir)/deploy_tool/keystore
	@mkdir -p $(BaseDir)/deploy_tool/leveldb
	@cp -r cmd/bridge_http/app_$(env).conf $(BaseDir)/bridge_http/app.conf
	@cp -r conf/config_$(env).json $(BaseDir)/ethereum_listen/config.json
	@cp -r conf/config_$(env).json $(BaseDir)/poly_listen/config.json
	@cp cmd/deploy_tool/config_$(env).json $(BaseDir)/deploy_tool/config.json
	@cp -R keystore/$(env)/ $(BaseDir)/deploy_tool/keystore/

bridge_http:
	@$(GOBUILD) -o $(BaseDir)/bridge_http/bridge_http cmd/bridge_http/main.go

ethereum_listen:
	@$(GOBUILD) -o $(BaseDir)/ethereum_listen/ethereum_listen cmd/ethereum_listen/*.go

poly_listen:
	@$(GOBUILD) -o $(BaseDir)/poly_listen/poly_listen cmd/poly_listen/main.go

deploy_tool:
	@$(GOBUILD) -o $(BaseDir)/deploy_tool/deploy_tool cmd/deploy_tool/*.go

all:
	make bridge_http ethereum_listen poly_listen deploy_tool

#
# compile-linux:
# 	GOOS=linux GOARCH=amd64 $(GOBUILD) -o build/$(ENV)/robot-linux cmd/main.go
#
# robot:
# 	@echo test case $(t)
# 	./build/$(ENV)/robot -config=build/$(ENV)/config.json -t=$(t)
#
# clean: