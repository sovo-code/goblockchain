package constcoe

//pow困难值设置
const (
	Difficuty = 12
	InitCoin  = 1000
	// for command manager and remember
	TransactionPoolFile = "./tmp/transaction_pool.data"
	BCPath              = "./tmp/blocks"
	BCFile              = "./tmp/blocks/MANIFEST"
	// for public secret and private secret and wallets
	ChecksumLength = 4
	NetWorkVersion = byte(0x00)
	Wallets        = "./tmp/wallets/"
	WalletsRefList = "./tmp/ref_list/"
)
