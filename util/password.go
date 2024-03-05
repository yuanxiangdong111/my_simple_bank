package util

import (
    "fmt"

    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    // 产生哈希密码, 默认的cost是10
    haashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        // %w 用于包装错误
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(haashedPassword), nil
}

func checkPassword(password string, hashedPassword string) error {
    // 检查密码是否正确
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
