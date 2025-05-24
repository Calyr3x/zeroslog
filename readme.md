# zeroslog

**zeroslog** — быстрый handler для стандартного Go-пакета `log/slog`. Ориентирован на минимальные задержки и малое количество аллокаций при интенсивном логировании.

---

## Возможности

- Быстрый форматтер, минимум аллокаций, синхронизация только на вывод.
- Настраиваемый формат времени, уровень логирования, цветной вывод, любой io.Writer.
- Поддержка структурированных логов: key-value поля, группы, вложенные атрибуты.
- Полная совместимость с API Go `log/slog`.

---

## Быстрый старт

```bash
go get github.com/calyrexx/zeroslog
```

```go
import (
	"log/slog"
	"os"
	"github.com/calyrexx/zeroslog"
)

func main() {
	logger := slog.New(zeroslog.New(
		zeroslog.WithOutput(os.Stdout),
		zeroslog.WithColors(),
		zeroslog.WithMinLevel(slog.LevelInfo),
		zeroslog.WithTimeFormat("2006-01-02 15:04:05.000 -07:00"),
	))

	logger.Info("Service started",
		"port", 8080,
		"mode", "prod",
	)
}
```

### Опции

| Опция                       | Описание                         | Значение по умолчанию      |
|-----------------------------|----------------------------------|----------------------------|
| `WithOutput(io.Writer)`     | Вывод логов                      | `os.Stderr`                |
| `WithColors()`              | Цветной режим вывода             | Отключено                  |
| `WithMinLevel(slog.Level)`  | Минимальный уровень логирования  | `slog.LevelInfo`           |
| `WithTimeFormat(string)`    | Формат времени                   | `RFC3339`                  |

---

## Производительность

Бенчмарки:

| Handler          | ns/op | B/op | allocs/op |
| ---------------- | ----- | ---- | --------- |
| slog.TextHandler | 1220  | 32   | 3         |
| **zeroslog**     | 934   | 64   | 6         |
| logrus           | 3290  | 1410 | 30        |

*zeroslog минимум в 3 раза быстрее logrus, и сопоставим по скорости с slog.TextHandler при большем количестве поддерживаемых фич.*

---