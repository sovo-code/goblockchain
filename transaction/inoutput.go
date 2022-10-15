package transaction

type TxOutput struct {
	Value      int    //转出的资产值
	HashPubKey []byte //目标地址
}

type TxInput struct {
	TxID   []byte //前置交易信息
	OutIdx int    //计数
	PubKey []byte //来源地址
	Sig    []byte //签名
}
