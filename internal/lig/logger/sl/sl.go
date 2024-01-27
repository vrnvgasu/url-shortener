package sl

import "log/slog"

// своя реализация добавления параметра "Ошибка" в лог
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
