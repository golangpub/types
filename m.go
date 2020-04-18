package types

import (
	"encoding/json"
	"math/big"
	"net/mail"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/golangpub/conv"
	"github.com/golangpub/log"
	"github.com/shopspring/decimal"
)

// M is a special map which provides convenient methods
type M map[string]interface{}

func (m M) Slice(key string) []interface{} {
	value := m[key]
	if value == nil {
		return []interface{}{}
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		length := v.Len()
		var values = make([]interface{}, length)
		for i := 0; i < length; i++ {
			values[i] = v.Index(i).Interface()
		}
		return values
	default:
		return []interface{}{value}
	}
}

func (m M) Get(key string) interface{} {
	value := m[key]
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() > 0 {
			return v.Index(0).Interface()
		}
		return nil
	default:
		return value
	}
}

func (m M) Contains(key string) bool {
	return m.Get(key) != nil
}

func (m M) ContainsString(key string) bool {
	switch m.Get(key).(type) {
	case string:
		return true
	default:
		return false
	}
}

func (m M) String(key string) string {
	switch v := m.Get(key).(type) {
	case string:
		return v
	case json.Number:
		return string(v)
	default:
		return ""
	}
}

func (m M) TrimmedString(key string) string {
	switch v := m.Get(key).(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return string(v)
	default:
		return ""
	}
}

func (m M) DefaultString(key string, defaultValue string) string {
	switch v := m.Get(key).(type) {
	case string:
		return v
	default:
		return defaultValue
	}
}

func (m M) MustString(key string) string {
	switch v := m.Get(key).(type) {
	case string:
		return v
	default:
		panic("No string value for key:" + key)
	}
}

func (m M) StringSlice(key string) []string {
	_, found := m[key]
	if !found {
		return nil
	}

	values := m.Slice(key)
	var result []string
	for _, v := range values {
		if str, ok := v.(string); ok {
			result = append(result, str)
		}
	}

	return result
}

func (m M) ContainsBool(key string) bool {
	_, err := conv.ToBool(m.Get(key))
	return err == nil
}

func (m M) Bool(key string) bool {
	v, _ := conv.ToBool(m.Get(key))
	return v
}

func (m M) DefaultBool(key string, defaultValue bool) bool {
	if v, err := conv.ToBool(m.Get(key)); err == nil {
		return v
	}
	return defaultValue
}

func (m M) MustBool(key string) bool {
	if v, err := conv.ToBool(m.Get(key)); err == nil {
		return v
	}
	panic("No bool value for key:" + key)
}

func (m M) Int(key string) int {
	v, _ := conv.ToInt64(m.Get(key))
	return int(v)
}

func (m M) DefaultInt(key string, defaultVal int) int {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return int(v)
	}
	return defaultVal
}

func (m M) MustInt(key string) int {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return int(v)
	}
	panic("No int value for key:" + key)
}

func (m M) IntSlice(key string) []int {
	l, _ := conv.ToIntSlice(m.Slice(key))
	return l
}

func (m M) ContainsInt64(key string) bool {
	_, err := conv.ToInt64(m.Get(key))
	return err == nil
}

func (m M) Int64(key string) int64 {
	v, _ := conv.ToInt64(m.Get(key))
	return v
}

func (m M) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return v
	}
	return defaultVal
}

func (m M) MustInt64(key string) int64 {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return v
	}
	panic("No int64 value for key:" + key)
}

func (m M) Int64Slice(key string) []int64 {
	values := m.Slice(key)
	var result []int64
	for _, v := range values {
		i, e := conv.ToInt64(v)
		if e == nil {
			result = append(result, i)
		}
	}

	return result
}

func (m M) ContainsFloat64(key string) bool {
	_, err := conv.ToFloat64(m.Get(key))
	return err == nil
}

func (m M) Float64(key string) float64 {
	v, _ := conv.ToFloat64(m.Get(key))
	return v
}

func (m M) DefaultFloat64(key string, defaultValue float64) float64 {
	if v, err := conv.ToFloat64(m.Get(key)); err == nil {
		return v
	}
	return defaultValue
}

func (m M) MustFloat64(key string) float64 {
	if v, err := conv.ToFloat64(m.Get(key)); err == nil {
		return v
	}
	panic("No float64 value for key:" + key)
}

func (m M) Float64Slice(key string) []float64 {
	values := m.Slice(key)
	var result []float64
	for _, val := range values {
		i, e := conv.ToFloat64(val)
		if e == nil {
			result = append(result, i)
		}
	}

	return result
}

