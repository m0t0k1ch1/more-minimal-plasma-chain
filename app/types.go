package app

const (
	BlockTypeNormal  = "normal"
	BlockTypeDeposit = "deposit"
)

type BlockType string

func (bt BlockType) IsValid() bool {
	return bt.IsNormal() || bt.IsDeposit()
}

func (bt BlockType) IsNormal() bool {
	return bt == BlockTypeNormal
}

func (bt BlockType) IsDeposit() bool {
	return bt == BlockTypeDeposit
}
