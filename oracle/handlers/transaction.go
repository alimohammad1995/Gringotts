package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/holiman/uint256"
	"gringotts/config"
	"gringotts/models"
	"gringotts/provider"
	"gringotts/service"
	"gringotts/utils"
	"math"
)

type OutboundTransactionItemRequest struct {
	Token           string `json:"token"`
	DistributionBPS int    `json:"distribution"`
	Recipient       string `json:"recipient"`
}

type InboundTransactionItemRequest struct {
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

type TransactionRequest struct {
	User        string                                                 `json:"user"`
	SrcItems    map[models.Blockchain][]InboundTransactionItemRequest  `json:"src_items"`
	DstItems    map[models.Blockchain][]OutboundTransactionItemRequest `json:"dst_items"`
	SlippageBPS int                                                    `json:"slippage_bps"`
}

type TransactionResponse struct {
	Marketplace     *models.Marketplace                         `json:"marketplace"`
	InTransaction   map[models.Blockchain][]*models.Transaction `json:"in_transactions"`
	OutTransactions map[models.Blockchain][]*models.Transaction `json:"out_transactions"`
	Transaction     *models.UnsignedTransaction                 `json:"transaction"`
}

/*
	FRONT SENDS TOKENS AMOUNT WITH CORRECT DECIMALS!!!
*/

func HandleTransaction(c *fiber.Ctx) error {
	var request TransactionRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if len(request.SrcItems) != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	request.SlippageBPS = utils.Max(request.SlippageBPS, config.DefaultSlippageBPS)

	/* Inbound transaction */
	outAmountUSDX := uint64(0)

	var sourceChain models.Blockchain
	chainInTransactions := make(map[models.Blockchain][]*models.Transaction)

	for chain, items := range request.SrcItems {
		sourceChain = chain
		inTransactions := make([]*models.Transaction, 0)

		for _, tx := range items {
			srcAmount, err := uint256.FromDecimal(tx.Amount)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
			}

			inTransaction, inErr := service.InboundTransaction(
				chain,
				tx.Token,
				srcAmount,
				request.SlippageBPS,
			)
			if inErr != nil {
				log.Errorw("Inbound error", "err", err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid swap request - in"})
			}

			outAmountUSDX = outAmountUSDX + utils.MoveDecimals(inTransaction.OutAmount, inTransaction.ToToken.Decimals, config.ChainTransferDecimals).Uint64()

			inTransaction.SrcAmount = srcAmount
			inTransactions = append(inTransactions, transform(inTransaction))
		}

		chainInTransactions[chain] = inTransactions
	}

	/* Estimate Calculation */
	dstMapping := make(map[models.Blockchain][]string)
	for chain, transactions := range request.DstItems {
		dstMapping[chain] = make([]string, 0)
		for _, transaction := range transactions {
			dstMapping[chain] = append(dstMapping[chain], transaction.Token)
		}
	}
	marketplace, err := service.EstimateMarketplace(sourceChain, dstMapping, uint256.NewInt(outAmountUSDX))
	if err != nil {
		log.Errorw("Estimate error", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	outAmountUSDX = outAmountUSDX - marketplace.GasPriceUSD.Uint64()
	outAmountUSDX = outAmountUSDX - marketplace.CommissionUSD.Uint64()
	outAmountUSDX = uint64(float64(outAmountUSDX) * config.ConversionFactor)

	marketplaceCommission := marketplace.CommissionUSD.Float64() / math.Pow10(config.ChainTransferDecimals)
	marketplaceGas := marketplace.GasPriceUSD.Float64() / math.Pow10(config.ChainTransferDecimals)

	/* Outbound transaction */
	chainOutTransactions := make(map[models.Blockchain][]*models.Transaction)
	for chain, transactions := range request.DstItems {
		outTransactions := make([]*models.Transaction, 0)

		for _, transaction := range transactions {
			outTransaction, outErr := service.OutboundTransaction(
				chain,
				transaction.Recipient,
				transaction.Token,
				utils.ApplyBPS(uint256.NewInt(outAmountUSDX), transaction.DistributionBPS),
				request.SlippageBPS,
			)

			if outErr != nil {
				log.Errorw("Outbound error", "err", outErr)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid outTransaction request - out"})
			}

			outTransaction.DistributionBPS = transaction.DistributionBPS
			outTransactions = append(outTransactions, transform(outTransaction))
		}

		chainOutTransactions[chain] = outTransactions
	}

	tx, err := service.CreateBlockchainTransaction(
		uint256.NewInt(outAmountUSDX).ToBig(),
		request.User,
		sourceChain,
		chainInTransactions,
		chainOutTransactions,
	)
	if err != nil {
		log.Errorw("create final transaction error", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	return c.JSON(TransactionResponse{
		InTransaction:   chainInTransactions,
		OutTransactions: chainOutTransactions,
		Marketplace: &models.Marketplace{
			Commission:  marketplaceCommission,
			GasPriceUSD: marketplaceGas,
		},
		Transaction: &models.UnsignedTransaction{
			Contract: tx.Contract,
			Data:     tx.Data,
			Value:    tx.Value,
		},
	})
}

func transform(transaction *provider.Transaction) *models.Transaction {
	out := &models.Transaction{
		FromToken:       transaction.FromToken.Address,
		ToToken:         transaction.ToToken.Address,
		OutAmount:       transaction.OutAmount,
		MinOutAmount:    transaction.MinOutAmount,
		Recipient:       transaction.Recipient,
		DistributionBPS: transaction.DistributionBPS,
	}

	if transaction.SrcAmount != nil {
		out.SrcAmount = transaction.SrcAmount.String()
	}

	if transaction.Swap != nil {
		out.Swap = &models.Swap{
			Address:       transaction.Swap.ExecutorAddress,
			Command:       transaction.Swap.Command,
			MetaData:      transaction.Swap.Metadata,
			AddressLookup: transaction.Swap.AddressLookup,
		}
	}

	return out
}
