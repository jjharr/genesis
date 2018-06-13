package validate

import (
	"fmt"
	"github.com/ansel1/merry"
	"reflect"
	"strconv"
	"strings"
)

// TODO [jjh]
// a regex validator that takes a pattern
//
// more dyanmic, non-tag way of specifying validations. Problem with tags is that they're static - can't compare
// to variable values for example, and not really a good way to compare to other fields. Maybe have a function
// that contains arrays of validators for each field? something like EmValidator.Op. Only the params need to be
// dynamic probably. For a dynamic IDE, we don't really care about the boilerplate, and we could create mirror
// structs with validation information. :
// x.AddValidation(struct.Field, []Validations{{IsUrl},{Between,5,10}, ...})

// tokens cannot have overlapping characters (e.g. => and = for different tokens)
const (
	tagName            = "valid"
	validatorSeparator = "|"
	settingsToken      = "="
	customMessageToken = "->"
	paramOpenToken     = "("
	paramCloseToken    = ")"

	// do not make this a backslash!
	paramSeparator = ","
)

// customValidatorsHolder is a singleton holder for all custom validations
type customValidatorsHolder map[string]map[string]string

func (h customValidatorsHolder) getKey(i interface{}, validatorName string) string {
	ty := reflect.TypeOf(i)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return fmt.Sprintf("%s %s %s", ty.PkgPath(), ty.Name(), validatorName)
}

// clear clears custom navigation, used only in tests
func (h *customValidatorsHolder) clear() {
	for k := range *h {
		delete(*h, k)
	}
}

func (h customValidatorsHolder) add(validatorName string, sampleStruct interface{}, validations map[string]string) error {
	ty := reflect.TypeOf(sampleStruct)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	for fieldName := range validations {
		if _, found := ty.FieldByName(fieldName); !found {
			return fmt.Errorf("Field %s in %s not found", fieldName, ty.Name())
		}
	}

	h[h.getKey(sampleStruct, validatorName)] = validations
	return nil
}

func (h customValidatorsHolder) get(validatorName string, i interface{}) (map[string]string, error) {
	key := h.getKey(i, validatorName)

	if validations, found := customValidators[key]; found {
		return validations, nil
	}
	return nil, fmt.Errorf("No custom validation %s", validatorName)
}

var customValidators = customValidatorsHolder(map[string]map[string]string{})

// MustAddCustomValidation adds a custom validation. This method must be called in init() methods and will panic if
// any of the fields isn't found in the sampleStruct
func MustAddCustomValidation(validatorName string, sampleStruct interface{}, validations map[string]string) {
	if err := customValidators.add(validatorName, sampleStruct, validations); err != nil {
		panic(err.Error())
	}
}

func AddCustomValidation(validatorName string, sampleStruct interface{}, validations map[string]string) error {
	return customValidators.add(validatorName, sampleStruct, validations)
}

// clearCustomValidation is used only in tests
func clearCustomValidation() {
	customValidators.clear()
}

// ValidateStruct use tags for fields.
// result will contain validation errors or be empty if there are none (HasErrors() returns false)
// error is set only if there is an internal Validator error (as opposed to a failed validation)
func ValidateStruct(s interface{}) (*ErrorBag, error) {
	return doValidateStruct(s, NewErrorBag(), nil)
}

// CustomValidateStruct validates the interfaces using custom validations (registered with AddCustomValidation)
// The default validation (defined with valid tags) can be used with "valid"
func CustomValidateStruct(s interface{}, customValidations ...string) (*ErrorBag, error) {
	bag := NewErrorBag()
	for _, customValidation := range customValidations {
		var validations map[string]string
		if customValidation != tagName {
			v, err := customValidators.get(customValidation, s)
			if err != nil {
				return nil, err
			}
			validations = v
		}
		_, err := doValidateStruct(s, bag, validations)
		if err != nil {
			return bag, err
		}
	}
	return bag, nil
}

// to allow recursive calling with the same error bag
func doValidateStruct(s interface{}, bag *ErrorBag, customFieldTags map[string]string) (*ErrorBag, error) {

	var internalError error

	// don't error out. since this is called recursively, it may be valid for a child struct to be nil.
	if s == nil {
		return bag, internalError
	}

	obj := reflect.ValueOf(s)
	if obj.Kind() == reflect.Interface || obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	// we only accept structs
	if obj.Kind() != reflect.Struct {
		return bag, fmt.Errorf("doValidateStruct only accepts structs; got %s", obj.Kind())
	}

	for i := 0; i < obj.NumField(); i++ {
		valueField := obj.Field(i)
		typeField := obj.Type().Field(i)
		if typeField.PkgPath != "" {
			continue // Private field
		}

		internalError = validateField(valueField, typeField, obj, bag, customFieldTags)
		if internalError != nil {
			break
		}
	}

	return bag, internalError
}

