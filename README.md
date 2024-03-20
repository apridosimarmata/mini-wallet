# Mini Wallet
#### submitted by Imam Aprido Simarmata

## Depedencies

### 1. Postgres

Used to store wallet & transaction data

### 2. Redis

Used to store tokens as key and wallet id as value
Also provide distributed lock to prevent race condition on `deposits` and `withdrawal`

## How to setup

goose -dir ./infrastructure/migrations/ create create_wallet_table sql