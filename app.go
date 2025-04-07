package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/daog1/pumpgo/generated/pump_amm"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/joho/godotenv"
)

func main() {
	// 连接 WebSocket
	// wss://api.mainnet-beta.solana.com,wss://solana-rpc.publicnode.com
	// add .env file rpc=wss://mainnet.helius-rpc.com/?api-key=
	//err := godotenv.Load(".env")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	rpcURL := os.Getenv("ws")
	client, err := ws.Connect(context.Background(), rpcURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	// PumpSwap 程序 ID
	programID := solana.MustPublicKeyFromBase58("pAMMBay6oceH9fJKBRHGP5D4bD4sWpmSwMn52FMfXEA")

	// 订阅日志
	sub, err := client.LogsSubscribeMentions(
		programID,
		rpc.CommitmentRecent,
	)
	if err != nil {
		//panic(err)
		return
	}
	if err != nil {
		log.Fatalf("订阅失败: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("开始监听 PumpSwap 事件...")

	// 接收日志
	for {
		logResult, err := sub.Recv(context.Background())
		if err != nil {
			log.Printf("接收失败: %v", err)
			continue
		}
		//fmt.Printf("ev: %s\n", logResult.Value.Logs)

		evts, err := pump_amm.DecodeEventsInLogs(logResult.Value.Logs)
		for _, ev := range evts {
			fmt.Printf("ev: %s\n", ev.Name)
			if ev.Name == "CreatePoolEvent" {
				trEv := ev.Data.(*pump_amm.CreatePoolEventEventData)
				fmt.Printf("creator %s pool %s tx:%s\n", trEv.Creator, trEv.Pool, logResult.Value.Signature)
			}
		}
	}
}