func (m M) BigInt(key string) *big.Int {
	s := m.String(key)
	n := big.NewInt(0)
	_, ok := n.SetString(s, 10)
	if !ok {
		_, ok = n.SetString(s, 16)
	}

	if ok {
		return n
	}

	if k, err := conv.ToInt64(m[key]); err == nil {
		return big.NewInt(k)
	}

	return nil
}

func (m M) DefaultBigInt(key string, defaultVal *big.Int) *big.Int {
	if n := m.BigInt(key); n != nil {
		return n
	}
	return defaultVal
}

func (m M) MustBigInt(key string) *big.Int {
	if n := m.BigInt(key); n != nil {
		return n
	}
	panic("No big.Int64 value for key:" + key)
}

func (m M) BigFloat(key string) *big.Float {
	s := m.String(key)
	n := big.NewFloat(0)
	_, ok := n.SetString(s)
	if ok {
		return n
	}

	if k, err := conv.ToFloat64(m[key]); err == nil {
		return big.NewFloat(k)
	}

	return nil
}

func (m M) DefaultBigFloat(key string, defaultVal *big.Float) *big.Float {
	if n := m.BigFloat(key); n != nil {
		return n
	}
	return defaultVal
}

func (m M) MustBigFloat(key string) *big.Float {
	if n := m.BigFloat(key); n != nil {
		return n
	}
	panic("No big.Float64 value for key:" + key)
}

func (m M) ContainsDecimal(key string) bool {
	s := m.String(key)
	if s == "" {
		return false
	}
	_, err := decimal.NewFromString(s)
	return err == nil
}

func (m M) Decimal(key string) decimal.Decimal {
	v, _ := decimal.NewFromString(m.String(key))
	return v
}

func (m M) DefaultDecimal(key string, defaultVal decimal.Decimal) decimal.Decimal {
	if v, err := decimal.NewFromString(m.String(key)); err == nil {
		return v
	}
	return defaultVal
}

func (m M) MustDecimal(key string) decimal.Decimal {
	if v, err := decimal.NewFromString(m.String(key)); err == nil {
		return v
	}
	panic("No decimal value for key:" + key)
}

func (m M) Map(key string) M {
	switch v := m.Get(key).(type) {
	case M:
		return v
	case map[string]interface{}:
		return v
	default:
		return M{}
	}
}

func (m M) Date(key string) (*time.Time, bool) {
	return m.DateInLocation(key, time.UTC)
}

func (m M) DateInLocation(key string, loc *time.Location) (*time.Time, bool) {
	s := strings.TrimSpace(m.String(key))
	if s == "" {
		return nil, false
	}
	d, err := conv.ToDateInLocation(s, loc)
	if err != nil {
		return nil, false
	}
	return d, true
}

func (m M) PhoneNumber(key string) *PhoneNumber {
	switch v := m[key].(type) {
	case string:
		pn, _ := NewPhoneNumber(v)
		return pn
	case M, map[string]interface{}:
		data, err := json.Marshal(v)
		if err != nil {
			log.Errorf("Marshal failed: %v", err)
		}
		var pn *PhoneNumber
		err = json.Unmarshal(data, &pn)
		if err != nil {
			log.Errorf("Unmarshal failed: %v", err)
		}
		pn, _ = NewPhoneNumber(pn.String())
		return pn
	default:
		return nil
	}
}

func (m M) EmailAddress(key string) *mail.Address {
	s := m.String(key)
	s = strings.TrimSpace(s)
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return nil
	}
	return addr
}

func (m M) URL(key string) string {
	s := m.String(key)
	s = strings.TrimSpace(s)
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	if len(u.Scheme) == 0 {
		return ""
	}

	if len(u.Host) == 0 {
		return ""
	}
	return s
}

func (m M) SetNX(k string, v interface{}) {
	if v == nil {
		return
	}
	if _, ok := m[k]; ok {
		return
	}
	m[k] = v
}

func (m M) AddMap(val M) {
	for k, v := range val {
		m.SetNX(k, v)
	}
}

func (m M) JSONString() string {
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}

func (m M) Remove(keys ...string) {
	for k := range m {
		if indexOfStr(keys, k) < 0 {
			delete(m, k)
		}
	}
}

func (m M) Keep(keys ...string) {
	for k := range m {
		if indexOfStr(keys, k) < 0 {
			delete(m, k)
		}
	}
}

func indexOfStr(l []string, s string) int {
	for i, str := range l {
		if s == str {
			return i
		}
	}
	return -1
}
