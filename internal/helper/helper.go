package helper

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/*
* Get Client IP
* @param c *gin.Context
* @return string
 */
func GetClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

/*
* Hash password using bcrypt
* @param password string
* @return string
 */
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

/*
* Check if password matches hash
* @param hashedPassword string
* @param password string
* @return bool
 */
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

/*
* Generate random token
* @param n int
* @return string
 */
func GenerateToken(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	sb := strings.Builder{}
	sb.Grow(n) // Pre-allocate memory for efficiency
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

/*
* Generate random user ID
* @return string
 */
func GenerateUserID() string {
	return strconv.Itoa(100000 + rand.Intn(900000)) // generates 6-digit number (100000–999999)
}

/*
* Generate 6-digit OTP
* @return int64, error
 */
func GenerateOTP() (int64, error) {
	max := big.NewInt(900000)              // 0 - 899999
	n, err := crand.Int(crand.Reader, max) //  two return values
	if err != nil {
		return 0, err
	}
	return n.Int64() + 100000, nil // 100000 - 999999 (always 6 digits)
}

/*
* Generate unique slug for model
* @param db *gorm.DB
* @param title string
* @param model interface{}
* @param field string
* @return string
 */
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

/*
* Mask email address for privacy
* @param email string
* @return string
 */
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	local := parts[0]
	domain := parts[1]

	if len(local) <= 3 {
		return email // Too short to mask safely based on rules
	}

	// Keep first 3 characters
	firstPart := local[:3]

	// Keep last 2 characters (if length permits)
	if len(local) > 5 {
		lastPart := local[len(local)-2:]
		return fmt.Sprintf("%s**%s@%s", firstPart, lastPart, domain)
	}

	// If length is between 5 and 6, just append **
	return fmt.Sprintf("%s**@%s", firstPart, domain)
}
