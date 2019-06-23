package amg_test

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"

	bindings "github.com/RTradeLtd/AMG/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
)

var (
	key     = `{"address":"7e4a2359c745a982a54653128085eac69e446de1","crypto":{"cipher":"aes-128-ctr","ciphertext":"eea2004c17292a9e94217bf53efbc31ff4ae62f3dd57f0938ab61c949a565dc1","cipherparams":{"iv":"6f6a7a89b556604940ac87ab1e78cfd1"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"8088e943ac0f37c8b4d01592d8bee96468853b6f1f13ca64d201cd68e7dc7b12"},"mac":"f856d734705f35e2acf854a44eb40796518730bd835ecaec01d1f3e7a7037813"},"id":"99e2cd49-4b51-4f01-b34c-aaa0efd332c3","version":3}`
	keypass = "password123"
)

func Test_AMG_Token_Transfer_Burn(t *testing.T) {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/22894d54ea1d4101ae951b14a69bc3ec")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	auth, err := bind.NewTransactor(strings.NewReader(key), keypass)
	if err != nil {
		log.Fatalf("Error unlocking account %v", err)
	}
	address, tx, tokenContract, err := bindings.DeployArenaMatchGold(auth, client)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bind.WaitDeployed(ctx, client, tx); err != nil {
		t.Fatal(err)
	}
	fmt.Println("contract address", address.String())
	balance, err := tokenContract.BalanceOf(nil, auth.From)
	if err != nil {
		t.Fatal(err)
	}
	tx, err = tokenContract.Transfer(auth, common.HexToAddress("0xe850da3D22061243E738FA07E20415e4B481257C"), balance)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bind.WaitMined(ctx, client, tx); err != nil {
		t.Fatal(err)
	}
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
