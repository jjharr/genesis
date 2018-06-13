package validate

/*
The MIT License (MIT)

Copyright (c) 2014 Alex Saskevich

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func IsEmail(val string, params ...interface{}) bool {
	// TODO uppercase letters are not supported
	return rxEmail.MatchString(val)
}

// IsURL check if the string is an URL.
func IsURL(str string, params ...interface{}) bool {

	// don't invalidate for zero length. Use 'required' validator for that
	if str == `` {
		return true
	}

	if len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return rxURL.MatchString(str)

}

// IsRequestURL check if the string rawurl, assuming
// it was recieved in an HTTP request, is a valid
// URL confirm to RFC 3986
func IsRequestURL(rawurl string, params ...interface{}) bool {

	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return false //Couldn't even parse the rawurl
	}
	if len(url.Scheme) == 0 {
		return false //No Scheme found
	}
	return true
}

// IsRequestURI check if the string rawurl, assuming
// it was recieved in an HTTP request, is an
// absolute URI or an absolute path.
func IsRequestURI(rawurl string, params ...interface{}) bool {
	_, err := url.ParseRequestURI(rawurl)
	return err == nil
}

// IsAlpha check if the string contains only letters (a-zA-Z). Empty string is valid.
func IsAlpha(str string, params ...interface{}) bool {

	if IsNull(str) {
		return true
	}
	return rxAlpha.MatchString(str)
}

//IsUTFLetter check if the string contains only unicode letter characters.
//Similar to IsAlpha but for all languages. Empty string is valid.
func IsUTFLetter(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}

	for _, c := range str {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true

}

func IsTitle(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxTitle.MatchString(str)
}

func IsName(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxName.MatchString(str)
}

func IsPhone(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxPhone.MatchString(str)
}

func IsSkype(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxSkype.MatchString(str)
}

// IsAlphanumeric check if the string contains only letters and numbers. Empty string is valid.
func IsAlphanumeric(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphanumeric.MatchString(str)
}

// IsUTFLetterNumeric check if the string contains only unicode letters and numbers. Empty string is valid.
func IsUTFLetterNumeric(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	for _, c := range str {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) { //letters && numbers are ok
			return false
		}
	}
	return true

}

// IsNumeric check if the string contains only numbers. Empty string is valid.
func IsNumeric(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxNumeric.MatchString(str)
}

// IsUTFNumeric check if the string contains only unicode numbers of any kind.
// Numbers can be 0-9 but also Fractions ¾,Roman Ⅸ and Hangzhou 〩. Empty string is valid.
func IsUTFNumeric(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	if strings.IndexAny(str, "+-") > 0 {
		return false
	}
	if len(str) > 1 {
		str = strings.TrimPrefix(str, "-")
		str = strings.TrimPrefix(str, "+")
	}
	for _, c := range str {
		if unicode.IsNumber(c) == false { //numbers && minus sign are ok
			return false
		}
	}
	return true

}

// IsUTFDigit check if the string contains only unicode radix-10 decimal digits. Empty string is valid.
func IsUTFDigit(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	if strings.IndexAny(str, "+-") > 0 {
		return false
	}
	if len(str) > 1 {
		str = strings.TrimPrefix(str, "-")
		str = strings.TrimPrefix(str, "+")
	}
	for _, c := range str {
		if !unicode.IsDigit(c) { //digits && minus sign are ok
			return false
		}
	}
	return true

}

// IsHexadecimal check if the string is a hexadecimal number.
func IsHexadecimal(str string, params ...interface{}) bool {
	return rxHexadecimal.MatchString(str)
}

// IsHexcolor check if the string is a hexadecimal color.
func IsHexcolor(str string, params ...interface{}) bool {
	return rxHexcolor.MatchString(str)
}

// IsRGBcolor check if the string is a valid RGB color in form rgb(RRR, GGG, BBB).
func IsRGBcolor(str string, params ...interface{}) bool {
	return rxRGBcolor.MatchString(str)
}

// IsLowerCase check if the string is lowercase. Empty string is valid.
func IsLowerCase(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return str == strings.ToLower(str)
}

// IsUpperCase check if the string is uppercase. Empty string is valid.
func IsUpperCase(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return str == strings.ToUpper(str)
}

// IsInt check if the string is an integer. Empty string is valid.
func IsInt(val interface{}, params ...interface{}) bool {

	switch v := val.(type) {
	case string:
		if IsNull(v) {
			return true
		}
		return rxInt.MatchString(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// IsFloat check if the string is a float.
func IsFloat(val interface{}, params ...interface{}) bool {

	str, ok := val.(string)

	if ok {
		return str != "" && rxFloat.MatchString(str)
	}

	val, ok = val.(float32)
	if ok {
		return true
	}

	return false
}

///////////////////
// Old Converters content
//////////////////

func toString(obj interface{}) string {
	res := fmt.Sprintf("%v", obj)
	return string(res)
}

// ToJSON convert the input to a valid JSON string
func toJSON(obj interface{}) (string, error) {
	res, err := json.Marshal(obj)
	if err != nil {
		res = []byte("")
	}
	return string(res), err
}

// ToFloat convert the input to a float
func toFloat(obj interface{}) (float64, error) {

	switch v := obj.(type) {
	case string:
		return strconv.ParseFloat(obj.(string), 64)
	// don't include int64,uint64 since it may not be possible to accurately convert int64 to float64
	case int, int8, int16, int32:
		return float64(reflect.ValueOf(v).Int()), nil
	case uint, uint8, uint16, uint32:
		return float64(reflect.ValueOf(v).Uint()), nil
	case float32, float64:
		return reflect.ValueOf(v).Float(), nil
	default:
		return 0.0, fmt.Errorf("Can't convert type %s to float", reflect.ValueOf(obj).Kind().String())
	}
}

// ToInt convert the input to an integer
func toInt(obj interface{}) (int64, error) {

	switch v := obj.(type) {
	case string:
		return strconv.ParseInt(obj.(string), 0, 64)
	// don't include int64,uint64 since it may not be possible to accurately convert int64 to float64
	case int, int8, int16, int32:
		return reflect.ValueOf(v).Int(), nil
	case uint, uint8, uint16, uint32:
		return int64(reflect.ValueOf(v).Uint()), nil
	case float32, float64:
		return int64(reflect.ValueOf(v).Float()), nil
	default:
		return 0, fmt.Errorf("Can't convert type %s to int", reflect.ValueOf(obj).Kind().String())
	}
}

// ToBoolean convert the input to a boolean. Uses strconv.ParseBool rules applied to string or (u)int
func toBoolean(obj interface{}) (bool, error) {

	switch v := obj.(type) {
	case string:
		return strconv.ParseBool(obj.(string))
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		return v == 1, nil
	default:
		return false, fmt.Errorf("Can't convert type %s to bool", reflect.ValueOf(obj).Kind().String())
	}
}

// TODO - create methods for common patterns, then use Matcher.matches for performance

// Matches check if string matches the pattern (pattern is regular expression)
// In case of error return false
func Matches(str, pattern string) bool {
	match, _ := regexp.MatchString(regexp.QuoteMeta(pattern), str)
	return match
}

// IsDivisibleBy check if the string is a number that's divisible by another.
// If second argument is not valid integer or zero, it's return false.
// Otherwise, if first argument is not valid integer or zero, it's return true (Invalid string converts to zero).
/*
func IsDivisibleBy(val interface{}, params ...interface{}) bool {

	v := reflect.TypeOf(val)
	switch v.Kind() {
	case reflect.String:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64 :
	case reflect.Float32, reflect.Float64
	}

	val,ok = val.(float32)
	if ok {
		fmt.Println("IsFloat, float")
		return true
	}

	f, _ := ToFloat(str)
	p := int64(f)
	q, _ := ToInt(num)
	if q == 0 {
		return false
	}
	return (p == 0) || (p%q == 0)
}
*/

