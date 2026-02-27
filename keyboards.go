package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetMainKeyboard() tgbotapi.InlineKeyboardMarkup {
    return tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("📅 Сегодня", "today"),
            tgbotapi.NewInlineKeyboardButtonData("📆 Неделя", "week"),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🔍 Поиск", "search"),
            tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
        ),
    )
}

func GetBackKeyboard() tgbotapi.InlineKeyboardMarkup {
    return tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("« Назад", "back"),
        ),
    )
}