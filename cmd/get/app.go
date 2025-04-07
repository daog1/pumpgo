package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/daog1/pumpgo/generated/pump_amm"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
)

func FindAccountIndex(accountKeys []solana.PublicKey, targetAccount solana.PublicKey) int {
	for i, key := range accountKeys {
		if key == targetAccount {
			return int(i)
		}
	}
	return -1
}
func main() {
	//err := godotenv.Load("../../.env")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	rpcURL := os.Getenv("rpc")
	fmt.Printf("\n rpcURL %v", rpcURL)
	client := rpc.New(rpcURL)
	MaxSupportedTransactionVersion := uint64(0)

	txHash := solana.MustSignatureFromBase58("3uw5EEEroMAyus9R7eJqd3UY3VjeGZWSwCNRTyC92d4ZsrLeWjZCQXoUCyTrwvdPBcdq8VwZ6Hdk1s3tfCsD4GWd")
	out, err := client.GetTransaction(context.Background(), txHash, &rpc.GetTransactionOpts{
		Commitment:                     rpc.CommitmentFinalized,
		MaxSupportedTransactionVersion: &MaxSupportedTransactionVersion,
	})
	if err != nil {
		println(err.Error())
		return
	}
	evts, _ := pump_amm.DecodeEvents(out, solana.MustPublicKeyFromBase58("pAMMBay6oceH9fJKBRHGP5D4bD4sWpmSwMn52FMfXEA"), nil)
	for _, ev := range evts {
		fmt.Printf("\nev: %s\n", ev.Name)
		fmt.Printf("\nev: %v\n", ev)
	}

}