func IsNonEmpty(val interface{}, params ...interface{}) bool {
	switch t := val.(type) {
	case int:
		return t > 0
	case int8:
		return t > 0
	case int16:
		return t > 0
	case int32:
		return t > 0
	case int64:
		return t > 0
	case string:
		return len(strings.TrimSpace(t)) > 0
	case float32:
		return t > 0
	case float64:
		return t > 0
	}

	len := getLen(val)
	if len < 0 {
		if val == nil {
			return false
		}
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Ptr { // Special case for pointers ("== nil" wil return false even when the pointer is nil)
			return !reflect.ValueOf(val).IsNil()
		}
		return val != nil
	}
	return len > 0
}

func getLen(i interface{}) int {
	if i == nil {
		return -1
	}
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice || v.Kind() == reflect.Map || v.Kind() == reflect.String {
		return v.Len()
	}
	return -1
}

// IsNull check if the string is null. Flexible definition of null...
func IsNull(val interface{}, params ...interface{}) bool {

	if val == nil {
		return true
	}

	str, ok := val.(string)
	if !ok {
		return false
	}

	return len(str) == 0 || strings.ToLower(str) == `null`
}

// IsUUIDv3 check if the string is a UUID version 3.
func IsUUIDv3(str string, params ...interface{}) bool {
	return rxUUID3.MatchString(str)
}

