use polyswap;

select * from wrapper_transactions; -- '8c0a5fc96d50b914f5ea9c3ad1dfbbf6b173bf539aa8990548ea56d003c40fb2';
select * from src_transactions where hash='8c0a5fc96d50b914f5ea9c3ad1dfbbf6b173bf539aa8990548ea56d003c40fb2';
select * from src_transfers where tx_hash='8c0a5fc96d50b914f5ea9c3ad1dfbbf6b173bf539aa8990548ea56d003c40fb2';
select * from poly_transactions where src_hash='8c0a5fc96d50b914f5ea9c3ad1dfbbf6b173bf539aa8990548ea56d003c40fb2';
select * from dst_transactions where poly_hash='98c9d48686b0760f1f8d81283e2acdc2457b96a6d21c3cf5db2dc9b76bb92983'; 
select * from dst_transfers where tx_hash='6e81ec2e6dddcc5739b84bf9ff1b5ec649420175e57cf8a794bd65788516ef5e';

-- SELECT * FROM `wrapper_transactions` 
-- WHERE standard = 1 and status = 2 
-- and (user in ('5fb03eb21303d39967a1a119b32dd744a0fa8986') or dst_user in ('5fb03eb21303d39967a1a119b32dd744a0fa8986')) 
-- ORDER BY time desc LIMIT 10






