package test

import (
	"encoding/hex"
	"github.com/farukterzioglu/goBlockchain"
	"testing"
)

func TestNewTransaction(t *testing.T){
	nodeId := "3000"
	wallets, _ := goBlockchain.NewWallets(nodeId)
	fromAddress := wallets.CreateWallet()
	wallets.SaveToFile(nodeId)
	wallet := wallets.GetWallet(fromAddress)
	pubKeyHash := goBlockchain.HashPubKey(wallet.PublicKey)

	toAddress := wallets.CreateWallet()
	wallets.SaveToFile(nodeId)

	//Coinbase & reward transactions to sender
	txin := goBlockchain.TXInput{Txid: []byte{}, Vout: -1, PubKey: []byte("Coinbase")}
	txout := goBlockchain.NewTXOutput(10, fromAddress)
	transaction1 := goBlockchain.Transaction{Vin: []goBlockchain.TXInput{txin}, Vout: []goBlockchain.TXOutput{*txout}}
	transaction1.ID = transaction1.Hash()

	txin = goBlockchain.TXInput{Txid: []byte{}, Vout: -1, PubKey: []byte("Block reward")}
	txout = goBlockchain.NewTXOutput(10, fromAddress)
	transaction2 := goBlockchain.Transaction{Vin: []goBlockchain.TXInput{txin}, Vout: []goBlockchain.TXOutput{*txout}}
	transaction2.ID = transaction2.Hash()

	//Inputs
	prevTXs := make(map[string]goBlockchain.Transaction)
	var inputs []goBlockchain.TXInput

	//Check is output belongs to sender
	if transaction1.Vout[0].IsLockedWithKey(pubKeyHash){
		input1 := goBlockchain.TXInput{Txid: transaction1.ID, Vout: 0, PubKey: wallet.PublicKey}
		inputs = append(inputs, input1)
		prevTXs[hex.EncodeToString(transaction1.ID)] = transaction1
	}
	if transaction2.Vout[0].IsLockedWithKey(pubKeyHash){
		input2 := goBlockchain.TXInput{Txid: transaction2.ID, Vout: 0, PubKey: wallet.PublicKey}
		inputs = append(inputs, input2)
		prevTXs[hex.EncodeToString(transaction2.ID)] = transaction2
	}

	//Create output (locks with address)
	output := *goBlockchain.NewTXOutput(20, toAddress)

	//New transaction
	newTransaction := goBlockchain.Transaction{Vout: []goBlockchain.TXOutput{ output}, Vin:inputs, ID:nil}
	newTransaction.ID = newTransaction.Hash()

	//Sign
	newTransaction.Sign(wallet.PrivateKey, prevTXs)

	//Verify transaction
	if !newTransaction.Verify(prevTXs) {
		panic("transaction is not valid!")
	}
}