// IsUUIDv4 check if the string is a UUID version 4.
func IsUUIDv4(str string, params ...interface{}) bool {
	return rxUUID4.MatchString(str)
}

// IsUUIDv5 check if the string is a UUID version 5.
func IsUUIDv5(str string, params ...interface{}) bool {
	return rxUUID5.MatchString(str)
}

// IsUUID check if the string is a UUID (version 3, 4 or 5).
func IsUUID(str string, params ...interface{}) bool {
	return rxUUID.MatchString(str)
}

// IsCreditCard check if the string is a credit card.
func IsCreditCard(str string, params ...interface{}) bool {
	r, _ := regexp.Compile("[^0-9]+")
	sanitized := r.ReplaceAll([]byte(str), []byte(""))
	if !rxCreditCard.MatchString(string(sanitized)) {
		return false
	}
	var sum int64
	var digit string
	var tmpNum int64
	var shouldDouble bool
	for i := len(sanitized) - 1; i >= 0; i-- {
		digit = string(sanitized[i:(i + 1)])
		tmpNum, _ = toInt(digit)
		if shouldDouble {
			tmpNum *= 2
			if tmpNum >= 10 {
				sum += ((tmpNum % 10) + 1)
			} else {
				sum += tmpNum
			}
		} else {
			sum += tmpNum
		}
		shouldDouble = !shouldDouble
	}

	if sum%10 == 0 {
		return true
	}
	return false
}

// IsISBN10 check if the string is an ISBN version 10.
func IsISBN10(str string, params ...interface{}) bool {
	return IsISBN(str, 10)
}

// IsISBN13 check if the string is an ISBN version 13.
func IsISBN13(str string, params ...interface{}) bool {
	return IsISBN(str, 13)
}

// IsISBN check if the string is an ISBN (version 10 or 13).
// If version value is not equal to 10 or 13, it will be check both variants.
func IsISBN(str string, version int) bool {
	r, _ := regexp.Compile("[\\s-]+")
	sanitized := r.ReplaceAll([]byte(str), []byte(""))
	var checksum int32
	var i int32
	if version == 10 {
		if !rxISBN10.MatchString(string(sanitized)) {
			return false
		}
		for i = 0; i < 9; i++ {
			checksum += (i + 1) * int32(sanitized[i]-'0')
		}
		if sanitized[9] == 'X' {
			checksum += 10 * 10
		} else {
			checksum += 10 * int32(sanitized[9]-'0')
		}
		if checksum%11 == 0 {
			return true
		}
		return false
	} else if version == 13 {
		if !rxISBN13.MatchString(string(sanitized)) {
			return false
		}
		factor := []int32{1, 3}
		for i = 0; i < 12; i++ {
			checksum += factor[i%2] * int32(sanitized[i]-'0')
		}
		if (int32(sanitized[12]-'0'))-((10-(checksum%10))%10) == 0 {
			return true
		}
		return false
	}
	return IsISBN(str, 10) || IsISBN(str, 13)
}

// IsJSON check if the string is valid JSON (note: uses json.Unmarshal).
func IsJSON(str string, params ...interface{}) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// IsMultibyte check if the string contains one or more multibyte chars. Empty string is valid.
func IsMultibyte(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxMultibyte.MatchString(str)
}

// IsASCII check if the string contains ASCII chars only. Empty string is valid.
func IsASCII(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxASCII.MatchString(str)
}

// IsPrintableASCII check if the string contains printable ASCII chars only. Empty string is valid.
func IsPrintableASCII(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxPrintableASCII.MatchString(str)
}

// IsFullWidth check if the string contains any full-width chars. Empty string is valid.
func IsFullWidth(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxFullWidth.MatchString(str)
}

// IsHalfWidth check if the string contains any half-width chars. Empty string is valid.
func IsHalfWidth(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxHalfWidth.MatchString(str)
}

// IsVariableWidth check if the string contains a mixture of full and half-width chars. Empty string is valid.
func IsVariableWidth(str string, params ...interface{}) bool {
	if IsNull(str) {
		return true
	}
	return rxHalfWidth.MatchString(str) && rxFullWidth.MatchString(str)
}