// parse struct field tags and return a slice of FieldValidators
//
// tags are always separated by commas
// there are two kinds of tags: validation directives and settings
// validation directives : are configured like : xxxx(a,b)=>this is a message, where the parameters and custom message parts are both optional
// settings : are configured like : xxxx=yyy
func getFieldValidators(v reflect.Value, t reflect.StructField, o reflect.Value, customFieldTags map[string]string) ([]FieldValidator, error) {

	fieldValidators := make([]FieldValidator, 0)

	var tag string
	if customFieldTags == nil {
		tag = t.Tag.Get(tagName)
	} else {
		tag = customFieldTags[t.Name]
	}

	rawKeys := strings.Split(tag, validatorSeparator)

	// set default field name using titlecase equivalent
	fieldName := strings.Title(t.Name)

	// process settings, only supported setting at the moment is field name
	rawKeys, name, err := extractFieldName(rawKeys)
	if len(name) > 0 {
		fieldName = name
	} else if err != nil {
		return nil, err
	}

	// handle validator directives
	for _, key := range rawKeys {
		if key == "-" || key == "" {
			continue
		}

		if len(key) == 0 {
			if customFieldTags != nil {
				// custom fields validatinos can be set for only a subset of fields
				continue
			}
			return nil, fmt.Errorf("Missing validator directive for field %s (nothing after the separator token).", fieldName)
		}

		validator := FieldValidator{
			FieldName:  fieldName,
			FieldValue: v.Interface(),
		}

		// after each operation for negation,message, and parameters we reset the key to exclude
		// the element we just processed
		if string(key[0]) == "!" {
			validator.IsNegated = true
			key = key[1:]
		}

		key, customMessagesPtr, err := extractMessage(key, fieldName, validator.IsNegated)
		if customMessagesPtr != nil {
			validator.FieldCustomMessages = *customMessagesPtr
		} else if err != nil {
			return nil, err
		}

		key, params, err := params(key, fieldName)

		if len(params) > 0 {
			validator.ValidatorParams = params
		} else if err != nil {
			return nil, err
		}

		v, ok := GetValidator(key)
		if !ok {
			return nil, fmt.Errorf("Invalid validation key for field %s: %s", fieldName, key)
		}

		validator.Validator = *v

		fieldValidators = append(fieldValidators, validator)
	}

	return fieldValidators, nil
}

func tmpLookup(tag string, key string) (value string, ok bool) {
	// When modifying this code, also update the validateStructTag code
	// in golang.org/x/tools/cmd/vet/structtag.go.

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}

// extract the field name from the tag, return the tag keys with the name removed, plus the field name and error
func extractFieldName(rawKeys []string) ([]string, string, error) {
	for idx, key := range rawKeys {

		settingTokenIndex := strings.Index(key, settingsToken)
		if settingTokenIndex > 0 {
			rawKeys = append(rawKeys[:idx], rawKeys[idx+1:]...)
			settingName := key[:settingTokenIndex]
			switch settingName {
			case "name":
				return rawKeys, key[settingTokenIndex+1:], nil
			default:
				return rawKeys, ``, fmt.Errorf("%s is not a valid validation option", settingName)
			}
		}
	}

	return rawKeys, ``, nil
}

func extractMessage(key string, fieldName string, negated bool) (string, *MessageSet, error) {
	customMessageIndex := strings.Index(key, customMessageToken)
	var ms MessageSet
	if customMessageIndex > 0 {
		cm := key[customMessageIndex+len(customMessageToken):]
		key = key[:customMessageIndex]
		if len(cm) == 0 {
			return key, nil, fmt.Errorf("Custom message indicated but not given for %s", fieldName)
		}
		isFmt := strings.Index(cm, "%s") > -1
		if !(negated && isFmt) {
			ms = MessageSet{Message: cm}
		} else if !negated {
			ms = MessageSet{MessageFmt: cm}
		} else if negated && !isFmt {
			ms = MessageSet{NegatedMessage: cm}
		} else {
			ms = MessageSet{NegatedMessageFmt: cm}
		}
		return key, &ms, nil
	}

	return key, nil, nil
}

// params return the parameters for validations functions that require them.
// key is the validation text, fieldName is the name of the struct field being validated
func params(validationText string, fieldName string) (string, []interface{}, error) {

	paramsOpenIndex := strings.Index(validationText, paramOpenToken)

	// if no open index, this is not a validator that requires params
	if paramsOpenIndex > 0 {

		validationKey := validationText[:paramsOpenIndex]

		paramsCloseIndex := strings.LastIndex(validationText, paramCloseToken)
		if paramsCloseIndex != (len(validationText) - len(paramCloseToken)) {
			return ``, make([]interface{}, 0, 0), merry.Errorf("The parameter close token for %s on field %s is incorrect or missing", validationText[:paramsOpenIndex], fieldName)
		}

		paramStr := validationText[paramsOpenIndex+len(paramOpenToken) : paramsCloseIndex]
		params := extractParams(paramStr)

		return validationKey, params, nil
	}

	return validationText, make([]interface{}, 0, 0), nil
}

