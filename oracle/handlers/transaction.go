package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/holiman/uint256"
	"gringotts/config"
	"gringotts/models"
	"gringotts/service"
	"gringotts/utils"
)

/*
	FRONT SENDS TOKENS AMOUNT WITH CORRECT DECIMALS!!!
*/

func HandleTransaction(c *fiber.Ctx) error {
	var request models.TransactionRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if len(request.SrcItems) != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	servingContext, err := models.NewServingContext(&request)
	if err != nil {
		log.Errorw("ServingContext", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	/* Inbound transaction */
	outAmountUSDX := uint64(0)

	for _, item := range servingContext.Inbounds {
		inErr := service.InboundTransaction(servingContext.SourceChain, item)

		if inErr != nil {
			log.Errorw("Inbound error", "err", inErr)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid swap request - in"})
		}

		if item.Swap == nil {
			outAmountUSDX = outAmountUSDX + utils.MoveDecimals(item.Amount, item.Token.Decimals, config.ChainTransferDecimals).Uint64()
		} else {
			outAmountUSDX = outAmountUSDX + utils.MoveDecimals(item.Swap.OutAmount, item.Swap.ToToken.Decimals, config.ChainTransferDecimals).Uint64()
		}
	}

	/* Estimate Calculation */
	//err = service.EstimateMarketplace(servingContext, uint256.NewInt(outAmountUSDX))
	//if err != nil {
	//	log.Errorw("Estimate error", "err", err)
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	//}
	//
	//outAmountUSDX = outAmountUSDX - servingContext.Marketplace.GasPriceUSDX.Uint64() - servingContext.Marketplace.CommissionUSDX.Uint64()
	//outAmountUSDX = outAmountUSDX + servingContext.Marketplace.GasPriceDiscountUSDX.Uint64() + servingContext.Marketplace.CommissionDiscountUSDX.Uint64()
	//
	//outAmountUSDX = uint64(float64(outAmountUSDX) * config.ConversionFactor)

	/* Outbound transaction */
	for chain, items := range servingContext.Outbounds {
		for _, item := range items {
			item.Amount = utils.ApplyBPS(uint256.NewInt(outAmountUSDX), item.DistributionBPS)

			outErr := service.OutboundTransaction(chain, item)

			if outErr != nil {
				log.Errorw("Outbound error", "err", outErr)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid outTransaction request - out"})
			}
		}
	}

	err = service.CreateBlockchainTransaction(uint256.NewInt(outAmountUSDX).ToBig(), servingContext)
	if err != nil {
		log.Errorw("create final transaction error", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	return c.JSON(BuildResponse(servingContext))
}

func BuildResponse(servingContext *models.ServingContext) *models.TransactionResponse {
	inputTransactions := make([]*models.TransactionItemResponse, 0)
	for _, item := range servingContext.Inbounds {
		transaction := &models.TransactionItemResponse{
			Token: item.Token.Address,
		}

		if item.Swap == nil {
			transaction.OutAmount = item.Amount
			transaction.MinOutAmount = item.Amount
			transaction.USDXToken = utils.MoveDecimals(item.Amount, item.Token.Decimals, config.ChainTransferDecimals)
		} else {
			transaction.OutAmount = item.Swap.OutAmount
			transaction.MinOutAmount = item.Swap.MinOutAmount
			transaction.USDXToken = utils.MoveDecimals(item.Swap.OutAmount, item.Swap.ToToken.Decimals, config.ChainTransferDecimals)
		}
		inputTransactions = append(inputTransactions, transaction)
	}

	outputTransactions := make(map[models.Blockchain][]*models.TransactionItemResponse)
	for chain, items := range servingContext.Outbounds {
		outputTransactions[chain] = make([]*models.TransactionItemResponse, 0)
		for _, item := range items {
			transaction := &models.TransactionItemResponse{
				Token: item.Token.Address,
			}

			if item.Swap == nil {
				transaction.OutAmount = item.Amount
				transaction.MinOutAmount = item.Amount
				transaction.USDXToken = utils.MoveDecimals(item.Amount, item.Token.Decimals, config.ChainTransferDecimals)
			} else {
				transaction.OutAmount = item.Swap.OutAmount
				transaction.MinOutAmount = item.Swap.MinOutAmount
				transaction.USDXToken = item.Amount
			}

			outputTransactions[chain] = append(outputTransactions[chain], transaction)
		}
	}

	return &models.TransactionResponse{
		InTransaction:   map[models.Blockchain][]*models.TransactionItemResponse{servingContext.SourceChain: inputTransactions},
		OutTransactions: outputTransactions,
		Marketplace: &models.MarketplaceResponse{
			CommissionUSDX:         servingContext.Marketplace.CommissionUSDX.Uint64(),
			CommissionDiscountUSDX: servingContext.Marketplace.CommissionDiscountUSDX.Uint64(),
			GasPriceUSDX:           servingContext.Marketplace.GasPriceUSDX.Uint64(),
			GasPriceDiscountUSDX:   servingContext.Marketplace.GasPriceDiscountUSDX.Uint64(),
		},
		Transaction: &models.UnsignedTransactionResponse{
			Contract: servingContext.SourceChain.GetContract(),
			Data:     utils.ToHex(servingContext.Tx.Data),
			Value:    servingContext.Tx.Value,
		},
	}
}
