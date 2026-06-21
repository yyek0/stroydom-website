package logger

import (
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 1. Создаем папку logs, если её еще нет (os.ModePerm дает нужные права)
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("Не удалось создать директорию для логов: %v", err)
	}

	// 2. Открываем файлы уже внутри новой папки
	infoFile, err := os.OpenFile(filepath.Join(logDir, "info.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть info.log: %v", err)
	}

	errFile, err := os.OpenFile(filepath.Join(logDir, "error.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть error.log: %v", err)
	}

	// 3. Создаем фильтры (какие уровни куда пускать)
	// В инфо-лог пускаем всё, что НИЖЕ предупреждения (Debug, Info)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	// В лог ошибок пускаем всё, что РАВНО или ВЫШЕ предупреждения (Warn, Error, Fatal)
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 4. Собираем "ядра" (cores)
	coreInfo := zapcore.NewCore(jsonEncoder, zapcore.AddSync(infoFile), infoLevel)
	coreError := zapcore.NewCore(jsonEncoder, zapcore.AddSync(errFile), errorLevel)

	// Опционально: оставляем вывод в консоль для удобства разработки
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	coreConsole := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)

	// 5. Соединяем всё в один тройник
	core := zapcore.NewTee(coreInfo, coreError, coreConsole)

	return zap.New(core)
}
