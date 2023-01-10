package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/spf13/viper"

	"web_app/validators"
)

func EncryptPassword(pwd string) string {
	h := md5.New()
	h.Write([]byte(viper.GetString("encrypt.secret")))
	return hex.EncodeToString(h.Sum([]byte(pwd)))
}

// 定义一个 removeTopStruct方法使返回错误的key只有具体字段，不含结构体名
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(validators.Trans)),
	})
	return
}

// escapeBytesBackslash escapes []byte with backslashes (\)
// This escapes the contents of a string (provided as []byte) by adding backslashes before special
// characters, and turning others into specific escape sequences, such as
// turning newlines into \n and null bytes into \0.
// https://github.com/mysql/mysql-server/blob/mysql-5.7.5/mysys/charset.c#L823-L932
func escapeBytesBackslash(v []byte) []byte {
	buf := bytes.NewBuffer(nil)
	for _, c := range v {
		switch c {
		case '\x00':
			buf.WriteByte('\\')
			buf.WriteByte('0')
		case '\n':
			buf.WriteByte('\\')
			buf.WriteByte('n')
		case '\r':
			buf.WriteByte('\\')
			buf.WriteByte('r')
		case '\x1a':
			buf.WriteByte('\\')
			buf.WriteByte('Z')
		case '\'':
			buf.WriteByte('\\')
			buf.WriteByte('\'')
		case '"':
			buf.WriteByte('\\')
			buf.WriteByte('"')
		case '\\':
			buf.WriteByte('\\')
			buf.WriteByte('\\')
		default:
			buf.WriteByte(c)
		}
	}
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result
}

// escapeStringBackslash is similar to escapeBytesBackslash but for string.
func escapeStringBackslash(v string) []byte {
	buf := bytes.NewBuffer(nil)
	for i := 0; i < len(v); i++ {
		c := v[i]
		switch c {
		case '\x00':
			buf.WriteByte('\\')
			buf.WriteByte('0')
		case '\n':
			buf.WriteByte('\\')
			buf.WriteByte('n')
		case '\r':
			buf.WriteByte('\\')
			buf.WriteByte('r')
		case '\x1a':
			buf.WriteByte('\\')
			buf.WriteByte('Z')
		case '\'':
			buf.WriteByte('\\')
			buf.WriteByte('\'')
		case '"':
			buf.WriteByte('\\')
			buf.WriteByte('"')
		case '\\':
			buf.WriteByte('\\')
			buf.WriteByte('\\')
		default:
			buf.WriteByte(c)
		}
	}
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result
}

func GenerateSqlLog(query string, args ...interface{}) string {
	// Number of ? should be same to len(args)
	if strings.Count(query, "?") != len(args) {
		return query + " args:" + fmt.Sprint(args) + " (param number not match)"
	}
	buf := bytes.NewBuffer(nil)

	argPos := 0

	for i := 0; i < len(query); i++ {
		q := strings.IndexByte(query[i:], '?')
		if q == -1 {
			buf.WriteString(query[i:])
			break
		}
		buf.WriteString(query[i : i+q])
		i += q

		arg := args[argPos]
		argPos++

		if arg == nil {
			buf.WriteString("NULL")
			continue
		}
		switch v := arg.(type) {
		case int64:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int32:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int16:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int8:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case uint64:
			// Handle uint64 explicitly because our custom ConvertValue emits unsigned values
			buf.WriteString(strconv.FormatUint(v, 10))
		case uint32:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint16:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint8:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case float64:
			buf.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
		case float32:
			buf.WriteString(strconv.FormatFloat(float64(v), 'g', -1, 32))
		case bool:
			if v {
				buf.WriteByte('1')
			} else {
				buf.WriteByte('0')
			}
		case []byte:
			if v == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString("_binary'")
				buf.Write(escapeBytesBackslash(v))
				buf.WriteByte('\'')
			}
		case string:
			buf.WriteByte('\'')
			buf.Write(escapeStringBackslash(v))
			buf.WriteByte('\'')
		default:
			rv := reflect.ValueOf(v)
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				buf.WriteString(strconv.FormatInt(rv.Int(), 10))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				buf.WriteString(strconv.FormatUint(rv.Uint(), 10))
			case reflect.Float64:
				buf.WriteString(strconv.FormatFloat(rv.Float(), 'g', -1, 64))
			case reflect.Float32:
				buf.WriteString(strconv.FormatFloat(rv.Float(), 'g', -1, 32))
			case reflect.Bool:
				buf.WriteString(strconv.FormatBool(rv.Bool()))
			case reflect.String:
				buf.WriteByte('\'')
				buf.Write(escapeStringBackslash(rv.String()))
				buf.WriteByte('\'')
			case reflect.Slice:
				if rv.Elem().Kind() == reflect.Uint8 {
					if b := rv.Bytes(); b == nil {
						buf.WriteString("NULL")
					} else {
						buf.WriteString("_binary'")
						buf.Write(escapeBytesBackslash(b))
						buf.WriteByte('\'')
					}
				} else {
					buf.WriteString(fmt.Sprintf("(unsupport driver type %T value %s)", v, v))
				}
			default:
				buf.WriteString(fmt.Sprintf("(unsupport driver type %T value %s)", v, v))
			}
		}

	}
	if argPos != len(args) {
		return query + " args:" + fmt.Sprint(args) + " (param number not match)"
	}
	return buf.String()
}
