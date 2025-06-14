-- Block Height
-- name: GetBlockHeight :one
SELECT number_int from block
where chain_id = $1
order by number_int desc
limit 1;

-- Block Insert
-- name: InsertBlock :exec
INSERT INTO block (chain_id, difficulty, hash, gas_limit, gas_used, miner, number,
                   number_int, parent_hash, timestamp, timestamp_int, created_at,
                   total_difficulty, transactions_root)
VALUES ($1, $2, $3, $4, $5, $6, $7,
        $8, $9, $10, $11, $12,
        $13, $14);

-- Transaction Insert
-- name: InsertTransaction :exec
INSERT INTO transaction (chain_id, block_hash, block_number, block_number_int, "from", "to",
                         gas, gas_int, gas_price, gas_price_int, hash, r, s, v, transaction_index,
                         value, value_int, nonce, input, contract_address, gas_used, gas_used_int,
                         status, type, timestamp, timestamp_int, created_at, coin_count, nft_count, erc20_count, erc721_count, erc1155_count)
VALUES ($1, $2, $3, $4, $5, $6,
        $7, $8, $9, $10, $11, $12, $13, $14, $15,
        $16, $17, $18, $19, $20, $21, $22,
        $23, $24, $25, $26, $27, $28, $29, $30, $31, $32);

-- Log Insert
-- name: InsertLog :exec
INSERT INTO log (chain_id, address, block_hash, block_number, block_number_int, data,
                 log_index, removed, topics, transaction_hash, transaction_index,
                 "from", "to", timestamp, timestamp_int, created_at)
VALUES ($1, $2, $3, $4, $5, $6,
        $7, $8, $9, $10, $11,
        $12, $13, $14, $15, $16);

-- Contract Insert
-- name: InsertContract :exec
INSERT INTO contract (chain_id, hash, name, symbol, decimals, total_supply, type, creator)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (hash, chain_id) DO UPDATE
SET name = EXCLUDED.name,
    symbol = EXCLUDED.symbol,
    decimals = EXCLUDED.decimals,
    total_supply = EXCLUDED.total_supply,
    type = EXCLUDED.type,
    creator = EXCLUDED.creator;

-- Coin Log Insert
-- name: InsertCoinLog :exec
INSERT INTO coin_log (chain_id, timestamp, timestamp_int, created_at, transaction_hash,
                      "from", "to", amount, gas, gas_int, gas_price, gas_price_int,
                      gas_used, gas_used_int)
VALUES ($1, $2, $3, $4, $5,
        $6, $7, $8, $9, $10, $11, $12,
        $13, $14);

-- ERC20 Log Insert
-- name: InsertERC20Log :exec
INSERT INTO erc20_log (chain_id, timestamp, timestamp_int, created_at, transaction_hash,
                       contract_address, "from", "to", amount, function, name, symbol)
VALUES ($1, $2, $3, $4, $5,
        $6, $7, $8, $9, $10, $11, $12);

-- ERC721 Log Insert
-- name: InsertERC721Log :exec
INSERT INTO erc721_log (chain_id, timestamp, timestamp_int, created_at, transaction_hash,
                        contract_address, "from", "to", token_id, function, name, symbol)
VALUES ($1, $2, $3, $4, $5,
        $6, $7, $8, $9, $10, $11, $12);

-- ERC1155 Log Insert
-- name: InsertERC1155Log :exec
INSERT INTO erc1155_log (chain_id, timestamp, timestamp_int, created_at, transaction_hash,
                         contract_address, "from", "to", token_id, amount, function, name, symbol)
VALUES ($1, $2, $3, $4, $5,
        $6, $7, $8, $9, $10, $11, $12, $13);

-- Wallet Insert
-- name: InsertWallet :exec
INSERT INTO wallet (chain_id, address)
VALUES ($1, $2) ON CONFLICT (chain_id, address) DO NOTHING;

-- Wallet Update Balance
-- name: UpsertWalletBalance :exec
INSERT INTO wallet (chain_id, address, balance)
VALUES ($1, $2, $3) ON CONFLICT (chain_id, address) DO UPDATE
SET balance = wallet.balance + EXCLUDED.balance;

-- ERC20 Balance UPSERT
-- name: UpsertERC20Balance :exec
INSERT INTO erc20_balance (chain_id, balance, hash, address)
VALUES ($1, $2, $3, $4)
ON CONFLICT (hash, chain_id, address) DO UPDATE
SET balance = erc20_balance.balance + EXCLUDED.balance;

-- -- ERC721 Balance INSERT
-- name: UpsertERC721Balance :exec
INSERT INTO erc721_balance (chain_id, hash, token_id, address)
VALUES ($1, $2, $3, $4)
ON CONFLICT (hash, token_id, chain_id)
DO UPDATE SET address = EXCLUDED.address;

-- ERC1155 Balance DELETE (소유권 이전 시)
-- name: SubtractERC1155Balance :exec
UPDATE erc1155_balance
SET amount = amount - $5
WHERE chain_id = $1 AND hash = $2 AND token_id = $3 AND address = $4
AND amount >= $5;

-- ERC1155 Balance UPSERT
-- name: UpsertERC1155Balance_Add :exec
INSERT INTO erc1155_balance (chain_id, hash, token_id, address, amount)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (hash, token_id, address, chain_id)
DO UPDATE SET amount = erc1155_balance.amount + EXCLUDED.amount;

-- name: CreateErc721 :exec
INSERT INTO erc721(chain_id, hash, token_id)
VALUES ($1, $2, $3)
ON CONFLICT (chain_id, hash, token_id)
DO NOTHING;

-- name: CreateErc1155 :exec
INSERT INTO erc721(chain_id, hash, token_id)
VALUES ($1, $2, $3)
    ON CONFLICT (chain_id, hash, token_id)
DO NOTHING;

-- name: UpdateContractType :exec
UPDATE contract
SET type = $1
WHERE chain_id = $2 AND hash = $3;