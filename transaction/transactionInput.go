package transaction

type TransactionInput struct {
	TransactionOutputId string
	UTXO                TransactionOutput
}

func (ti *TransactionInput) Create(toId string) {
	ti.TransactionOutputId = toId
	ti.UTXO = TransactionOutput{}
}