// IsBase64 check if a string is base64 encoded.
func IsBase64(str string, params ...interface{}) bool {
	return rxBase64.MatchString(str)
}

// IsFilePath check is a string is Win or Unix file path and returns it's type.
func IsFilePath(val interface{}, params ...interface{}) (bool, int) {

	str, ok := val.(string)
	if !ok {
		return false, Unknown
	}

	if rxWinPath.MatchString(str) {
		//check windows path limit see:
		//  http://msdn.microsoft.com/en-us/library/aa365247(VS.85).aspx#maxpath
		if len(str[3:]) > 32767 {
			return false, Win
		}
		return true, Win
	} else if rxUnixPath.MatchString(str) {
		return true, Unix
	}
	return false, Unknown
}

// IsDataURI checks if a string is base64 encoded data URI such as an image
func IsDataURI(str string, params ...interface{}) bool {
	dataURI := strings.Split(str, ",")
	if !rxDataURI.MatchString(dataURI[0]) {
		return false
	}
	return IsBase64(dataURI[1])
}

// IsISO3166Alpha2 checks if a string is valid two-letter country code
func IsISO3166Alpha2(str string, params ...interface{}) bool {

	for _, entry := range ISO3166List {
		if str == entry.Alpha2Code {
			return true
		}
	}
	return false
}

// IsISO3166Alpha3 checks if a string is valid three-letter country code
func IsISO3166Alpha3(str string, params ...interface{}) bool {

	for _, entry := range ISO3166List {
		if str == entry.Alpha3Code {
			return true
		}
	}
	return false
}

// IsDNSName will validate the given string as a DNS name
func IsDNSName(str string, params ...interface{}) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		// constraints already violated
		return false
	}
	return rxDNSName.MatchString(str)
}

// IsDialString validates the given string for usage with the various Dial() functions
func IsDialString(str string, params ...interface{}) bool {
	if h, p, err := net.SplitHostPort(str); err == nil && h != "" && p != "" && (IsDNSName(h) || IsIP(h)) && IsPort(p) {
		return true
	}

	return false
}

// IsIP checks if a string is either IP version 4 or 6.
func IsIP(str string, params ...interface{}) bool {
	return net.ParseIP(str) != nil
}

// IsPort checks if a string represents a valid port
func IsPort(val interface{}, params ...interface{}) bool {

	switch val.(type) {
	case string:
		x, _ := val.(string)
		if i, err := strconv.Atoi(x); err == nil && i > 0 && i < 65536 {
			return true
		}
		return false
	case int, int8, int16, int64, uint, uint8, uint16, uint32, uint64:
		x, _ := val.(int)
		return x > 0 && x < 65536
	default:
		return false
	}
}

// IsIPv4 check if the string is an IP version 4.
func IsIPv4(str string, params ...interface{}) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ".")
}

