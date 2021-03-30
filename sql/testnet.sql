use polyswap;

-- desc src_transactions;
-- desc src_transfers;
-- desc tokens;

-- set @txhash='00631d440e7a140c87434a9f441aea56d5b56a5b81892b0e6caeb58333b487cf';
-- set @user='a107c23029c31da1b5ab19eab8228a2a44024c7d';
-- set @property=1;
-- set @chainid=3;

-- explain select 
-- srctx.chain_id, srctx.state as tx_state, srctx.hash as tx_hash,
-- srctran.from as tx_from, srctran.to as tx_to, srctran.amount,
-- tks.token_basic_name, tks.property
-- from src_transactions as srctx
-- inner join src_transfers as srctran on srctx.hash=srctran.tx_hash
-- inner join tokens as tks on tks.hash=srctran.asset and tks.chain_id=@chainid
-- where srctran.from=@user
-- and tks.property=@property;

select * from token_maps where src_token_hash='49a98BBb058b666886F661dDeE6B431C5Df9d9Fd';






