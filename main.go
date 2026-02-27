package main

import (
    "log"
    "telegram-schedule-bot/bot"
    "telegram-schedule-bot/config"
    "telegram-schedule-bot/database"
    "telegram-schedule-bot/scheduler"
)

func main() {
    cfg := config.LoadConfig()
    
    if cfg.BotToken == "" {
        log.Fatal("BOT_TOKEN не установлен!")
    }
    
    db, err := database.NewDatabase(cfg.DatabasePath)
    if err != nil {
        log.Fatal("Ошибка БД:", err)
    }
    defer db.Close()
    
    groupID := "13499"
    
    parser := scheduler.NewBrowserParser(cfg.ScheduleURL)
    
    b, err := bot.NewBot(cfg.BotToken, db, parser, groupID)
    if err != nil {
        log.Fatal("Ошибка создания бота:", err)
    }
    
    log.Println("✅ Бот запущен")
    b.HandleUpdates()
}