// IsIPv6 check if the string is an IP version 6.
func IsIPv6(val interface{}, params ...interface{}) bool {

	str, ok := val.(string)
	if !ok {
		return false
	}

	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

// IsMAC check if a string is valid MAC address.
// Possible MAC formats:
// 01:23:45:67:89:ab
// 01:23:45:67:89:ab:cd:ef
// 01-23-45-67-89-ab
// 01-23-45-67-89-ab-cd-ef
// 0123.4567.89ab
// 0123.4567.89ab.cdef
func IsMAC(str string, params ...interface{}) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// IsHost checks if the string is a valid IP (both v4 and v6) or a valid DNS name
func IsHost(str string, params ...interface{}) bool {
	return IsIP(str) || IsDNSName(str)
}

// IsMongoID check if the string is a valid hex-encoded representation of a MongoDB ObjectId.
func IsMongoID(val interface{}, params ...interface{}) bool {

	str, ok := val.(string)
	if !ok {
		return false
	}

	return rxHexadecimal.MatchString(str) && (len(str) == 24)
}

// todo - handle float
// IsLatitude check if a string is valid latitude.
func IsLatitude(str string, params ...interface{}) bool {
	return rxLatitude.MatchString(str)
}

// IsLongitude check if a string is valid longitude.
func IsLongitude(str string, params ...interface{}) bool {
	return rxLongitude.MatchString(str)
}

// IsSSN will validate the given string as a U.S. Social Security Number
func IsSSN(str string, params ...interface{}) bool {
	if str == "" || len(str) != 11 {
		return false
	}
	return rxSSN.MatchString(str)
}

// IsSemver check if string is valid semantic version
func IsSemver(str string, params ...interface{}) bool {
	return rxSemver.MatchString(str)
}

// ByteLength check if the string's length (in bytes) falls in a range.
func ByteLength(val interface{}, params ...interface{}) bool {

	str, ok := val.(string)
	if !ok {
		return false
	}

	if len(params) != 2 {
		return false
	}

	var err error
	var min, max int
	min, err = getIntOrError(params[0])
	if err != nil {
		return false
	}
	max, err = getIntOrError(params[1])
	if err != nil {
		return false
	}

	return len(str) >= min && len(str) <= max
}

func getIntOrError(i interface{}) (int, error) {
	switch p := i.(type) {
	case int:
		return p, nil
	case string:
		i, err := strconv.ParseInt(p, 10, 32)
		if err != nil {
			return 0, err
		}
		return int(i), nil
	}
	return 0, fmt.Errorf("Invalid type %v", i)
}

// StringMatches checks if a string matches a given pattern.
func StringMatches(val string, params ...interface{}) bool {

	if len(params) != 1 {
		return false
	}

	pattern, ok := params[0].(string)
	if !ok {
		return false
	}

	return Matches(val, pattern)
}

// Between check params's length (including multi byte for strings) against supplied parameters. Parameters
// are inclusive. Expects ints for params. Handles string, int types for val.
func Between(val interface{}, params ...interface{}) bool {

	// better to use between or min/max?
	// between enforces min/max which is good
	// but we can't tell user which is wrong
	// possibly return 3 values?

	if len(params) != 2 {
		return false
	}

	// we ignore empty values in this validator. Otherwise, it'd be the same as required, which it's not.
	// use required if you need required!
	switch val.(type) {
	case string:
		x, _ := val.(string)

		if len(x) == 0 {
			return true
		}

		min, max, err := minMaxToInt(params[0], params[1])
		if err != nil {
			return false
		}

		strLength := int64(utf8.RuneCountInString(x))
		return strLength >= min && strLength <= max
	case int, int8, int16, int64, uint, uint8, uint16, uint32, uint64:
		x, _ := toInt(val)

		if x == 0 {
			return true
		}

		min, max, err := minMaxToInt(params[0], params[1])
		if err != nil {
			return false
		}

		return x >= min && x <= max
	case float32, float64:
		x, _ := toFloat(val)

		if x == 0.0 {
			return true
		}

		min, max, err := minMaxToFloat(params[0], params[1])
		if err != nil {
			return false
		}

		return x >= min && x <= max
	//case time.Time:

	default:
		return false
	}

	return false
}

func minMaxToInt(min, max interface{}) (int64, int64, error) {
	intMin, err := toInt(min)
	if err != nil {
		return 0, 0, err
	}

	intMax, err := toInt(max)
	if err != nil {
		return 0, 0, err
	}

	return intMin, intMax, nil
}

func minMaxToFloat(min, max interface{}) (float64, float64, error) {
	floatMin, err := toFloat(min)
	if err != nil {
		return 0.0, 0.0, err
	}

	floatMax, err := toFloat(max)
	if err != nil {
		return 0.0, 0.0, err
	}

	return floatMin, floatMax, nil
}

// Abs returns absolute value of number
func Abs(value float64) float64 {
	return value * Sign(value)
}

// Sign returns signum of number: 1 in case of value > 0, -1 in case of value < 0, 0 otherwise
func Sign(value float64) float64 {
	if value > 0 {
		return 1
	} else if value < 0 {
		return -1
	} else {
		return 0
	}
}

// IsNegative returns true if value < 0
func IsNegative(value float64) bool {
	return value < 0
}

// IsPositive returns true if value > 0
func IsPositive(value float64) bool {
	return value > 0
}

// IsNonNegative returns true if value >= 0
func IsNonNegative(value float64) bool {
	return value >= 0
}

// IsNonPositive returns true if value <= 0
func IsNonPositive(value float64) bool {
	return value <= 0
}

// InRange returns true if value lies between left and right border
//func InRange(value, left, right float64) bool {
//	if left > right {
//		left, right = right, left
//	}
//	return value >= left && value <= right
//}

// IsWhole returns true if value is whole number
func IsWhole(value float64) bool {
	return Abs(math.Remainder(value, 1)) == 0
}
