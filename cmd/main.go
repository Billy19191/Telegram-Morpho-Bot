package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Billy19191/Telegram-Morpho-Bot/model"
	"github.com/Billy19191/Telegram-Morpho-Bot/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var morphoService *service.MorphoService

func main() {
	if err := godotenv.Load(); err != nil {
		panic("ENV not found")
	}

	tgBotToken := getEnvKey("TG_BOT_TOKEN", "")
	morphoService = service.NewMorphoService(
		getEnvKey("BASE_URL", ""),
		getEnvKey("WALLET_ADDRESS", ""),
		getEnvKey("CHAIN_ID", ""),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	tgBot, err := bot.New(tgBotToken, opts...)

	if err != nil {
		panic(err)
	}

	tgBot.Start(ctx)
}

func getEnvKey(key string, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallbackValue
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("User message: " + update.Message.Text)

	result, err := morphoService.GetVaultPositions()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Error: %s", err.Error()),
		})
		return
	}

	riskReport := service.EvaluateVaultRisk(result.Data[0])

	reasonLine := ""
	if riskReport.OverallStatus != model.StatusSafe {
		reasonLine = fmt.Sprintf("🔥 Reason: %s\n", riskReport.Reason)
	}

	messageTemplate := fmt.Sprintf(
		"⚡ Overall Status: %s\n"+
			"%s"+
			"----------------------\n"+
			"🏦 Vault Name: %s\n"+
			"📈 Current APY: %s%%\n"+
			"💰 My Asset USD: $%s\n"+
			"----------------------\n"+
			"🏛️ Total Asset USD: $%s\n"+
			"💧 Liquidity: %s\n"+
			"📊 Utilization: %s%%\n"+
			"🤝 Shared In Vault: %s%%\n",
		riskReport.OverallStatus,
		reasonLine,
		result.Data[0].VaultName,
		formatNumberWithSeparator(result.Data[0].NetApy),
		formatNumberWithSeparator(result.Data[0].MyAssetUsd),
		formatNumberWithSeparator(result.Data[0].TotalAssetUsd),
		formatNumberWithSeparator(result.Data[0].Liquidity),
		riskReport.Metrics[3].Value,
		formatNumberWithSeparator(result.Data[0].SharedInVault),
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messageTemplate,
	})
}

func formatNumberWithSeparator(number float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", number)
}
