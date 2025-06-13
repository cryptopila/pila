# C++ to Go Module Mapping

This table lists the main C++ directories in the current project and suggests an equivalent Go module or third-party package that could be used during the migration.

| C++ Directory | Purpose | Suggested Go Replacement |
|---------------|---------|--------------------------|
| `coin` | Core blockchain logic, wallet functionality, networking and custom features such as ZeroTime, Incentive Reward and ChainBlender. | Use the `btcsuite/btcd` packages for Bitcoin-style networking and blockchain structures along with `golang.org/x/crypto` for cryptography. |
| `database` | Provides the O(1) block and transaction propagation layer, wrapping LevelDB and related database utilities. | The `goleveldb/leveldb` library can provide similar storage capabilities. |
| `crawler` | Discovers peers and monitors node availability using Boost.Asio networking. | Replace with Go's networking primitives or a P2P framework such as `libp2p` or the `btcd` peer package. |

These mappings are preliminary suggestions and may be adjusted as the migration progresses.
