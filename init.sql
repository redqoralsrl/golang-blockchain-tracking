-- command + alt + l

create table block
(
    id                serial primary key,
    chain_id          numeric      not null,
    difficulty        varchar(255),
    hash              varchar(255) not null,
    gas_limit         varchar(255),
    gas_used          varchar(255),
    miner             varchar(42),
    number            varchar(255) not null,
    number_int        numeric      not null,
    parent_hash       varchar(255),
    timestamp         varchar(255) not null,
    timestamp_int     numeric      not null,
    created_at        timestamp    not null,
    total_difficulty  varchar(255),
    transactions_root varchar(255),
    unique (hash, chain_id)
);

create table transaction
(
    id                serial primary key,
    chain_id          numeric           not null,
    block_hash        varchar(255)      not null,
    block_number      varchar(255)      not null,
    block_number_int  numeric           not null,
    "from"            varchar(42),
    "to"              varchar(42),
    gas               varchar(255),
    gas_int           numeric,
    gas_price         varchar(255),
    gas_price_int     numeric,
    hash              varchar(255)      not null,
    r                 varchar(255),
    s                 varchar(255),
    v                 varchar(255),
    transaction_index varchar(255),
    value             varchar(255),
    value_int         numeric,
    nonce             varchar(255),
    input             text,
    contract_address  varchar(255),
    gas_used          varchar(255),
    gas_used_int      numeric,
    status            varchar(255),
    type              varchar(255),
    timestamp         varchar(255)      not null,
    timestamp_int     numeric           not null,
    created_at        timestamp         not null,
    coin_count        numeric default 0 not null,
    nft_count         numeric default 0 not null,
    erc20_count       numeric default 0 not null,
    erc721_count      numeric default 0 not null,
    erc1155_count     numeric default 0 not null,
    unique (hash, chain_id)
);

create table log
(
    id                serial primary key,
    chain_id          numeric      not null,
    address           varchar(42),
    block_hash        varchar(255) not null,
    block_number      varchar(255) not null,
    block_number_int  numeric,
    data              text,
    log_index         varchar(255),
    removed           boolean      not null,
    topics            text[],
    transaction_hash  varchar(255) not null,
    transaction_index varchar(255) not null,
    "from"            varchar(42),
    "to"              varchar(42),
    timestamp         varchar(255) not null,
    timestamp_int     numeric      not null,
    created_at        timestamp    not null
);

create table contract
(
    id             serial primary key,
    chain_id       numeric      not null,
    hash           varchar(255) not null,
    name           text,
    symbol         text,
    decimals       integer,
    total_supply   numeric,
    type           varchar(255), -- erc20: 20 erc721: 721 erc1155: 1155
    creator        varchar(42),
    logo_url       text,
    background_url text,
    unique (hash, chain_id)
);

create table coin_log
(
    id               serial primary key,
    chain_id         numeric      not null,
    timestamp        varchar(255) not null,
    timestamp_int    numeric      not null,
    created_at       timestamp    not null,
    transaction_hash varchar(255),
    "from"           varchar(42),
    "to"             varchar(42),
    amount           numeric,
    gas              varchar(255),
    gas_int          numeric,
    gas_price        varchar(255),
    gas_price_int    numeric,
    gas_used         varchar(255),
    gas_used_int     numeric
);

create table erc20_log
(
    id               serial primary key,
    chain_id         numeric      not null,
    timestamp        varchar(255) not null,
    timestamp_int    numeric      not null,
    created_at       timestamp    not null,
    transaction_hash varchar(255) not null,
    contract_address varchar(255) not null,
    "from"           varchar(42)  not null,
    "to"             varchar(42)  not null,
    amount           numeric      not null,
    function         varchar(255) not null,
    name             text,
    symbol           text
);

create table erc721_log
(
    id               serial primary key,
    chain_id         numeric      not null,
    timestamp        varchar(255) not null,
    timestamp_int    numeric      not null,
    created_at       timestamp    not null,
    transaction_hash varchar(255) not null,
    contract_address varchar(255) not null,
    "from"           varchar(42)  not null,
    "to"             varchar(42)  not null,
    token_id         numeric      not null,
    function         varchar(255),
    name             text,
    symbol           text
);

create table erc1155_log
(
    id               serial primary key,
    chain_id         numeric      not null,
    timestamp        varchar(255) not null,
    timestamp_int    numeric      not null,
    created_at       timestamp    not null,
    transaction_hash varchar(255) not null,
    contract_address varchar(255),
    "from"           varchar(42),
    "to"             varchar(42),
    token_id         numeric,
    amount           numeric,
    function         varchar(255),
    name             text,
    symbol           text
);

create table wallet
(
    id       serial primary key,
    chain_id numeric           not null,
    address  varchar(42)       not null,
    balance  numeric default 0 not null,
    unique (chain_id, address)
);

create table erc20_balance
(
    id       serial primary key,
    chain_id numeric           not null,
    balance  numeric default 0 not null,
    hash     varchar(255)      not null,
    address  varchar(42)       not null,
    unique (hash, chain_id, address)
);

create table erc721_balance
(
    id       serial primary key,
    chain_id numeric      not null,
    hash     varchar(255) not null,
    token_id numeric      not null,
    address  varchar(42)  not null,
    unique (hash, token_id, chain_id)
);

create table erc1155_balance
(
    id       serial primary key,
    chain_id numeric      not null,
    hash     varchar(255) not null,
    token_id numeric      not null,
    address  varchar(42)  not null not null,
    amount   numeric      not null,
    unique (hash, token_id, address, chain_id)
);

create table erc721
(
    id        serial primary key,
    chain_id  numeric      not null,
    hash      varchar(255) not null,
    token_id  numeric      not null,
    url text,
    image_url text,
    unique (chain_id, hash, token_id)
);

create table erc1155
(
    id        serial primary key,
    chain_id  numeric      not null,
    hash      varchar(255) not null,
    token_id  numeric      not null,
    url text,
    image_url text,
    unique (chain_id, hash, token_id)
);

INSERT INTO wallet (chain_id, address, balance)
VALUES (8989, '0xc90cf9b42cbd122810d3b4661c3e1083d24bb97b', 456719261665419348619673600),
       (8989, '0xca59d509d074d28c3fa89f976aaf7dc3302b327c', 457796877875121402945740800),
       (8989, '0xaee6ea43376868951542a32b7629b83bd07f03ae', 1469146820925806344693760000),
       (8989, '0xab63b99c7bbbaa50d5e172e57242c0a3a0cae3fc', 568172893402942721024000000),
       (8989, '0x602928e12cf7c77c61f50ae9437b0a082725c1b3', 935000000000000000000000000),
       (8989, '0xbb16fae5808ff3dd77557448280a900f72b6991a', 192000000000000000000000);

INSERT INTO wallet (chain_id, address, balance)
VALUES (898989, '0xbb16fae5808ff3dd77557448280a900f72b6991a', 192000000000000000000000);
