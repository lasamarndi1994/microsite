package helper

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gosimple/slug"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateToken(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	sb := strings.Builder{}
	sb.Grow(n) // Pre-allocate memory for efficiency
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func GenerateUserID() string {
	return strconv.Itoa(100000 + rand.Intn(900000)) // generates 6-digit number (100000–999999)
}

func GenerateOTP() (int64, error) {
	max := big.NewInt(1000000)             // 0 - 999999
	n, err := crand.Int(crand.Reader, max) //  two return values
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil // always 6 digits
}

func GenerateUniqueSlug(db *gorm.DB, title string, model interface{}, field string) string {
	baseSlug := slug.Make(title)
	uniqueSlug := baseSlug
	counter := 1

	for {
		var count int64
		db.Model(model).Where(fmt.Sprintf("%s = ?", field), uniqueSlug).Count(&count)

		if count == 0 {
			break
		}

		uniqueSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return uniqueSlug
}
