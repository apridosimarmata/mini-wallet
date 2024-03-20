# Mini Wallet
#### submitted by Imam Aprido Simarmata

## Depedencies

### 1. Postgres

Used to store wallet & transaction data

### 2. Redis

Used to store tokens as key and wallet id as value\
Also provide distributed lock (in case locks are being used in multiple pod/machine) to prevent race condition on `deposits` and `withdrawal`\
Redlock (implemented with redsync) also provide TTL for each lock to prevent deadlock.\

## How to setup

goose -dir ./infrastructure/migrations/ create create_wallet_table sql