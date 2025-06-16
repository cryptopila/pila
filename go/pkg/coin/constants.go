package coin

// Client and blockchain constants derived from the C++ sources.
const (
	VersionClientMajor    = 0
	VersionClientMinor    = 6
	VersionClientRevision = 0
	VersionClientBuild    = 4
	VersionClient         = 1000000*VersionClientMajor + 10000*VersionClientMinor + 100*VersionClientRevision + VersionClientBuild

	// TestNet indicates whether the library should operate in test network mode.
	TestNet = false

	VersionString = "0.6.0.4"
	ClientName    = "Pila"

	Coin int64 = 1000000
	Cent int64 = 10000

	MaxMoneySupply int64 = 30735360 * Coin

	MinTxFee       int64 = Cent / 20 // 0.05 * cent
	MinRelayTxFee  int64 = MinTxFee
	MinTxOutAmount int64 = MinTxFee

	ChainStartTime int64 = 1419310800

	CoinbaseMaturity            = 200
	CoinbaseMaturityTestNetwork = 1

	MaxMintProofOfStake int64 = Coin * 7 / 1000 // 0.007 * coin

	MinStakeAge = 60 * 60 * 8
	MaxStakeAge = 60 * 60 * 24 * 365

	MaxClockDrift = 2 * 60 * 60

	LockTimeThreshold uint64 = 500000000

	WorkAndStakeTargetSpacing = 200

	PowCutoffBlock = 2147483647 - 1
)
