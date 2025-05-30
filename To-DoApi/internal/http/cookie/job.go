package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

type User struct {
	UserID   int64  `json:"id,omitempty" validate:"numeric,omitempty"`
	Username string `json:"username,omitempty" validate:"alphanum,required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func cashKey(key string) []byte {
	bKey := []byte(key)
	hKey := sha256.Sum256(bKey)
	return hKey[:]
}

func GetCookieUser(key string, dataUser User) (*http.Cookie, error) {
	user, err := json.Marshal(dataUser)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(cashKey(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherCookie := gcm.Seal(nonce, nonce, user, nil)
	dataCookie := base64.URLEncoding.EncodeToString(cipherCookie)

	cookie := &http.Cookie{
		Name:     "session_cookie",
		Value:    dataCookie,
		Path:     "/",
		MaxAge:   (3600 * 24) * 365,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie, nil
}

func ValidateCookieUser(key string, cookieValue string) (string, error) {
	value, err := base64.URLEncoding.DecodeString(cookieValue)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(cashKey(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce, cipher := value[:gcm.NonceSize()], value[gcm.NonceSize():]
	str, err := gcm.Open(nil, nonce, cipher, nil)
	if err != nil {
		return "", err
	}

	return string(str), nil
}