// extract valid params from contents supplied to paramaterized validation func. Handles escaped paramSeparators
// for regexp functions
func extractParams(rawParamStr string) []interface{} {

	paramSeparatorIndex := strings.Index(rawParamStr, paramSeparator)

	if paramSeparatorIndex < 0 {
		return []interface{}{
			rawParamStr,
		}
	}

	// split instead of iterating over string to make multi-char separator easier
	re := strings.Split(rawParamStr, paramSeparator)
	params := make([]interface{}, 0, len(re))

	param := ``
	for _, v := range re {

		// skip escaped paramSeparator instances
		if string(v[len(v)-1]) == `\` && string(v[len(v)-2]) != `\\` {
			param += v + paramSeparator
			continue
		}

		param = param + v

		params = append(params, param)
		param = ``
	}

	return params
}

func validateField(v reflect.Value, t reflect.StructField, o reflect.Value, validationErrs *ErrorBag, customFieldTags map[string]string) error {

	if !v.IsValid() {
		return nil
	}

	fieldValidators, err := getFieldValidators(v, t, o, customFieldTags)
	if err != nil {
		return err
	}

	// todo add time.Time
	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		err = validateBasicType(v, t, fieldValidators, validationErrs)
		if err != nil {
			return err
		}

	case reflect.Map:
		if err := validateComplexType(v, t, fieldValidators, validationErrs); err != nil {
			return err
		}
		if err := validateMap(v, validationErrs, customFieldTags); err != nil {
			return err
		}
	case reflect.Slice:
		if err := validateComplexType(v, t, fieldValidators, validationErrs); err != nil {
			return err
		}
		if err := validateArrayOrSlice(v, t, o, validationErrs, customFieldTags); err != nil {
			return err
		}
	case reflect.Array:
		if err := validateComplexType(v, t, fieldValidators, validationErrs); err != nil {
			return err
		}
		if err := validateArrayOrSlice(v, t, o, validationErrs, customFieldTags); err != nil {
			return err
		}
	case reflect.Interface:
		if err := validateComplexType(v, t, fieldValidators, validationErrs); err != nil {
			return err
		}
		// If the value is an interface then encode its element
		if v.IsNil() {
			return nil
		}

		if _, err := doValidateStruct(v.Interface(), validationErrs, customFieldTags); err != nil {
			return err
		}
	case reflect.Ptr:
		if err := validateComplexType(v, t, fieldValidators, validationErrs); err != nil {
			return err
		}
		// If the value is a pointer then check its element
		if v.IsNil() {
			return nil
		}
		return validateField(v.Elem(), t, o, validationErrs, customFieldTags)
	case reflect.Struct:
		if err := validateComplexType(v, t, fieldValidators, validationErrs); err != nil {
			return err
		}
		if _, err = doValidateStruct(v.Interface(), validationErrs, customFieldTags); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s is not a supported validation type", v.Kind().String())
	}

	return nil
}

func errorKey(t reflect.StructField, v FieldValidator) string {

	// usually we don't want to use the field names as error bag keys. Try form and json first
	var errorBagKey = t.Tag.Get(`form`)
	if errorBagKey != `` {
		return errorBagKey
	}

	errorBagKey = t.Tag.Get(`json`)
	if errorBagKey != `` {
		return errorBagKey
	}

	return v.Validator.Key
}

func validateBasicType(v reflect.Value, t reflect.StructField, fieldValidators []FieldValidator, validationErrs *ErrorBag) error {

	for _, validator := range fieldValidators {

		valid, err := validator.Validator.Validate(v, validator.ValidatorParams)
		if err != nil {
			return fmt.Errorf("Error validating %s: %s", t.Name, err.Error())
		}

		if !valid {
			validationErrs.Add(errorKey(t, validator), validator.Message(), validator.FieldName)
		}
	}

	return nil
}

func validateComplexType(v reflect.Value, t reflect.StructField, fieldValidators []FieldValidator, validationErrs *ErrorBag) error {

	for _, validator := range fieldValidators {
		if validator.CanValidateComplexTypes() {
			valid, err := validator.Validator.Validate(v, validator.ValidatorParams)
			if err != nil {
				return fmt.Errorf("Error validating %s: %s", t.Name, err.Error())
			}

			if !valid {
				validationErrs.Add(errorKey(t, validator), validator.Message(), validator.FieldName)
			}
		}
	}

	return nil
}

// fixme - currently only works for maps where values are structs. modify to also handle basic types
func validateMap(v reflect.Value, validationErrs *ErrorBag, customFieldTags map[string]string) error {

	// check len

	if v.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("%s is not a supported validation type", v.Kind().String())
	}

	// check value type
	var sv = v.MapKeys()
	for _, k := range sv {
		if v.MapIndex(k).Kind() == reflect.Struct {
			_, err := doValidateStruct(v.MapIndex(k).Interface(), validationErrs, customFieldTags)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateArrayOrSlice(v reflect.Value, t reflect.StructField, o reflect.Value, validationErrs *ErrorBag, customFieldTags map[string]string) error {
	for i := 0; i < v.Len(); i++ {
		var err error
		if v.Index(i).Kind() != reflect.Struct {
			err = validateField(v.Index(i), t, o, validationErrs, customFieldTags)
			if err != nil {
				return err
			}
		} else {
			_, err = doValidateStruct(v.Index(i).Interface(), validationErrs, customFieldTags)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func IsEmptyValue(val interface{}) bool {
	return isEmptyValue(reflect.ValueOf(val))
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
