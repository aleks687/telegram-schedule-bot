package bot

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "telegram-schedule-bot/database"
)

type Lesson struct {
    DayOfWeek    string
    Date         string
    StartTime    string
    EndTime      string
    LessonType   string
    LessonName   string
    Classroom    string  // Добавили аудиторию
    Teacher      string  // Добавили преподавателя
}

type ScheduleProvider interface {
    GetTodaySchedule(groupID string) ([]Lesson, error)
    GetWeekSchedule(groupID string) (map[string][]Lesson, error)
}

type Bot struct {
    API        *tgbotapi.BotAPI
    DB         *database.Database
    Parser     ScheduleProvider
    GroupID    string
}

func NewBot(token string, db *database.Database, parser ScheduleProvider, groupID string) (*Bot, error) {
    api, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, err
    }
    
    api.Debug = true
    
    return &Bot{
        API:     api,
        DB:      db,
        Parser:  parser,
        GroupID: groupID,
    }, nil
}

func (b *Bot) HandleUpdates() {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    
    updates := b.API.GetUpdatesChan(u)
    
    for update := range updates {
        if update.Message != nil {
            b.handleMessage(update.Message)
        }
    }
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
    switch message.Text {
    case "/start":
        b.sendMainKeyboard(message.Chat.ID)
    case "📅 Сегодня":
        b.sendTodaySchedule(message.Chat.ID)
    case "📆 Неделя":
        b.sendWeekSchedule(message.Chat.ID)
    case "❓ Помощь":
        b.sendHelp(message.Chat.ID)
    }
}

func (b *Bot) sendMainKeyboard(chatID int64) {
    text := "👋 Привет! Я бот расписания ТГПИ.\n\nВыбери действие:"
    
    keyboard := tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("📅 Сегодня"),
            tgbotapi.NewKeyboardButton("📆 Неделя"),
        ),
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("❓ Помощь"),
        ),
    )
    
    msg := tgbotapi.NewMessage(chatID, text)
    msg.ReplyMarkup = keyboard
    b.API.Send(msg)
}

func (b *Bot) sendTodaySchedule(chatID int64) {
    b.API.Send(tgbotapi.NewMessage(chatID, "⏳ Загружаю..."))
    
    lessons, err := b.Parser.GetTodaySchedule(b.GroupID)
    if err != nil {
        b.API.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки"))
        return
    }
    
    if len(lessons) == 0 {
        b.API.Send(tgbotapi.NewMessage(chatID, "📅 Сегодня пар нет"))
        return
    }
    
    text := "📅 СЕГОДНЯ:\n\n"
    for i, l := range lessons {
        text += string(rune(49+i)) + ". " + cleanText(l.LessonName) + "\n"
        text += "   " + cleanText(l.LessonType) + "\n"
        text += "   🕒 " + cleanText(l.StartTime) + "-" + cleanText(l.EndTime) + "\n"
        if l.Classroom != "" {
            text += "   📍 " + cleanText(l.Classroom) + "\n"
        }
        if l.Teacher != "" {
            text += "   👨‍🏫 " + cleanText(l.Teacher) + "\n"
        }
        text += "\n"
    }
    
    msg := tgbotapi.NewMessage(chatID, text)
    b.API.Send(msg)
}

func (b *Bot) sendWeekSchedule(chatID int64) {
    b.API.Send(tgbotapi.NewMessage(chatID, "⏳ Загружаю..."))
    
    week, err := b.Parser.GetWeekSchedule(b.GroupID)
    if err != nil {
        b.API.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки"))
        return
    }
    
    if len(week) == 0 {
        b.API.Send(tgbotapi.NewMessage(chatID, "📆 На неделе пар нет"))
        return
    }
    
    text := "📆 НЕДЕЛЯ:\n\n"
    days := []string{"понедельник", "вторник", "среда", "четверг", "пятница", "суббота"}
    
    for _, day := range days {
        if lessons, ok := week[day]; ok && len(lessons) > 0 {
            text += "🔹 " + cleanText(day) + ":\n"
            for _, l := range lessons {
                text += "   • " + cleanText(l.LessonName) + "\n"
                text += "     " + cleanText(l.LessonType) + " 🕒 " + cleanText(l.StartTime) + "-" + cleanText(l.EndTime) + "\n"
                if l.Classroom != "" {
                    text += "     📍 " + cleanText(l.Classroom) + "\n"
                }
                if l.Teacher != "" {
                    text += "     👨‍🏫 " + cleanText(l.Teacher) + "\n"
                }
            }
            text += "\n"
        }
    }
    
    msg := tgbotapi.NewMessage(chatID, text)
    b.API.Send(msg)
}

func (b *Bot) sendHelp(chatID int64) {
    text := "❓ Помощь\n\n📅 Сегодня\n📆 Неделя\n\nГруппа: 13499"
    b.API.Send(tgbotapi.NewMessage(chatID, text))
}

func cleanText(s string) string {
    result := make([]rune, 0, len(s))
    for _, r := range s {
        if r >= 32 && r != 65533 {
            result = append(result, r)
        }
    }
    return string(result)
}