package cli

import (
	"bytes"
	"flag"
	"fmt"
	"goblockchain/blockchain"
	"goblockchain/utils"
	"goblockchain/wallet"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Welcome to sovo's tiny blockchain system, usage is as follows:")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("All you need is to first create a blockchain and declare the owner.")
	fmt.Println("And then you can make transactions.")
	fmt.Println("In addition, don't forget to run mine function after transatcions are collected.")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("createwallet -refname REFNAME                       ----> Creates and save a wallet. The refname is optional.")
	fmt.Println("walletinfo -refname NAME -address Address           ----> Print the information of a wallet. At least one of the refname and address is required.")
	fmt.Println("walletsupdate                                       ----> Registrate and update all the wallets (especially when you have added an existed .wlt file).")
	fmt.Println("walletslist                                         ----> List all the wallets found (make sure you have run walletsupdate first).")
	fmt.Println("createblockchain -refname NAME -address ADDRESS     ----> Creates a blockchain with the owner you input")
	fmt.Println("balance -refname NAME -address ADDRESS              ----> Back the balance of the address you input")
	fmt.Println("blockchaininfo                                      ----> Prints the blocks in the chain")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> Make a transaction and put it into candidate block")
	fmt.Println("sendbyrefname -from NAME1 -to NAME2 -amount AMOUNT  ----> Make a transaction and put it into candidate block using refname.")
	fmt.Println("mine                                                ----> Mine and add a block to the chain")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------")
}

// create wallet
func (cli *CommandLine) createWallet(refname string) {
	newWallet := wallet.NewWallet()
	newWallet.Save()
	reList := wallet.LoadRefList()
	reList.BindRef(string(newWallet.Address()), refname)
	reList.Save()
	fmt.Println("Succeed in creating wallet.")
}

// walletinfoRefName
func (cli *CommandLine) walletInfoRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.walletInfo(address)
}

// walletinfo
func (cli *CommandLine) walletInfo(address string) {
	wlt := wallet.LoadWallet(address)
	refList := wallet.LoadRefList()
	fmt.Printf("Wallet address:%x\n", wlt.Address())
	fmt.Printf("Public Key:%x\n", wlt.PublicKey)
	fmt.Printf("Reference Name:%x\n", (*refList)[address])
}

// walletupdate
func (cli *CommandLine) walletUpdate() {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
	fmt.Println("Succeed in updating wallets")
}

// walletslist
func (cli *CommandLine) walletsList() {
	refList := wallet.LoadRefList()
	for address := range *refList {
		wlt := wallet.LoadWallet(address)
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Wallet address:%s\n", address)
		fmt.Printf("Public Key:%x\n", wlt.PublicKey)
		fmt.Printf("Reference Name:%s\n", (*refList)[address])
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println()
	}
}

// send by refname
func (cli *CommandLine) sendRefName(fromRefname, toRefname string, amount int) {
	refList := wallet.LoadRefList()
	fromAddress, err := refList.FindRef(fromRefname)
	utils.Handle(err)
	toAddress, err := refList.FindRef(toRefname)
	utils.Handle(err)
	cli.send(fromAddress, toAddress, amount)
}

// createblockchain
func (cli *CommandLine) createBlockChain(address string) {
	newChain := blockchain.InitBlockChain(utils.Address2PubHash([]byte(address)))
	defer newChain.DataBase.Close()
	fmt.Println("Finished creating a blockchain, and the owner is: ", address)
}

// createblockchainrefname
func (cli *CommandLine) createBlockChainRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.createBlockChain(address)
}

// balance address
func (cli *CommandLine) balance(address string) {
	chain := blockchain.ContinueBlockChain()
	defer chain.DataBase.Close()
	// 这里使用了go语言的defer关键字，其后的代码将会在函数运行结束前最后执行，也就是我们最后将关闭数据库。
	wlt := wallet.LoadWallet(address)
	balance, _ := chain.FindUTXOs(wlt.PublicKey)
	fmt.Printf("Address:%s, Blance:%d \n", address, balance)
}

// balance refname
func (cli *CommandLine) balanceRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.balance(address)
}

// blockchaininfo
func (cli *CommandLine) getBlockChainInfo() {
	chain := blockchain.ContinueBlockChain()
	defer chain.DataBase.Close()
	iterator := chain.Iterator()
	ogprevhash := chain.BackOgPrevHash()
	for {
		block := iterator.Next()
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Timestamp:%d\n", block.Timestamp)
		fmt.Printf("Previous hash:%x\n", block.PreHash)
		fmt.Printf("Transactions:%v\n", block.Transactions)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidataPow()))
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		if bytes.Equal(block.PreHash, ogprevhash) {
			break
		}
	}
}

