package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
)

func GenerateSessionToken(email string) string {
	data := fmt.Sprintf("%s:%s", constants.ForgotPasswordKey, email)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
