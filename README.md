# BITCOIN

# Bitcoin RPC Explorer (Go)

A lightweight Bitcoin blockchain explorer built in Go using Bitcoin Core JSON-RPC.
This project was developed in a step-by-step bootcamp style, progressively building from basic RPC calls to a full blockchain inspection tool.

---

## Overview

This project connects to a **Bitcoin Core regtest node** and allows you to:

* Query wallet balances
* List wallet transactions
* Decode raw transactions
* Inspect blocks
* Fetch blockchain metadata
* Interact with multiple wallets (alice, bob, etc.)

It acts as a **minimal blockchain explorer backend**.

---

## Features

### 1. Wallet Balance

Fetch balances for specific wallets:

* `getbalance`

### 2. Transaction Listing

View wallet transaction history:

* `listtransactions`

Includes:

* Category (receive/send/mining)
* Amount
* Transaction ID
* Confirmations

---

### 3. Transaction Decoder

Decode raw transactions:

* `getrawtransaction`

Shows:

* VIN (inputs)
* VOUT (outputs)
* CoinBase detection (mining rewards)
* Address extraction

Handles both:

* `addresses[]`
* fallback `address`

---

### 4. Blockchain Info

Fetch network status:

* `getblockchaininfo`

Displays:

* Chain type (regtest/testnet/mainnet)
* Block height
* Difficulty

---

### 5. Block Explorer

Inspect blocks using:

* `getblockhash`
* `getblock`

Features:

* Auto-fetch best block if no hash provided
* Block height
* Timestamp
* Transaction count
* Transaction list

---

## Architecture

```
Go Application
   ↓
RPC Client (HTTP + JSON)
   ↓
Bitcoin Core (regtest)
   ↓
Blockchain Data
```

---

## RPC Helper Function

All blockchain interactions go through a reusable RPC wrapper:

* Handles authentication (rpcuser / rpcpassword)
* Supports wallet-specific routing
* Decodes JSON responses into Go structs

---

## Example Workflow

1. Start Bitcoin Core in regtest:

```bash
bitcoind -regtest -daemon
```

2. Create/load wallets:

```bash
bitcoin-cli -regtest createwallet "alice"
bitcoin-cli -regtest createwallet "bob"
```

3. Run the Go app:

```bash
go run .
```

---

## Sample Output

```
=== Wallet: alice ===
Balance: 89.9999859 BTC

=== Transactions for alice ===
IN  +50.00000000 BTC | 9 confs
        TXID: dbe6ff95...
OUT -10.00000000 BTC | 1 confs
        TXID: a4171cb8...

=== BLOCK INFO ===
Height: 102
Hash: 000000...
Time: 1710000000
Tx count: 1
```

---

## Key Concepts Learned

### Bitcoin RPC

* JSON-RPC communication with Bitcoin Core
* Methods like `getblock`, `getbalance`, `listtransactions`

### Wallet Management

* Loading wallets into memory
* Using `/wallet/<name>` RPC routing

### Blockchain Data Structures

* VIN / VOUT model (UTXO system)
* ScriptPubKey decoding
* Coinbase transactions

### Go Concepts

* Struct mapping to JSON
* Error handling patterns
* Reusable RPC client design

---

## Common Issues Solved

* Wallet not loaded errors → solved via `loadwallet`
* Multiple wallet selection → solved via `/wallet/<name>`
* JSON mismatches (address vs addresses[])
* Regtest configuration issues
* RPC authentication setup

---

## Requirements

* Go 1.18+
* Bitcoin Core (regtest mode enabled)
* RPC enabled in `bitcoin.conf`

Example config:

```
server=1
regtest=1
rpcuser=place_username_here
rpcpassword=place_password_here
```

---

## Future Improvements

* REST API wrapper
* Mempool explorer
* UTXO tracker
* Frontend dashboard
* Real-time block updates (WebSockets)
* Transaction graph visualization

---

## Author

Built as part of a Bitcoin RPC bootcamp project — evolving from basic RPC calls to a full blockchain inspection tool in Go.
