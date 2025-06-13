# Migration to Go - Initial Plan

This document outlines the initial ideas for migrating the current C++ codebase to Go.

## Overview of the Current Project
- Rewrite/enhancement based on Peercoin PoS and Bitcoin PoW
- Contains custom features: UDP layer, ZeroTime transaction locking, Incentive Reward, ChainBlender
- Relies heavily on Boost, Berkeley DB and OpenSSL
- Codebase has roughly 100k lines of C/C++ spread across `coin`, `database` and `crawler`.
  See [module_mapping.md](module_mapping.md) for a proposed mapping of these directories to Go packages.

## Probability Estimate
Due to the large size and complexity of the existing code, the migration to Go is estimated to have a **low probability of success (around 20%)** without significant resources and a dedicated team.

## Updated Analysis

Since the initial plan was drafted, only a minimal Go stub has been added under
`cmd/pila`. The bulk of the project remains C++ and still relies on Boost and
LevelDB. Migrating the UDP layer and custom consensus logic will require
carefully mapping the existing abstractions to Go packages. The PoW/PoS
implementation is tightly coupled with the current networking code, so a
straight rewrite is unlikely to succeed without further refactoring of the C++
codebase.

## Suggested Approach
1. **Evaluation**: map all current modules and dependencies.
2. **Architecture Definition**: choose Go packages to replace C++ components:
   - `golang.org/x/crypto` for cryptography
   - `goleveldb` for storage
   - a P2P network layer inspired by `btcd`
3. **Modeling**: design Go structures for blocks, transactions, wallets and networking.
4. **Networking**: implement P2P communication using Go libraries, replacing Boost Asio code.
5. **Consensus**: reimplement PoW/PoS and reward logic in Go.
6. **Custom Features**: port ZeroTime, ChainBlender and Incentive.
7. **RPC Services**: expose wallet and blockchain operations via Go's HTTP or gRPC.
8. **Build & Tests**: use `go mod`, write unit and integration tests for each module.

## Immediate Task List

- [ ] Evaluate dependencies in `coin`, `database` and `crawler`.
- [ ] Prototype block and transaction structures in Go.
- [ ] Implement a basic P2P handshake.
- [ ] Stub LevelDB interactions for the `database` package.
- [ ] Replace Boost-based networking in the `crawler` component.

## What Can Be Ignored or Replaced
- Custom cryptographic utility implementations can be replaced by Go's `crypto` packages.
- LevelDB C++ code can be replaced by the Go wrapper `goleveldb`.
- The O(1) propagation layer and crawler can be substituted with existing P2P frameworks used by `btcd`.
- C++ build scripts for Windows/Unix and JNI bindings can be dropped in favour of Go's cross-platform build tools.

## Next Steps

The following tasks will help advance the migration:

1. Audit each C++ directory and document which features must be ported.
2. Extend the Go module with `btcsuite/btcd` to replace the custom P2P layer.
3. Prototype the P2P handshake alongside basic block and transaction types in Go.
4. Implement a small LevelDB-backed storage layer using `goleveldb`.
5. Set up CI scripts running `go vet` and unit tests to verify the new code.

## Verifying the Go Environment

After running `go mod init pila`, a tiny executable can be built to
confirm that Go is installed correctly. From the repository root run:

```bash
go run ./cmd/pila
```

This should print `pila go stub running` confirming the environment is ready
for further development.

