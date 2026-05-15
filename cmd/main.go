package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Billy19191/Telegram-Morpho-Bot/model"
	"github.com/Billy19191/Telegram-Morpho-Bot/service"
	"github.com/Billy19191/Telegram-Morpho-Bot/util"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	// "github.com/joho/godotenv"
)

var morphoService *service.MorphoService

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	panic("ENV not found")
	// }

	tgBotToken := getEnvKey("TG_BOT_TOKEN", "")
	morphoService = service.NewMorphoService(
		getEnvKey("BASE_URL", ""),
		getEnvKey("WALLET_ADDRESS", ""),
		getEnvKey("CHAIN_ID", ""),
	)

	chatID, err := strconv.ParseInt(getEnvKey("TG_CHAT_ID", "0"), 10, 64)
	if err != nil || chatID == 0 {
		panic("TG_CHAT_ID is required and must be a valid number")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	tgBot, err := bot.New(tgBotToken, opts...)

	if err != nil {
		panic(err)
	}

	go startCronMonitor(ctx, tgBot, chatID)

	log.Println("🤖 Bot started. Cron monitor running every 5 minutes.")
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
	msg := formatVaultMessage(result.Data[0], riskReport)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}

// startCronMonitor checks vault positions every 5 minutes.
// - Critical status → sends alert immediately
// - Normal status → sends routine report every 8 hours
func startCronMonitor(ctx context.Context, b *bot.Bot, chatID int64) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	lastRoutineReport := time.Now()

	checkAndNotify(ctx, b, chatID, &lastRoutineReport)

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Cron monitor stopped.")
			return
		case <-ticker.C:
			checkAndNotify(ctx, b, chatID, &lastRoutineReport)
		}
	}
}

func checkAndNotify(ctx context.Context, b *bot.Bot, chatID int64, lastRoutineReport *time.Time) {
	result, err := morphoService.GetVaultPositions()
	if err != nil {
		log.Printf("❌ Cron check failed: %v", err)
		return
	}

	riskReport := service.EvaluateVaultRisk(result.Data[0])
	log.Printf("📋 Cron check — Status: %s", riskReport.OverallStatus)

	if riskReport.OverallStatus == model.StatusCritical {
		msg := "🚨 CRITICAL ALERT 🚨\n" + formatVaultMessage(result.Data[0], riskReport)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   msg,
		})
		*lastRoutineReport = time.Now()
		log.Println("🚨 Critical alert sent!")
		return
	}

	if time.Since(*lastRoutineReport) >= 8*time.Hour {
		msg := "📊 Routine Monitor Report\n\n" + formatVaultMessage(result.Data[0], riskReport)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   msg,
		})
		*lastRoutineReport = time.Now()
		log.Println("📊 Routine report sent!")
	}
}

func formatVaultMessage(vault model.VaultEntity, riskReport model.RiskReport) string {
	reasonLine := ""
	if riskReport.OverallStatus != model.StatusSafe {
		reasonLine = fmt.Sprintf("🔥 Reason: %s\n", riskReport.Reason)
	}

	return fmt.Sprintf(
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
		vault.VaultName,
		util.FormatNumberWithSeparator(vault.NetApy),
		util.FormatNumberWithSeparator(vault.MyAssetUsd),
		util.FormatNumberWithSeparator(vault.TotalAssetUsd),
		util.FormatNumberWithSeparator(vault.Liquidity),
		riskReport.Metrics[3].Value,
		util.FormatNumberWithSeparator(vault.SharedInVault),
	)
}
