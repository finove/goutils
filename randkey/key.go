package randkey

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

var numberUpperEncode = base32.NewEncoding("0123456789ABCDEFGHJKLMNPQRSTUVWX").WithPadding(base32.NoPadding)
var numberLowerEncode = base32.NewEncoding("0123456789abcdefghijkmnpqrstuvwx").WithPadding(base32.NoPadding)
var numberPassEncode = base64.NewEncoding("0123456789abcdefghijkmnpqrstuvwx!@#$%^&*()ABCDEFGHJKLMNPQRSTUVWX").WithPadding(base64.NoPadding)

func randNum() (num string) {
	var err error
	var b = make([]byte, 2)
	if _, err = rand.Read(b); err == nil {
		num = fmt.Sprintf("%05d", binary.BigEndian.Uint16(b))
	} else {
		num = fmt.Sprintf("%05d", time.Now().UnixNano()%100000)
	}
	return
}

// NumbersOnly 只返回数字
func NumbersOnly(count int) (code string) {
	var nums []string
	if count <= 0 {
		return
	}
	for i := 0; i < count; i += 5 {
		nums = append(nums, randNum())
	}
	code = strings.Join(nums, "")
	if len(code) > count {
		code = code[:count]
	}
	return
}

// NumberUpper 数字和大写字母
func NumberUpper(count int) (code string) {
	var err error
	var b = make([]byte, count)
	if _, err = rand.Read(b); err != nil {
		return
	}
	code = numberUpperEncode.EncodeToString(b)
	if len(code) > count {
		code = code[:count]
	}
	return
}

// NumberLower 数字和小写字母
func NumberLower(count int) (code string) {
	var err error
	var b = make([]byte, count)
	if _, err = rand.Read(b); err != nil {
		return
	}
	code = numberLowerEncode.EncodeToString(b)
	if len(code) > count {
		code = code[:count]
	}
	return
}

// NumberPass 数字和大写字母、特殊字符
func NumberPass(count int) (code string) {
	var err error
	var b = make([]byte, count)
	if _, err = rand.Read(b); err != nil {
		return
	}
	code = numberPassEncode.EncodeToString(b)
	if len(code) > count {
		code = code[:count]
	}
	return
}
