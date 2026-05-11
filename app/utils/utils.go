package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/btcsuite/btcd/btcutil/base58"
	nid "github.com/matoous/go-nanoid/v2"
)

// IsExist 判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {

		return true
	}

	if os.IsExist(err) {

		return true
	}

	return false
}

func EpusdtSign(data map[string]interface{}, token string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		if k == "signature" {

			continue
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)
	var sign strings.Builder
	for _, k := range keys {
		v := data[k]
		if v == nil || v == "" {

			continue
		}

		sign.WriteString(k)
		sign.WriteString("=")
		sign.WriteString(fmt.Sprintf("%v", v))
		sign.WriteString("&")
	}

	signString := strings.TrimRight(sign.String(), "&")

	return Md5String(signString + token)
}

func GenerateTradeId() (string, error) {
	var defaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	return nid.Generate(defaultAlphabet, 18)
}

func Md5String(text string) string {

	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}

func Ec(str string) string {
	escapeChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	for _, char := range escapeChars {
		str = strings.ReplaceAll(str, char, "\\"+char)
	}

	return str
}

func IsNumber(s string) bool {
	match, err := regexp.MatchString(`^\d+\.?\d*$`, s)

	return match && err == nil
}

func IsValidTronAddress(address string) bool {
	match, err := regexp.MatchString(`^T[a-zA-Z0-9]{33}$`, address)

	return match && err == nil
}

// IsValidTonAddress 校验 TON 地址是否为 raw 或 user-friendly 格式。
func IsValidTonAddress(address string) bool {
	_, ok := NormalizeTonAddress(address)

	return ok
}

// NormalizeTonAddress 将 TON raw/user-friendly 地址统一转换为 raw 地址。
func NormalizeTonAddress(address string) (string, bool) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", false
	}

	rawMatch, _ := regexp.MatchString(`^-?\d+:[0-9a-fA-F]{64}$`, address)
	if rawMatch {
		parts := strings.SplitN(address, ":", 2)

		return parts[0] + ":" + strings.ToLower(parts[1]), true
	}

	data, err := base64.RawURLEncoding.DecodeString(address)
	if err != nil {
		data, err = base64.RawStdEncoding.DecodeString(address)
	}
	if err != nil || len(data) != 36 {
		return "", false
	}

	if tonCRC16(data[:34]) != binary.BigEndian.Uint16(data[34:]) {
		return "", false
	}

	workchain := int(int8(data[1]))

	return fmt.Sprintf("%d:%s", workchain, hex.EncodeToString(data[2:34])), true
}

// FormatTonAddress 将 TON 地址转换为 non-bounceable user-friendly 格式用于展示和付款。
func FormatTonAddress(address string) (string, bool) {
	raw, ok := NormalizeTonAddress(address)
	if !ok {
		return "", false
	}

	parts := strings.SplitN(raw, ":", 2)
	workchain, err := strconv.Atoi(parts[0])
	if err != nil || workchain < -128 || workchain > 127 {
		return "", false
	}

	hash, err := hex.DecodeString(parts[1])
	if err != nil || len(hash) != 32 {
		return "", false
	}

	data := make([]byte, 34)
	data[0] = 0x51
	data[1] = byte(int8(workchain))
	copy(data[2:], hash)

	crc := tonCRC16(data)
	full := append(data, byte(crc>>8), byte(crc))

	return base64.RawURLEncoding.EncodeToString(full), true
}

// tonCRC16 计算 TON user-friendly 地址使用的 CRC16-CCITT 校验值。
func tonCRC16(data []byte) uint16 {
	var crc uint16
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}

	return crc
}

func IsValidEvmAddress(address string) bool {
	if len(address) != 42 || !strings.HasPrefix(address, "0x") {

		return false
	}

	addrWithoutPrefix := address[2:]
	if _, err := hex.DecodeString(addrWithoutPrefix); err != nil {

		return false
	}

	return true
}

func IsValidSolanaAddress(address string) bool {
	data := base58.Decode(address)

	return len(data) == 32
}

func IsValidAptosAddress(address string) bool {
	if strings.HasPrefix(address, "0x") || strings.HasPrefix(address, "0X") {

		address = address[2:]
	}

	if len(address) != 64 {

		return false
	}

	matched, _ := regexp.MatchString("^[0-9a-fA-F]{64}$", address)

	return matched
}

func MaskAddress(address string) string {
	if len(address) <= 20 {

		return address
	}

	return address[:8] + " ***** " + address[len(address)-10:]
}

func MaskAddress2(address string) string {
	if len(address) <= 20 {

		return address
	}

	return "*** " + address[len(address)-8:]
}

func MaskHash(hash string) string {
	if len(hash) <= 20 {

		return hash
	}

	return hash[:6] + " ***** " + hash[len(hash)-8:]
}

func CalcNextNotifyTime(base time.Time, num int) time.Time {

	return base.Add(time.Minute * time.Duration(math.Pow(2, float64(num))))
}

func HexStr2Int(str string) *big.Int {
	var n = new(big.Int)
	var val = strings.TrimLeft(strings.TrimPrefix(str, "0x"), "0")

	n.SetString(val, 16)

	return n
}

func InStrings(str string, list []string) bool {
	for _, item := range list {
		if item == str {

			return true
		}
	}

	return false
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

func StrSha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func IsAllowedCallbackURL(raw string) bool {
	// IsAllowedCallbackURL 校验回调/跳转地址格式是否合法
	// 规则：必须是合法 URL，且 scheme 只允许 http 或 https
	if raw == "" {
		return false
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}

// GetRequestHost 识别完整的请求主机地址
func GetRequestHost(r *http.Request) string {
	scheme := "http"

	if r.TLS != nil {
		scheme = "https"
	} else {
		// 检查各种代理转发的协议头
		// Nginx、Apache、Caddy 通用
		if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
			scheme = "https"
		} else if r.Header.Get("X-Forwarded-Ssl") == "on" {
			scheme = "https"
		} else if r.Header.Get("Front-End-Https") == "on" {
			// Apache mod_proxy
			scheme = "https"
		} else if r.Header.Get("X-Url-Scheme") == "https" {
			// Caddy
			scheme = "https"
		} else if r.Header.Get("CF-Visitor") != "" {
			// Cloudflare
			// CF-Visitor 格式: {"scheme":"https"}
			if strings.Contains(r.Header.Get("CF-Visitor"), `"scheme":"https"`) {
				scheme = "https"
			}
		}
	}

	return scheme + "://" + r.Host
}
