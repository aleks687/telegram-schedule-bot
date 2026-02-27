package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    BotToken     string
    DatabasePath string
    ScheduleURL  string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("Файл .env не найден, используем переменные окружения")
    }

    return &Config{
        BotToken:     getEnv("BOT_TOKEN", ""),
        DatabasePath: getEnv("DB_PATH", "schedule.db"),
        ScheduleURL:  getEnv("SCHEDULE_URL", ""),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}