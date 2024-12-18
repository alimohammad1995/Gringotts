package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/holiman/uint256"
	"gringotts/models"
	"gringotts/provider"
	"gringotts/utils"
	"math"
	"math/big"
	"strings"
	"time"
)

func HandleTest(c *fiber.Ctx) error {
	TestSol()
	testTron2()
	return testTron(c)
}

func TestSol() {
	swap := provider.Jupiter{}

	x, err := swap.GetSwap(
		&provider.SwapParams{
			Chain:       models.Solana,
			FromToken:   models.GetToken(models.Solana, ""),
			ToToken:     models.GetToken(models.Solana, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
			Amount:      uint256.NewInt(uint64(math.Pow10(9))),
			Recipient:   "H6ieTjyWqcRFv1RwD9LghEFDxVtMaEMgVbgnYrdHMjr5",
			SlippageBPS: 100,
		},
	)

	fmt.Println(x, err)

	return
}

func testTron2() {
	swap := provider.SunSwap{}

	x, err := swap.GetSwap(&provider.SwapParams{
		Chain:       models.TRON,
		FromToken:   models.GetToken(models.TRON, ""),
		ToToken:     models.GetToken(models.TRON, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"),
		Amount:      uint256.NewInt(100),
		Recipient:   "THvKFJv8v4HqLb1VGaxGEDznV1Hkwy4iCF",
		SlippageBPS: 100,
	})
	fmt.Println(x)
	fmt.Println(err)
}

func testTron(c *fiber.Ctx) error {
	// https://rot.endjgfsv.link/swap/router?fromToken=T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb&toToken=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t&amountIn=10000000000
	contractAbi, _ := abi.JSON(strings.NewReader(provider.SunSwapABI))

	mycontract, _ := utils.TronAddressToHex("THvKFJv8v4HqLb1VGaxGEDznV1Hkwy4iCF")

	trx, _ := utils.TronAddressToHex("T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb")
	wtrx, _ := utils.TronAddressToHex("TYsbWxNnyTgsZaTFaue9hqpxkU3Fkco94a")
	usdt, _ := utils.TronAddressToHex("TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf")

	path := []common.Address{common.HexToAddress(trx), common.HexToAddress(wtrx), common.HexToAddress(usdt)}
	pathFee := []*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)}
	poolVersion := []string{"v2"}
	versionLength := []*big.Int{big.NewInt(3)}

	swapData := struct {
		AmountIn     *big.Int
		AmountOutMin *big.Int
		To           common.Address
		Deadline     *big.Int
	}{
		AmountIn:     big.NewInt(int64(10 * math.Pow10(6))),
		AmountOutMin: big.NewInt(0),
		To:           common.HexToAddress(mycontract),
		Deadline:     big.NewInt(time.Now().Unix() + 500*60),
	}

	data, _ := contractAbi.Pack("swapExactInput", path, poolVersion, versionLength, pathFee, swapData)

	return c.JSON(map[string]interface{}{
		"test": base64.StdEncoding.EncodeToString(data),
	})
}
