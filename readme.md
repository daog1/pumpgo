
# 用 Go 监控 PumpSwap 流动性：三步搞定

最近我基于 Pump.fun 的 pAMMBay6oceH9fJKBRHGP5D4bD4sWpmSwMn52FMfXEA.json IDL，用我的 solana-anchor-go 项目写了个简单工具，实时监控 PumpSwap 的流动性事件（比如新池子创建），特别适合抢开盘或分析市场。我把过程浓缩成三步，代码清晰又实用。如果觉得不错，欢迎到 solana-anchor-go 点个 star 支持！

## 前提准备
环境：Go 1.18+，Solana RPC 节点（推荐helius申请个免费的账号,或者https://api.mainnet-beta.solana.com）。
依赖：安装我的 solana-anchor-go：
```
go install github.com/daog1/solana-anchor-go@latest
```
IDL：PumpSwap 的 IDL 文件（假设你已下载为 pAMMBay6oceH9fJKBRHGP5D4bD4sWpmSwMn52FMfXEA.json）。

## 步骤 1：生成 IDL 的 Go 接口
用 solana-anchor-go 把 pAMMBay6oceH9fJKBRHGP5D4bD4sWpmSwMn52FMfXEA.json 转换为 Go 代码，生成 PumpSwap 的接口。

运行以下命令：
```
solana-anchor-go -src=pAMMBay6oceH9fJKBRHGP5D4bD4sWpmSwMn52FMfXEA.json
```

这会生成 generated/pump_amm包,包里面的内容和idl文件中的指令定义、账号、事件一一对应.比如 CreatePoolEventEventData（新池创建事件）。



## 步骤 2：监听 PumpSwap 事件
PumpSwap 的程序 ID 是 pAMMBay6oceH9fJKBRHGP5D4bD4sWpmSwMn52FMfXEA。我们用 Go 订阅它的日志。

代码如下：

```
package main

import (
    "context"
    "log"

    "github.com/gagliardetto/solana-go"
    "github.com/gagliardetto/solana-go/rpc"
    "github.com/gagliardetto/solana-go/rpc/ws"
)

func main() {
    // 连接 WebSocket
    client, err := ws.Connect(context.Background(), "wss://api.mainnet-beta.solana.com")
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
        for _, logLine := range logResult.Value.Logs {
            log.Printf("新日志: %s", logLine)
        }
    }
}
```

运行后，它会实时打印 PumpSwap 的日志，比如 PoolCreatedEvent。

## 步骤 3：解析 PoolCreatedEvent
从 IDL 看，PoolCreatedEvent 是我们要监控的流动性事件，包含字段如 mint（代币地址）和 initial_liquidity（初始流动性）。用生成的 pumpswap.go 解析日志。

更新 for 循环：

```
	for {
		logResult, err := sub.Recv(context.Background())
		if err != nil {
			log.Printf("接收失败: %v", err)
			continue
		}
		evts, err := pump_amm.DecodeEvents(logResult.Value.Logs)
		for _, ev := range evts {
			if ev.Name == "CreatePoolEvent" {
				trEv := ev.Data.(*pump_amm.CreatePoolEventEventData)
				fmt.Printf("creator %s pool %s tx:%s\n", trEv.Creator, trEv.Pool, logResult.Value.Signature)
			}
		}
	}
```

完整代码
将 生成的idl解析代码和监听逻辑结合起来，保存为 app.go，然后：

```
go run main.go
```
输出示例：

```
2025/04/04 10:37:02 开始监听 PumpSwap 事件...
creator HwEfk69xeLu9MnLuqnrJEwn6iAkVg9QLqyVoLPdfSqGP pool Ed88ivvjJc2YqNQmyPXmwCjghhCe5kHa5wYAM5Z4Encs tx:EpytFhamRhk5QDdj85qYJ2yiirsg48gW6XuXrJiRjeqEGfSo1xK11NqL2TzBvXsSdVoJ832Kgx1pA7CyxNBGKHG
```
需要更多的事件类型和字段解析,可以在这个代码的基础上进行添加修改.
## 总结
其实写solana监听和evm监听没多大区别,主要是找对工具.

solana-anchor-go 是一个挺强大的golang solana链交互工具,只要是anchor开发的solana应用有idl都可以生成应用代码进行链交互.

这个工具根据我的应用需求进行了修复,发交易,解交易均没问题.目前anchor 0.30.1的idl都可以使用,anchor 以前的版本,可以按照,solana-anchor-go的文档进行先转换,再使用.


完整代码放在:
https://github.com/daog1/pumpgo
