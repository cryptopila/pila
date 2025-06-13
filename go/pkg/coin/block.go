package coin

// Transaction represents a very simplified transaction placeholder.
type Transaction struct {
	ID string
}

// Block represents a minimal block structure used for the migration demo.
type Block struct {
	Hash         string
	Transactions []Transaction
}
