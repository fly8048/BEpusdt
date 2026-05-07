package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/shopspring/decimal"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/utils"
	"github.com/v03413/tronprotocol/api"
	"github.com/v03413/tronprotocol/core"
	"google.golang.org/grpc"
)

var gasFreeUsdtTokenAddress = []byte{0xa6, 0x14, 0xf8, 0x03, 0xb6, 0xfd, 0x78, 0x09, 0x86, 0xa4, 0x2c, 0x78, 0xec, 0x9c, 0x7f, 0x77, 0xe6, 0xde, 0xd1, 0x3c}
var gasFreeOwnerAddress = []byte{0x41, 0x3b, 0x41, 0x50, 0x50, 0xb1, 0xe7, 0x9e, 0x38, 0x50, 0x7c, 0xb6, 0xe4, 0x8d, 0xac, 0xc2, 0x27, 0xaf, 0xfd, 0xd5, 0x0c}
var gasFreeContractAddress = []byte{0x41, 0x39, 0xdd, 0x12, 0xa5, 0x4e, 0x2b, 0xab, 0x7c, 0x82, 0xaa, 0x14, 0xa1, 0xe1, 0x58, 0xb3, 0x42, 0x63, 0xd2, 0xd5, 0x10}
var usdtTrc20ContractAddress = []byte{0x41, 0xa6, 0x14, 0xf8, 0x03, 0xb6, 0xfd, 0x78, 0x09, 0x86, 0xa4, 0x2c, 0x78, 0xec, 0x9c, 0x7f, 0x77, 0xe6, 0xde, 0xd1, 0x3c}
var usdcTrc20ContractAddress = []byte{0x41, 0x34, 0x87, 0xb6, 0x3d, 0x30, 0xb5, 0xb2, 0xc8, 0x7f, 0xb7, 0xff, 0xa8, 0xbc, 0xfa, 0xde, 0x38, 0xea, 0xac, 0x1a, 0xbe}

func base58CheckEncode(input []byte) string {
	checksum := chainhash.DoubleHashB(input)
	checksum = checksum[:4]

	input = append(input, checksum...)

	return base58.Encode(input)
}

func base58CheckDecode(address string) ([]byte, error) {
	// 1️⃣ Base58 decode
	decoded := base58.Decode(address)
	if len(decoded) != 25 {
		return nil, fmt.Errorf("invalid address length")
	}

	// 2️⃣ 拆分 payload 和 checksum
	payload := decoded[:21]
	checksum := decoded[21:]

	// 3️⃣ 重新计算 checksum
	hash := chainhash.DoubleHashB(payload)
	expectedChecksum := hash[:4]

	// 4️⃣ 校验 checksum
	if !bytes.Equal(checksum, expectedChecksum) {
		return nil, fmt.Errorf("checksum mismatch")
	}

	// 5️⃣ 返回 21 字节真实地址
	return payload, nil
}

func main() {

	//gasFreeUsdtTokenAddressResult := base58CheckEncode(gasFreeUsdtTokenAddress)
	//fmt.Println(gasFreeUsdtTokenAddressResult)
	//
	//gasFreeOwnerAddressResult := base58CheckEncode(gasFreeOwnerAddress)
	//fmt.Println(gasFreeOwnerAddressResult)
	//
	//gasFreeContractAddressResult := base58CheckEncode(gasFreeContractAddress)
	//fmt.Println(gasFreeContractAddressResult)
	//
	//usdtTrc20ContractAddressResult := base58CheckEncode(usdtTrc20ContractAddress)
	//fmt.Println(usdtTrc20ContractAddressResult)
	//
	//usdcTrc20ContractAddressResult := base58CheckEncode(usdcTrc20ContractAddress)
	//fmt.Println(usdcTrc20ContractAddressResult)
	//
	req, err := base58CheckDecode("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	if err != nil {
		fmt.Println("Create request error:", err)
	}
	printInline(req)

	//reqT, errT := base58CheckDecode("TU1ntBzpGPp7GJkzxLTKwYsneJ9JKUmBCK")
	//if errT != nil {
	//	fmt.Println("Create request error:", errT)
	//}
	//printInline(reqT)

	reqT, errT := base58CheckDecode("TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf")
	if errT != nil {
		fmt.Println("Create request error:", errT)
	}
	printInline(reqT)

	//DebugParseBlock(64858786)
	//DebugParseBlock(64873094)

}

func printInline(b []byte) {
	fmt.Print("[]byte{")
	for i, v := range b {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("0x%02x", v)
	}
	fmt.Println("}")
}

