package randkey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// PasswordValid hash验证密码
type PasswordValid string

// NewPasswordValid 构造一个新的验证器
func NewPasswordValid(key string) (p *PasswordValid) {
	p = new(PasswordValid)
	*p = PasswordValid(key)
	return
}

// Encode 加密密码
func (p *PasswordValid) Encode(password string) (enc string) {
	var encKey, randSeed, savePass string
	var crypt [32]byte
	if password == "" {
		return
	}
	encKey, randSeed = p.getEncKey()
	crypt = sha256.Sum256([]byte(fmt.Sprintf("%s%s%s", encKey, randSeed, password)))
	savePass = fmt.Sprintf("%s:%02x", randSeed, crypt)
	enc = base64.StdEncoding.EncodeToString([]byte(savePass))
	return
}

// Valid 验证是否一致
func (p *PasswordValid) Valid(originPassword, password string) (ok bool) {
	var keyPass []string
	var crypt [32]byte
	var encKey string
	ppp, err := base64.StdEncoding.DecodeString(originPassword)
	if err != nil {
		return
	}
	keyPass = strings.Split(string(ppp), ":")
	if len(keyPass) != 2 {
		return
	}
	encKey, _ = p.getEncKey()
	crypt = sha256.Sum256([]byte(fmt.Sprintf("%s%s%s", encKey, keyPass[0], password)))
	if fmt.Sprintf("%02x", crypt) == keyPass[1] {
		ok = true
	}
	return
}

func (p *PasswordValid) getEncKey() (k, seed string) {
	var tmpBuff = make([]byte, 8)
	n, err := rand.Read(tmpBuff)
	if err != nil || n != 8 {
		binary.BigEndian.PutUint64(tmpBuff, uint64(time.Now().UnixNano()))
	}
	seed = fmt.Sprintf("%02x", tmpBuff)
	k = string(*p)
	return
}
