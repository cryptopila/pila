package pkg

type Block struct {
    Hash     []byte
    PrevHash []byte
    Txns     []Transaction
}
