package scheduler

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/chromedp/chromedp"
    "telegram-schedule-bot/bot"
)

type BrowserParser struct {
    BaseURL string
}

func NewBrowserParser(baseURL string) *BrowserParser {
    return &BrowserParser{BaseURL: baseURL}
}

func (p *BrowserParser) GetTodaySchedule(groupID string) ([]bot.Lesson, error) {
    url := fmt.Sprintf("%s/schedule/group/%s/", p.BaseURL, groupID)
    
    fmt.Println("==========================================")
    fmt.Println("🚀 ЗАПУСК БРАУЗЕРА ДЛЯ СЕГОДНЯ:", url)
    fmt.Println("==========================================")
    
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", false),
        chromedp.Flag("disable-gpu", true),
    )
    
    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()
    
    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()
    
    ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
    defer cancel()
    
    now := time.Now()
    todayDate := now.Format("02.01.2006")
    todayName := getCurrentDayRussian()
    
    fmt.Printf("📅 Ищем: %s, %s\n", todayName, todayDate)
    
    var rawLessons []struct {
        Day     string `json:"day"`
        Date    string `json:"date"`
        Start   string `json:"start"`
        End     string `json:"end"`
        Subject string `json:"subject"`
        Type    string `json:"type"`
        Room    string `json:"room"`
        Teacher string `json:"teacher"`
    }
    
    jsCode := `
        (function() {
            let result = [];
            let days = document.querySelectorAll('.day');
            
            for(let day of days) {
                let dayHead = day.querySelector('.head');
                let dayText = dayHead ? dayHead.textContent : '';
                
                let dayMatch = dayText.match(/([а-я]+),\s+(\d{2}\.\d{2}\.\d{4})/);
                if(!dayMatch) continue;
                
                let dayName = dayMatch[1];
                let dayDate = dayMatch[2];
                
                let rows = day.querySelectorAll('tr');
                
                for(let row of rows) {
                    let timeCell = row.querySelector('.time');
                    if(!timeCell) continue;
                    
                    let timeText = timeCell.textContent.trim();
                    let timeMatch = timeText.match(/(\d{2}:\d{2})\s+(\d{2}:\d{2})/);
                    if(!timeMatch) continue;
                    
                    let startTime = timeMatch[1];
                    let endTime = timeMatch[2];
                    
                    let subjectSpan = row.querySelector('.sbtype');
                    if(!subjectSpan) continue;
                    
                    let typeText = subjectSpan.textContent.trim();
                    let nextText = subjectSpan.nextSibling ? subjectSpan.nextSibling.textContent : '';
                    
                    // Ищем аудиторию
                    let roomLink = row.querySelector('a[href*="/schedule/aud/"]');
                    let roomText = roomLink ? roomLink.textContent.trim() : '';
                    
                    // Ищем преподавателя
                    let teacherLink = row.querySelector('a[href*="/schedule/teacher/"]');
                    let teacherText = teacherLink ? teacherLink.textContent.trim() : '';
                    
                    result.push({
                        day: dayName,
                        date: dayDate,
                        start: startTime,
                        end: endTime,
                        subject: nextText,
                        type: typeText,
                        room: roomText,
                        teacher: teacherText
                    });
                }
            }
            return result;
        })();
    `
    
    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.Sleep(3*time.Second),
        chromedp.WaitVisible(`table.typelinks`, chromedp.ByQuery),
        chromedp.Sleep(2*time.Second),
        chromedp.Evaluate(jsCode, &rawLessons),
    )
    
    if err != nil {
        fmt.Println("❌ ОШИБКА БРАУЗЕРА:", err)
        return nil, err
    }
    
    fmt.Printf("✅ ВСЕГО НАЙДЕНО ЗАНЯТИЙ: %d\n", len(rawLessons))
    
    var result []bot.Lesson
    for _, l := range rawLessons {
        if strings.EqualFold(l.Day, todayName) && l.Date == todayDate {
            lessonType := parseLessonType(l.Type)
            lessonName := cleanSubject(l.Subject)
            
            result = append(result, bot.Lesson{
                DayOfWeek:  l.Day,
                Date:       l.Date,
                StartTime:  l.Start,
                EndTime:    l.End,
                LessonType: lessonType,
                LessonName: lessonName,
                Classroom:  cleanText(l.Room),
                Teacher:    cleanTeacherName(l.Teacher),
            })
            fmt.Printf("   ✅ %s %s %s [%s] %s\n", l.Start, lessonType, lessonName, l.Room, l.Teacher)
        }
    }
    
    sortLessonsByTime(result)
    fmt.Printf("📊 ИТОГО НА СЕГОДНЯ: %d пар\n", len(result))
    
    return result, nil
}

