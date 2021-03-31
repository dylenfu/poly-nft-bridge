#!/bin/bash

# eth
./deploy_tool --chain=2 deployERC20
./deploy_tool --chain=2 mintERC20 --to=0x31c0dd87B33Dcd66f9a255Cf4CF39287F8AE593C --amount=100000000000000000000000000

./deploy_tool --chain=2 deployECCD
./deploy_tool --chain=2 deployECCM
./deploy_tool --chain=2 deployCCMP

./deploy_tool --chain=2 transferECCDOwnership
./deploy_tool --chain=2 transferECCMOwnership

./deploy_tool --chain=2 registerSideChain
./deploy_tool --chain=2 approveSideChain

# bsc
./deploy_tool --chain=6 deployERC20
./deploy_tool --chain=6 mintERC20 --to=0x896fB9Dd4Bddd1C4ea2cab3df66C632AD736a9D1 --amount=100000000000000000000000000

./deploy_tool --chain=6 deployECCD
./deploy_tool --chain=6 deployECCM
./deploy_tool --chain=6 deployCCMP

./deploy_tool --chain=6 transferECCDOwnership
./deploy_tool --chain=6 transferECCMOwnership

./deploy_tool --chain=6 registerSideChain
./deploy_tool --chain=6 approveSideChain
