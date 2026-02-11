package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
)

// ========== 你原来的结构精简版 ==========

type block struct {
	RollDelayOffset int64
	ConfirmedOffset int64
}

type evmNative struct {
	Parse     bool
	Decimal   int32
	TradeType model.TradeType
}

type evm struct {
	Network      string
	Block        block
	Native       evmNative
	Client       *http.Client
	AvgBlockTime int64
}

const (
	// ERC20 Transfer 事件签名
	evmTransferEvent = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

// 你需要补一个 rpcEndpoint 方法
func (e *evm) rpcEndpoint() string {
	// ⚠️ 换成你真实的 RPC
	return "https://bsc-testnet-rpc.publicnode.com"
}

// ========== 核心测试方法 ==========

func (e *evm) getBlockByNumberTest(blockNum int64) {
	payload := fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_getBlockByNumber",
		"params":["0x%x", true],
		"id":1
	}`, blockNum)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", e.rpcEndpoint(), bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Create request error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Client.Do(req)
	if err != nil {
		fmt.Println("Send request error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body error:", err)
		return
	}

	fmt.Println("======== Block Raw JSON ========")
	fmt.Println(string(body))
	fmt.Println("================================")

	blockJson := gjson.ParseBytes(body)
	result := blockJson.Get("result")

	transactions := result.Get("transactions").Array()
	blockNumberHex := result.Get("number").String()
	blockTimestampHex := result.Get("timestamp").String()

	blockTime := hexToTime(blockTimestampHex)

	fmt.Println("\n====== 解析 Native Transfer ======")
	parseNativeTransfer(transactions, blockNumberHex, blockTime)

	fmt.Println("\n====== 解析 ERC20 Event Transfer ======")
	parseEventTransfer(e, blockNumberHex)

}

func (e *evm) getBlockByNumberListTest(blockFromNum int64, blockToNum int64) {
	//payload := fmt.Sprintf(`{
	//	"jsonrpc":"2.0",
	//	"method":"eth_getBlockByNumber",
	//	"params":["0x%x", true],
	//	"id":1
	//}`, blockNum)

	items := make([]string, 0)
	for i := blockFromNum; i <= blockToNum; i++ {
		items = append(items, fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x%x",%t],"id":%d}`, i, e.Native.Parse, i))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//req, err := http.NewRequestWithContext(ctx, "POST", e.rpcEndpoint(), bytes.NewBuffer([]byte(payload)))
	req, err := http.NewRequestWithContext(ctx, "POST", e.rpcEndpoint(), bytes.NewBuffer([]byte(fmt.Sprintf(`[%s]`, strings.Join(items, ",")))))
	if err != nil {
		fmt.Println("Create request error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Client.Do(req)
	if err != nil {
		fmt.Println("Send request error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body error:", err)
		return
	}

	fmt.Println("======== Block Raw JSON ========")
	fmt.Println(string(body))
	fmt.Println("================================")

	for _, blockJson := range gjson.ParseBytes(body).Array() {
		//blockJson := gjson.ParseBytes(body)
		result := blockJson.Get("result")

		transactions := result.Get("transactions").Array()
		blockNumberHex := result.Get("number").String()
		blockTimestampHex := result.Get("timestamp").String()

		blockTime := hexToTime(blockTimestampHex)

		fmt.Println("\n====== 解析 Native Transfer ======")
		parseNativeTransfer(transactions, blockNumberHex, blockTime)

		fmt.Println("\n====== 解析 ERC20 Event Transfer ======")
		parseEventTransfer(e, blockNumberHex)
	}
}

// ================== Native Transfer 解析 ==================

func parseNativeTransfer(txs []gjson.Result, blockNumHex string, blockTime time.Time) {

	for _, tx := range txs {

		// 过滤非原生转账（input != 0x 表示调用合约）
		if tx.Get("input").String() != "0x" {
			continue
		}

		valueHex := tx.Get("value").String()
		if valueHex == "0x0" {
			continue
		}

		amount, ok := new(big.Int).SetString(valueHex[2:], 16)
		if !ok || amount.Sign() <= 0 {
			continue
		}

		fmt.Println("---- Native Transfer ----")
		fmt.Println("TxHash:", tx.Get("hash").String())
		fmt.Println("From:", tx.Get("from").String())
		fmt.Println("To:", tx.Get("to").String())
		fmt.Println("Amount (wei):", amount.String())
		fmt.Println("Block:", blockNumHex)
		fmt.Println("Time:", blockTime)
		fmt.Println()
	}
}

// ================== ERC20 Event Transfer 解析 ==================

func parseEventTransfer(e *evm, blockNumHex string) {

	payload := fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_getLogs",
		"params":[{
			"fromBlock":"%s",
			"toBlock":"%s",
			"topics":["%s"]
		}],
		"id":1
	}`, blockNumHex, blockNumHex, evmTransferEvent)

	resp, err := e.Client.Post(e.rpcEndpoint(), "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("eth_getLogs error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	data := gjson.ParseBytes(body)

	if data.Get("error").Exists() {
		fmt.Println("eth_getLogs RPC error:", data.Get("error").String())
		return
	}

	for _, logItem := range data.Get("result").Array() {

		topics := logItem.Get("topics").Array()
		if len(topics) < 3 {
			continue
		}

		from := "0x" + topics[1].String()[26:]
		to := "0x" + topics[2].String()[26:]

		amountHex := logItem.Get("data").String()
		amount, ok := new(big.Int).SetString(amountHex[2:], 16)
		if !ok {
			continue
		}

		fmt.Println("---- ERC20 Transfer ----")
		fmt.Println("TxHash:", logItem.Get("transactionHash").String())
		fmt.Println("Token Contract:", logItem.Get("address").String())
		fmt.Println("From:", from)
		fmt.Println("To:", to)
		fmt.Println("Amount (raw):", amount.String())
		fmt.Println("Block:", logItem.Get("blockNumber").String())
		fmt.Println()
	}
}

// ================== 工具函数 ==================

func hexToTime(hexStr string) time.Time {
	val, _ := new(big.Int).SetString(strings.TrimPrefix(hexStr, "0x"), 16)
	return time.Unix(val.Int64(), 0)
}

func main() {

	//e := &evm{
	//	Network: "ETH",
	//	Client: &http.Client{
	//		Timeout: 10 * time.Second,
	//	},
	//	Native: evmNative{
	//		Parse: true,
	//	},
	//}

	e := evm{
		Network: conf.Bsc,
		Block: block{
			ConfirmedOffset: 15,
		},
		Native: evmNative{
			Parse:     true,
			Decimal:   conf.BscBnbDecimals,
			TradeType: model.BscBnb,
		},
		Client: utils.NewHttpClient(),
	}

	// 🎯 你要测试的区块
	//blockNumber := int64(89669376)
	//
	//e.getBlockByNumberTest(blockNumber)

	blockFromNumber := int64(89669373)
	blockToNumber := int64(89669379)

	e.getBlockByNumberListTest(blockFromNumber, blockToNumber)
}
