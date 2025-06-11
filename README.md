V(anilla)cash
===========

A decentralized currency for the internet.

[![Gitter](https://badges.gitter.im/cryptopila/pila.svg)](https://gitter.im/cryptopila/pila?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

This project is a codebase rewrite/enhancment using:
* [Peercoin](https://github.com/ppcoin/ppcoin) PoS & [Bitcoin](https://github.com/bitcoin/bitcoin) PoW consensus
* an UDP layer
* an off-chain transaction lock voting system called ZeroTime
* an Incentive Reward voting system
* a client-side blending system called ChainBlender


Dependencies:

* Boost == 1.53.0
* Berkeley DB == 6.1.29.NC
* OpenSSL == 1.0.2k

Linux: use https://github.com/cryptopila/pila-scripts
For the ongoing migration to Go, see [docs/migration_plan.md](docs/migration_plan.md).
The prototype Go module resides in the `go/` directory.
