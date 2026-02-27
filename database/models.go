package database

import "time"

type User struct {
    ID              int64     `db:"id"`
    TelegramID      int64     `db:"telegram_id"`
    GroupName       string    `db:"group_name"`
    Notifications   bool      `db:"notifications"`
    NotificationTime string   `db:"notification_time"`
    CreatedAt       time.Time `db:"created_at"`
}

type Lesson struct {
    ID          int    `db:"id"`
    GroupName   string `db:"group_name"`
    DayOfWeek   int    `db:"day_of_week"`
    LessonNumber int   `db:"lesson_number"`
    LessonName  string `db:"lesson_name"`
    Teacher     string `db:"teacher"`
    Classroom   string `db:"classroom"`
    LessonType  string `db:"lesson_type"` // lecture, practice, lab
    StartTime   string `db:"start_time"`
    EndTime     string `db:"end_time"`
    WeekType    string `db:"week_type"` // numerator, denominator, both
}