func client() (*grpc.ClientConn, error) {

	conn, err := utils.NewTronGrpcClient(tronGrpcEndpoint, "")
	if err != nil {

		return nil, fmt.Errorf("连接失败: %w", err)
	}

	return conn, nil
}

//func base58CheckEncode(input []byte) string {
//	checksum := chainhash.DoubleHashB(input)
//	checksum = checksum[:4]
//
//	input = append(input, checksum...)
//
//	return base58.Encode(input)
//}

func parseTrc20ContractTransferFrom(data []byte) (string, string, *big.Int) {
	if len(data) != 100 {

		return "", "", nil
	}

	from := base58CheckEncode(append([]byte{0x41}, data[16:36]...))
	to := base58CheckEncode(append([]byte{0x41}, data[48:68]...))
	amount := big.NewInt(0).SetBytes(data[68:100])

	return from, to, amount
}

func parseTrc20ContractTransfer(data []byte) (string, *big.Int) {
	if len(data) != 68 {

		return "", nil
	}

	receiver := base58CheckEncode(append([]byte{0x41}, data[16:36]...))
	amount := big.NewInt(0).SetBytes(data[36:68])

	return receiver, amount
}

func DebugParseBlock(blockNum int64) {
	conn, err := client()
	if err != nil {
		log.Task.Error("grpc.NewClient", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	block, err := api.NewWalletClient(conn).GetBlockByNum2(ctx, &api.NumberMessage{
		Num: blockNum,
	})
	if err != nil {
		log.Task.Error("GetBlockByNum2", err)
		return
	}

	fmt.Println("====================================================")
	fmt.Println("Block Number:", blockNum)
	fmt.Println("Timestamp   :", time.UnixMilli(block.GetBlockHeader().GetRawData().GetTimestamp()))
	fmt.Println("Tx Count    :", len(block.GetTransactions()))
	fmt.Println("====================================================")

	for _, trans := range block.GetTransactions() {

		if !trans.Result.Result {
			continue
		}

		txID := hex.EncodeToString(trans.Txid)
		raw := trans.GetTransaction()

		for _, contract := range raw.GetRawData().GetContract() {

			switch contract.GetType() {

			// ========================
			// TRX 转账
			// ========================
			case core.Transaction_Contract_TransferContract:
				var tc core.TransferContract
				if err := contract.GetParameter().UnmarshalTo(&tc); err != nil {
					continue
				}

				fmt.Println("---- TRX Transfer ----")
				fmt.Println("TxID      :", txID)
				fmt.Println("From      :", base58CheckEncode(tc.OwnerAddress))
				fmt.Println("To        :", base58CheckEncode(tc.ToAddress))
				fmt.Println("Amount    :", decimal.NewFromBigInt(
					new(big.Int).SetInt64(tc.Amount), -6))
				fmt.Println()

			// ========================
			// TRC20
			// ========================
			case core.Transaction_Contract_TriggerSmartContract:

				var sc core.TriggerSmartContract
				if err := contract.GetParameter().UnmarshalTo(&sc); err != nil {
					continue
				}

				data := sc.GetData()
				if len(data) < 4 {
					continue
				}

				contractAddr := base58CheckEncode(sc.GetContractAddress())

				// transfer
				if bytes.Equal(data[:4], []byte{0xa9, 0x05, 0x9c, 0xbb}) {

					receiver, amount := parseTrc20ContractTransfer(data)
					if amount == nil {
						continue
					}

					fmt.Println("---- TRC20 Transfer ----")
					fmt.Println("TxID        :", txID)
					fmt.Println("Token       :", contractAddr)
					fmt.Println("From        :", base58CheckEncode(sc.OwnerAddress))
					fmt.Println("To          :", receiver)
					fmt.Println("Amount(raw) :", decimal.NewFromBigInt(amount, conf.UsdtTronDecimals))
					fmt.Println()

				}

				// transferFrom
				if bytes.Equal(data[:4], []byte{0x23, 0xb8, 0x72, 0xdd}) {

					from, to, amount := parseTrc20ContractTransferFrom(data)
					if amount == nil {
						continue
					}

					fmt.Println("---- TRC20 TransferFrom ----")
					fmt.Println("TxID        :", txID)
					fmt.Println("Token       :", contractAddr)
					fmt.Println("From        :", from)
					fmt.Println("To          :", to)
					fmt.Println("Amount(raw) :", amount.String())
					fmt.Println()
				}

			// ========================
			// 资源代理
			// ========================
			case core.Transaction_Contract_DelegateResourceContract:
				var rc core.DelegateResourceContract
				if err := contract.GetParameter().UnmarshalTo(&rc); err != nil {
					continue
				}

				fmt.Println("---- Delegate Resource ----")
				fmt.Println("TxID    :", txID)
				fmt.Println("From    :", base58CheckEncode(rc.OwnerAddress))
				fmt.Println("To      :", base58CheckEncode(rc.ReceiverAddress))
				fmt.Println("Balance :", rc.Balance)
				fmt.Println()

			case core.Transaction_Contract_UnDelegateResourceContract:
				var rc core.UnDelegateResourceContract
				if err := contract.GetParameter().UnmarshalTo(&rc); err != nil {
					continue
				}

				fmt.Println("---- UnDelegate Resource ----")
				fmt.Println("TxID    :", txID)
				fmt.Println("From    :", base58CheckEncode(rc.OwnerAddress))
				fmt.Println("To      :", base58CheckEncode(rc.ReceiverAddress))
				fmt.Println("Balance :", rc.Balance)
				fmt.Println()
			}
		}
	}
}

const tronGrpcEndpoint = "grpc.nile.trongrid.io:50051" // 主网

//func main() {
//
//	if len(os.Args) < 2 {
//		fmt.Println("Usage: go run main.go <blockNumber>")
//		return
//	}
//
//	var blockNum int64
//	fmt.Sscan(os.Args[1], &blockNum)
//
//	conn, err := grpc.Dial(tronGrpcEndpoint, grpc.WithInsecure())
//	if err != nil {
//		panic(err)
//	}
//	defer conn.Close()
//
//	client := api.NewWalletClient(conn)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	block, err := client.GetBlockByNum2(ctx, &api.NumberMessage{
//		Num: blockNum,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println("======================================")
//	fmt.Println("Block Number:", blockNum)
//	fmt.Println("Timestamp   :", time.UnixMilli(block.GetBlockHeader().GetRawData().GetTimestamp()))
//	fmt.Println("Tx Count    :", len(block.GetTransactions()))
//	fmt.Println("======================================")
//
//	for _, trans := range block.GetTransactions() {
//
//		if !trans.Result.Result {
//			continue
//		}
//
//		txID := hex.EncodeToString(trans.Txid)
//		raw := trans.GetTransaction()
//
//		for _, contract := range raw.GetRawData().GetContract() {
//
//			switch contract.GetType() {
//
//			// TRX 转账
//			case core.Transaction_Contract_TransferContract:
//				var tc core.TransferContract
//				if err := contract.GetParameter().UnmarshalTo(&tc); err != nil {
//					continue
//				}
//
//				fmt.Println("---- TRX Transfer ----")
//				fmt.Println("TxID :", txID)
//				fmt.Println("From :", base58CheckEncode(tc.OwnerAddress))
//				fmt.Println("To   :", base58CheckEncode(tc.ToAddress))
//				fmt.Println("Amt  :", tc.Amount/1_000_000, "TRX")
//				fmt.Println()
//
//			// TRC20
//			case core.Transaction_Contract_TriggerSmartContract:
//
//				var sc core.TriggerSmartContract
//				if err := contract.GetParameter().UnmarshalTo(&sc); err != nil {
//					continue
//				}
//
//				data := sc.GetData()
//				if len(data) < 4 {
//					continue
//				}
//
//				// transfer
//				if bytes.Equal(data[:4], []byte{0xa9, 0x05, 0x9c, 0xbb}) {
//
//					if len(data) != 68 {
//						continue
//					}
//
//					receiver := base58CheckEncode(
//						append([]byte{0x41}, data[16:36]...))
//
//					amount := big.NewInt(0).SetBytes(data[36:68])
//
//					fmt.Println("---- TRC20 Transfer ----")
//					fmt.Println("TxID   :", txID)
//					fmt.Println("Token  :", base58CheckEncode(sc.GetContractAddress()))
//					fmt.Println("From   :", base58CheckEncode(sc.OwnerAddress))
//					fmt.Println("To     :", receiver)
//					fmt.Println("Amount :", amount.String())
//					fmt.Println()
//				}
//			}
//		}
//	}
//}
//
////func base58CheckEncode(input []byte) string {
////	checksum := chainhash.DoubleHashB(input)
////	checksum = checksum[:4]
////	input = append(input, checksum...)
////	return base58.Encode(input)
////}
