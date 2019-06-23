package amg_test

import (
	"context"
	"fmt" 
	"math/big"
	"testing"

	bindings "github.com/RTradeLtd/AMG/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)


func Test_AMG_Token_Transfer_Burn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	auth, client := NewBlockchain(t)
	_, tx, tokenContract, err := bindings.DeployArenaMatchGold(auth, client)
	if err != nil {
		t.Fatal(err)
	}
	client.Commit()
	if _, err := bind.WaitDeployed(ctx, client, tx); err != nil {
		t.Fatal(err)
	}
	balance, err := tokenContract.BalanceOf(nil, auth.From)
	if err != nil {
		t.Fatal(err)
	}
	if balance.Int64() != 100000000 {
			t.Fatal("incorrect balance")
	}
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	auth2 := bind.NewKeyedTransactor(key)
	tx, err = tokenContract.Transfer(auth, auth2.From, big.NewInt(10000))
	if err != nil {
		t.Fatal(err)
	}
	client.Commit()
	if _, err := bind.WaitMined(ctx, client, tx); err != nil {
		t.Fatal(err)
	}
	balance, err = tokenContract.BalanceOf(nil, auth2.From)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(balance.Int64())
}

// NewBlockchain is used to generate a simulated blockchain
func NewBlockchain(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	auth := bind.NewKeyedTransactor(key)
	// https://medium.com/coinmonks/unit-testing-solidity-contracts-on-ethereum-with-go-3cc924091281
	gAlloc := map[common.Address]core.GenesisAccount{
		auth.From: {Balance: big.NewInt(10000000000)},
	}
	sim := backends.NewSimulatedBackend(gAlloc, 8000000)
	return auth, sim
}