func (p *BrowserParser) GetWeekSchedule(groupID string) (map[string][]bot.Lesson, error) {
    url := fmt.Sprintf("%s/schedule/group/%s/", p.BaseURL, groupID)
    
    fmt.Println("==========================================")
    fmt.Println("🚀 ЗАПУСК БРАУЗЕРА ДЛЯ НЕДЕЛИ:", url)
    fmt.Println("==========================================")
    
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", false),
        chromedp.Flag("disable-gpu", true),
    )
    
    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()
    
    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()
    
    ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
    defer cancel()
    
    var rawLessons []struct {
        Day     string `json:"day"`
        Date    string `json:"date"`
        Start   string `json:"start"`
        End     string `json:"end"`
        Subject string `json:"subject"`
        Type    string `json:"type"`
        Room    string `json:"room"`
        Teacher string `json:"teacher"`
    }
    
    jsCode := `
        (function() {
            let result = [];
            let days = document.querySelectorAll('.day');
            
            for(let day of days) {
                let dayHead = day.querySelector('.head');
                let dayText = dayHead ? dayHead.textContent : '';
                
                let dayMatch = dayText.match(/([а-я]+),\s+(\d{2}\.\d{2}\.\d{4})/);
                if(!dayMatch) continue;
                
                let dayName = dayMatch[1];
                let dayDate = dayMatch[2];
                
                let rows = day.querySelectorAll('tr');
                
                for(let row of rows) {
                    let timeCell = row.querySelector('.time');
                    if(!timeCell) continue;
                    
                    let timeText = timeCell.textContent.trim();
                    let timeMatch = timeText.match(/(\d{2}:\d{2})\s+(\d{2}:\d{2})/);
                    if(!timeMatch) continue;
                    
                    let startTime = timeMatch[1];
                    let endTime = timeMatch[2];
                    
                    let subjectSpan = row.querySelector('.sbtype');
                    if(!subjectSpan) continue;
                    
                    let typeText = subjectSpan.textContent.trim();
                    let nextText = subjectSpan.nextSibling ? subjectSpan.nextSibling.textContent : '';
                    
                    let roomLink = row.querySelector('a[href*="/schedule/aud/"]');
                    let roomText = roomLink ? roomLink.textContent.trim() : '';
                    
                    let teacherLink = row.querySelector('a[href*="/schedule/teacher/"]');
                    let teacherText = teacherLink ? teacherLink.textContent.trim() : '';
                    
                    result.push({
                        day: dayName,
                        date: dayDate,
                        start: startTime,
                        end: endTime,
                        subject: nextText,
                        type: typeText,
                        room: roomText,
                        teacher: teacherText
                    });
                }
            }
            return result;
        })();
    `
    
    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.Sleep(3*time.Second),
        chromedp.WaitVisible(`table.typelinks`, chromedp.ByQuery),
        chromedp.Sleep(2*time.Second),
        chromedp.Evaluate(jsCode, &rawLessons),
    )
    
    if err != nil {
        fmt.Println("❌ ОШИБКА БРАУЗЕРА:", err)
        return nil, err
    }
    
    fmt.Printf("✅ ВСЕГО НАЙДЕНО ЗАНЯТИЙ: %d\n", len(rawLessons))
    
    now := time.Now()
    week := make(map[string][]bot.Lesson)
    
    for _, l := range rawLessons {
        dateParts := strings.Split(l.Date, ".")
        if len(dateParts) != 3 {
            continue
        }
        
        day := parseInt(dateParts[0])
        month := parseInt(dateParts[1])
        year := parseInt(dateParts[2])
        
        lessonDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
        diff := int(lessonDate.Sub(now).Hours() / 24)
        
        if diff >= 0 && diff <= 7 {
            dayName := strings.ToLower(l.Day)
            lessonType := parseLessonType(l.Type)
            lessonName := cleanSubject(l.Subject)
            
            week[dayName] = append(week[dayName], bot.Lesson{
                DayOfWeek:  l.Day,
                Date:       l.Date,
                StartTime:  l.Start,
                EndTime:    l.End,
                LessonType: lessonType,
                LessonName: lessonName,
                Classroom:  cleanText(l.Room),
                Teacher:    cleanTeacherName(l.Teacher),
            })
            fmt.Printf("   ✅ %s %s %s %s [%s] %s\n", l.Date, l.Day, l.Start, lessonName, l.Room, l.Teacher)
        }
    }
    
    for day := range week {
        sortLessonsByTime(week[day])
    }
    
    return week, nil
}

func parseInt(s string) int {
    var result int
    fmt.Sscanf(s, "%d", &result)
    return result
}

func parseLessonType(typeStr string) string {
    typeStr = strings.TrimSpace(typeStr)
    
    switch typeStr {
    case "лек":
        return "📚 Лекция"
    case "пр":
        return "✍️ Практика"
    case "лаб":
        return "🔬 Лабораторная"
    default:
        return "📝 Занятие"
    }
}

func cleanSubject(subject string) string {
    subject = strings.TrimSpace(subject)
    subject = strings.Join(strings.Fields(subject), " ")
    return subject
}

func cleanTeacherName(teacher string) string {
    teacher = strings.TrimSpace(teacher)
    teacher = strings.ReplaceAll(teacher, "доц.", "👨‍🏫")
    teacher = strings.ReplaceAll(teacher, "ст. преп.", "👨‍🏫")
    teacher = strings.ReplaceAll(teacher, "гпх спец.", "👨‍🔬")
    return teacher
}

func sortLessonsByTime(lessons []bot.Lesson) {
    for i := 0; i < len(lessons)-1; i++ {
        for j := i + 1; j < len(lessons); j++ {
            if lessons[i].StartTime > lessons[j].StartTime {
                lessons[i], lessons[j] = lessons[j], lessons[i]
            }
        }
    }
}

func getCurrentDayRussian() string {
    weekdays := map[time.Weekday]string{
        time.Monday:    "понедельник",
        time.Tuesday:   "вторник",
        time.Wednesday: "среда",
        time.Thursday:  "четверг",
        time.Friday:    "пятница",
        time.Saturday:  "суббота",
        time.Sunday:    "воскресенье",
    }
    return weekdays[time.Now().Weekday()]
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