package other

import (
	"os"
)

// GetEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