// send -from FROADDRESS -to TOADDRESS -amount AMOUNT
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer chain.DataBase.Close()

	fromWallet := wallet.LoadWallet(from)
	tx, ok := chain.CreateTransaction(fromWallet.PublicKey, utils.Address2PubHash([]byte(to)), amount, fromWallet.Privtekey)
	if !ok {
		fmt.Println("Failed creating transaction")
		return
	}

	tp := blockchain.CreateTransactionPool()
	tp.AddTransaction(tx)
	tp.SaveFile()
	fmt.Println("Success!")
}

// mine
func (cli *CommandLine) mine() {
	chain := blockchain.ContinueBlockChain()
	defer chain.DataBase.Close()
	chain.RunMine()
	fmt.Println("Finish Mining")
}

// 使用go自带的flag库将各个命令注册
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// 操作
func (cli *CommandLine) Run() {
	cli.validateArgs()

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	walletInfoCmd := flag.NewFlagSet("walletinfo", flag.ExitOnError)
	wallestUpdate := flag.NewFlagSet("walletsupdate", flag.ExitOnError)
	walletsList := flag.NewFlagSet("walletlist", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	getBlockChainInfoCmd := flag.NewFlagSet("blockchaininfo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendByRefNameCmd := flag.NewFlagSet("sendbyrefname", flag.ExitOnError)
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)

	createWalletRefName := createWalletCmd.String("refname", "", "The refname of the wallet, and this is optimal")
	walletInfoRefName := walletInfoCmd.String("refname", "", "The refname of the wallet")
	walletInfoAddress := walletInfoCmd.String("address", "", "The address os the wallet")
	createBlockChainOwner := createBlockChainCmd.String("address", "", "The address refer to the owner of blockchain")
	createBlockChainRefNameOwner := createBlockChainCmd.String("refname", "", "The name refer to the owner of blockchain")
	balanceAddress := balanceCmd.String("address", "", "Who need to get balance amount")
	balanceRefName := balanceCmd.String("refname", "", "Who needs to get balance amount")
	sendFromRefName := sendByRefNameCmd.String("from", "", "Source refname")
	sendToRefName := sendByRefNameCmd.String("to", "", "Destination refname")
	sendByRefNameAmount := sendByRefNameCmd.Int("amount", 0, "Amount to send")
	sendFromAddress := sendCmd.String("from", "", "Source address")
	sendToAddress := sendCmd.String("to", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletinfo":
		err := walletInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletsupdate":
		err := wallestUpdate.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletslist":
		err := walletsList.Parse(os.Args[2:])
		utils.Handle(err)

	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "balance":
		err := balanceCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "blockchaininfo":
		err := getBlockChainInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "sendbyrefname":
		err := sendByRefNameCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "mine":
		err := mineCmd.Parse(os.Args[2:])
		utils.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(*createWalletRefName)
	}

	if walletInfoCmd.Parsed() {
		if *walletInfoAddress == "" {
			if *walletInfoRefName == "" {
				walletInfoCmd.Usage()
				runtime.Goexit()
			} else {
				cli.walletInfoRefName(*walletInfoRefName)
			}
		} else {
			cli.walletInfo(*walletInfoAddress)
		}
	}

	if wallestUpdate.Parsed() {
		cli.walletUpdate()
	}

	if walletsList.Parsed() {
		cli.walletsList()
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainOwner == "" {
			if *createBlockChainRefNameOwner == "" {
				createBlockChainCmd.Usage()
				runtime.Goexit()
			} else {
				cli.createBlockChainRefName(*createBlockChainRefNameOwner)
			}
		}
		cli.createBlockChain(*createBlockChainOwner)
	}

	if balanceCmd.Parsed() {
		if *balanceAddress == "" {
			if *balanceRefName == "" {
				balanceCmd.Usage()
				runtime.Goexit()
			} else {
				cli.balanceRefName(*balanceRefName)
			}
		} else {
			cli.balance(*balanceAddress)
		}
	}

	if sendCmd.Parsed() {
		if *sendFromAddress == "" || *sendToAddress == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFromAddress, *sendToAddress, *sendAmount)
	}

	if sendByRefNameCmd.Parsed() {
		if *sendFromRefName == "" || *sendToRefName == "" || *sendByRefNameAmount <= 0 {
			sendByRefNameCmd.Usage()
			runtime.Goexit()
		}
		cli.sendRefName(*sendFromRefName, *sendToRefName, *sendByRefNameAmount)
	}

	if getBlockChainInfoCmd.Parsed() {
		cli.getBlockChainInfo()
	}

	if mineCmd.Parsed() {
		cli.mine()
	}
}
