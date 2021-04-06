# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build -tags $(env)
GOTEST=$(GOCMD) test

# Set and confirm environment `BRIDGE`, it should be one of devnet/testnet/mainnet
env=$(BRIDGE)
BaseDir=build/$(env)

prepare:
	@mkdir -p $(BaseDir)/bridge_http
	@mkdir -p $(BaseDir)/eth_listen/
	@mkdir -p $(BaseDir)/bsc_listen/
	@mkdir -p $(BaseDir)/heco_listen/
	@mkdir -p $(BaseDir)/poly_listen/
	@mkdir -p $(BaseDir)/deploy_tool/keystore
	@mkdir -p $(BaseDir)/deploy_tool/leveldb
	@cp -r cmd/bridge_http/app_$(env).conf $(BaseDir)/bridge_http/app.conf
	@cp -r conf/config_$(env).json $(BaseDir)/eth_listen/config.json
	@cp -r conf/config_$(env).json $(BaseDir)/bsc_listen/config.json
	@cp -r conf/config_$(env).json $(BaseDir)/heco_listen/config.json
	@cp -r conf/config_$(env).json $(BaseDir)/poly_listen/config.json
	@cp -r cmd/deploy_tool/config_$(env).json $(BaseDir)/deploy_tool/config.json

bridge_http:
	@$(GOBUILD) -o $(BaseDir)/bridge_http/bridge_http cmd/bridge_http/main.go

eth_listen:
	@$(GOBUILD) -o $(BaseDir)/eth_listen/listener cmd/eth_listen/*.go
	@cp $(BaseDir)/eth_listen/listener $(BaseDir)/bsc_listen/listener
	@cp $(BaseDir)/eth_listen/listener $(BaseDir)/heco_listen/listener

poly_listen:
	@$(GOBUILD) -o $(BaseDir)/poly_listen/listener cmd/poly_listen/main.go

deploy_tool:
	@cp -R keystore/$(env)/ $(BaseDir)/deploy_tool/keystore/
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