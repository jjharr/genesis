package validate

import (
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	"reflect"
	"regexp"
	"strings"
)

var emKeyMap = hashmap.New()

var customMessageVarRegex = regexp.MustCompile("{.*?}")

// Messages take precedence over MessageFmts
type MessageSet struct {
	Message           string
	MessageFmt        string
	NegatedMessage    string
	NegatedMessageFmt string
}

type EmValidator struct {
	Key string

	// Not about Op* fields. If the Op is != nil, it will be executed regardless of Op* validators. Otherwise
	// the appropriate Op* will be executed by type. If you add a new Op* type, don't forget to add a switch in
	// EmValidator.Validate()!!!
	Op                      func(val interface{}, params ...interface{}) bool
	OpString                func(val string, params ...interface{}) bool
	CanValidateComplexTypes bool

	DefaultMessages         MessageSet
	ValidatorCustomMessages MessageSet
}

func (ev EmValidator) Validate(v reflect.Value, params []interface{}) (bool, error) {

	if ev.Op != nil {
		return ev.Op(v.Interface(), params...), nil
	}

	switch t := v.Interface().(type) {
	case string:
		if ev.OpString != nil {
			return ev.OpString(t, params...), nil
		}
	}

	return false, fmt.Errorf("No default validator for field of type %s", v.Type().Name())
}

type FieldValidator struct {
	FieldName           string
	FieldValue          interface{}
	Validator           EmValidator
	ValidatorParams     []interface{}
	IsNegated           bool
	FieldCustomMessages MessageSet
}

func (ms FieldValidator) CanValidateComplexTypes() bool {
	return ms.Validator.CanValidateComplexTypes
}

func (ms FieldValidator) Message() string {

	if len(ms.FieldCustomMessages.Message) > 0 {
		return ms.fillMessagePlaceholders(ms.FieldCustomMessages.Message)
	} else if len(ms.FieldCustomMessages.MessageFmt) > 0 {
		return fmt.Sprintf(ms.FieldCustomMessages.MessageFmt, ms.FieldName)
	} else if len(ms.Validator.DefaultMessages.Message) > 0 {
		return ms.fillMessagePlaceholders(ms.Validator.DefaultMessages.Message)
	}

	return fmt.Sprintf(ms.Validator.DefaultMessages.MessageFmt, ms.FieldName)
}

func (ms FieldValidator) fillMessagePlaceholders(msg string) string {
	if strings.Index(msg, "{") < 0 {
		return msg
	}
	replacements := map[string]interface{}{
		"field": ms.FieldName,
		"value": ms.FieldValue,
	}
	return customMessageVarRegex.ReplaceAllStringFunc(msg, func(val string) string {
		key := val[1 : len(val)-1]
		formatString := "%v"
		parts := strings.Split(key, ":")
		if len(parts) > 1 {
			key = parts[0]
			formatString = parts[1]
		}
		if replacement, found := replacements[key]; found {
			return fmt.Sprintf(formatString, replacement)
		}
		return "???"
	})
}

func (ms FieldValidator) NegatedMessage() string {

	if len(ms.FieldCustomMessages.NegatedMessage) > 0 {
		return ms.FieldCustomMessages.Message
	} else if len(ms.FieldCustomMessages.NegatedMessageFmt) > 0 {
		return fmt.Sprintf(ms.FieldCustomMessages.MessageFmt, ms.FieldName)
	} else if len(ms.Validator.DefaultMessages.NegatedMessage) > 0 {
		return ms.Validator.DefaultMessages.Message
	}

	return fmt.Sprintf(ms.Validator.DefaultMessages.NegatedMessageFmt, ms.FieldName)
}

func GetValidator(key string) (*EmValidator, bool) {
	val, found := emKeyMap.Get(key)

	if found {
		return val.(*EmValidator), found
	}

	return nil, found
}

func SetCustomMessage(key string, msg string) error {
	v, exists := GetValidator(key)
	if !exists {
		return fmt.Errorf("Validator with key %s doesn't exist", key)
	}

	ms := MessageSet{}

	if strings.Index(key, "%s") > -1 {
		ms.MessageFmt = msg
	} else {
		ms.Message = msg
	}

	v.ValidatorCustomMessages = ms

	return nil
}

func SetCustomNegationMessage(key string, msg string) error {
	v, exists := GetValidator(key)
	if !exists {
		return fmt.Errorf("Validator with key %s doesn't exist", key)
	}

	ms := MessageSet{}

	if strings.Index(key, "%s") > -1 {
		ms.NegatedMessageFmt = msg
	} else {
		ms.NegatedMessage = msg
	}

	v.ValidatorCustomMessages = ms

	return nil
}
func SetCustomMessages(key string, ms MessageSet) error {
	v, exists := GetValidator(key)
	if !exists {
		return fmt.Errorf("Validator with key %s doesn't exist", key)
	}

	v.ValidatorCustomMessages = ms

	return nil
}

func ClearCustomMessages(validationKey string) {
	v, exists := GetValidator(validationKey)
	if !exists {
		return
	}

	v.ValidatorCustomMessages = MessageSet{}
}

func init() {
	emKeyMap.Clear()
	emKeyMap.Put("required", &EmValidator{Op: IsNonEmpty, CanValidateComplexTypes: true, DefaultMessages: MessageSet{Message: "{field} must not be empty"}})
	emKeyMap.Put("between", &EmValidator{Op: Between, DefaultMessages: MessageSet{Message: "{field} is out of range"}})
	emKeyMap.Put("matches", &EmValidator{OpString: StringMatches}) // can't use random regexes in
	emKeyMap.Put("title", &EmValidator{OpString: IsTitle})
	emKeyMap.Put("name", &EmValidator{OpString: IsName})
	emKeyMap.Put("phone", &EmValidator{OpString: IsPhone})
	emKeyMap.Put("skype", &EmValidator{OpString: IsSkype})
	emKeyMap.Put("email", &EmValidator{OpString: IsEmail})
	emKeyMap.Put("url", &EmValidator{OpString: IsURL})
	emKeyMap.Put("dialstring", &EmValidator{OpString: IsDialString})
	emKeyMap.Put("requrl", &EmValidator{OpString: IsRequestURL})
	emKeyMap.Put("requri", &EmValidator{OpString: IsRequestURI})
	emKeyMap.Put("alpha", &EmValidator{OpString: IsAlpha})
	emKeyMap.Put("utfletter", &EmValidator{OpString: IsUTFLetter})
	emKeyMap.Put("alphanum", &EmValidator{OpString: IsAlphanumeric})
	emKeyMap.Put("utfletternum", &EmValidator{OpString: IsUTFLetterNumeric})
	emKeyMap.Put("utfnumeric", &EmValidator{OpString: IsUTFNumeric})
	emKeyMap.Put("numeric", &EmValidator{OpString: IsNumeric})
	emKeyMap.Put("utfdigit", &EmValidator{OpString: IsUTFDigit})
	emKeyMap.Put("hexadecimal", &EmValidator{OpString: IsHexadecimal})
	emKeyMap.Put("hexcolor", &EmValidator{OpString: IsHexcolor})
	emKeyMap.Put("rgbcolor", &EmValidator{OpString: IsRGBcolor})
	emKeyMap.Put("lowercase", &EmValidator{OpString: IsLowerCase})
	emKeyMap.Put("uppercase", &EmValidator{OpString: IsUpperCase})
	emKeyMap.Put("int", &EmValidator{Op: IsInt})
	emKeyMap.Put("float", &EmValidator{Op: IsFloat})
	emKeyMap.Put("null", &EmValidator{Op: IsNull})
	emKeyMap.Put("uuid", &EmValidator{OpString: IsUUID})
	emKeyMap.Put("uuidv3", &EmValidator{OpString: IsUUIDv3})
	emKeyMap.Put("uuidv4", &EmValidator{OpString: IsUUIDv4})
	emKeyMap.Put("uuidv5", &EmValidator{OpString: IsUUIDv5})
	emKeyMap.Put("isoalpha2", &EmValidator{OpString: IsISO3166Alpha2})
	emKeyMap.Put("isoalpha3", &EmValidator{OpString: IsISO3166Alpha3})
	emKeyMap.Put("creditcard", &EmValidator{OpString: IsCreditCard})
	emKeyMap.Put("isbn10", &EmValidator{OpString: IsISBN10})
	emKeyMap.Put("isbn13", &EmValidator{OpString: IsISBN13})
	emKeyMap.Put("json", &EmValidator{OpString: IsJSON})
	emKeyMap.Put("multibyte", &EmValidator{OpString: IsMultibyte})
	emKeyMap.Put("ascii", &EmValidator{OpString: IsASCII})
	emKeyMap.Put("printableascii", &EmValidator{OpString: IsPrintableASCII})
	emKeyMap.Put("fullwidth", &EmValidator{OpString: IsFullWidth})
	emKeyMap.Put("halfwidth", &EmValidator{OpString: IsHalfWidth})
	emKeyMap.Put("variablewidth", &EmValidator{OpString: IsVariableWidth})
	emKeyMap.Put("base64", &EmValidator{OpString: IsBase64})
	emKeyMap.Put("datauri", &EmValidator{OpString: IsDataURI})
	emKeyMap.Put("ip", &EmValidator{OpString: IsIP})
	emKeyMap.Put("port", &EmValidator{Op: IsPort})
	emKeyMap.Put("ipv4", &EmValidator{OpString: IsIPv4})
	emKeyMap.Put("dns", &EmValidator{OpString: IsDNSName})
	emKeyMap.Put("host", &EmValidator{OpString: IsHost})
	emKeyMap.Put("mac", &EmValidator{OpString: IsMAC})
	emKeyMap.Put("latitude", &EmValidator{OpString: IsLatitude})
	emKeyMap.Put("longitude", &EmValidator{OpString: IsLongitude})
	emKeyMap.Put("ssn", &EmValidator{OpString: IsSSN})
	emKeyMap.Put("semver", &EmValidator{OpString: IsSemver})

	err := SetMessagesLocale(`en`)
	if err != nil {
		panic(err.Error())
	}
}

func SetMessagesLocale(locale string) error {

	var messageMap map[string]string

	switch locale {
	case `en`:
		messageMap = messages_en
	default:
		return fmt.Errorf("Locale %s is not implemented (we'd love it if you could help us fix that!)", locale)
	}

	// extract parts from the compound key format
	var keyRoot = func(fullKey string) string {
		idx := strings.Index(fullKey, ".")
		return fullKey[:idx]
	}

	var keyMsgType = func(fullKey string) string {
		idx := strings.Index(fullKey, ".")
		return fullKey[idx+1:]
	}

	// set messages
	for key, msg := range messageMap {

		v, exists := emKeyMap.Get(keyRoot(key))
		if !exists {
			continue
		}

		validator, _ := v.(*EmValidator)

		messageType := keyMsgType(key)

		switch messageType {
		case `messagefmt`:
			validator.DefaultMessages.MessageFmt = msg
		case `negatedmessagefmt`:
			validator.DefaultMessages.NegatedMessageFmt = msg
		case `message`:
			validator.DefaultMessages.Message = msg
		case `negatedmessage`:
			validator.DefaultMessages.NegatedMessage = msg
		default:
			return fmt.Errorf("%s is not a valid message type identifier", messageType)
		}
	}

	return nil
}

// ISO3166Entry stores country codes
type ISO3166Entry struct {
	EnglishShortName string
	FrenchShortName  string
	Alpha2Code       string
	Alpha3Code       string
	Numeric          string
}

//ISO3166List based on https://www.iso.org/obp/ui/#search/code/ Code Type "Officially Assigned Codes"
var ISO3166List = []ISO3166Entry{
	{"Afghanistan", "Afghanistan (l')", "AF", "AFG", "004"},
	{"Albania", "Albanie (l')", "AL", "ALB", "008"},
	{"Antarctica", "Antarctique (l')", "AQ", "ATA", "010"},
	{"Algeria", "Algérie (l')", "DZ", "DZA", "012"},
	{"American Samoa", "Samoa américaines (les)", "AS", "ASM", "016"},
	{"Andorra", "Andorre (l')", "AD", "AND", "020"},
	{"Angola", "Angola (l')", "AO", "AGO", "024"},
	{"Antigua and Barbuda", "Antigua-et-Barbuda", "AG", "ATG", "028"},
	{"Azerbaijan", "Azerbaïdjan (l')", "AZ", "AZE", "031"},
	{"Argentina", "Argentine (l')", "AR", "ARG", "032"},
	{"Australia", "Australie (l')", "AU", "AUS", "036"},
	{"Austria", "Autriche (l')", "AT", "AUT", "040"},
	{"Bahamas (the)", "Bahamas (les)", "BS", "BHS", "044"},
	{"Bahrain", "Bahreïn", "BH", "BHR", "048"},
	{"Bangladesh", "Bangladesh (le)", "BD", "BGD", "050"},
	{"Armenia", "Arménie (l')", "AM", "ARM", "051"},
	{"Barbados", "Barbade (la)", "BB", "BRB", "052"},
	{"Belgium", "Belgique (la)", "BE", "BEL", "056"},
	{"Bermuda", "Bermudes (les)", "BM", "BMU", "060"},
	{"Bhutan", "Bhoutan (le)", "BT", "BTN", "064"},
	{"Bolivia (Plurinational State of)", "Bolivie (État plurinational de)", "BO", "BOL", "068"},
	{"Bosnia and Herzegovina", "Bosnie-Herzégovine (la)", "BA", "BIH", "070"},
	{"Botswana", "Botswana (le)", "BW", "BWA", "072"},
	{"Bouvet Island", "Bouvet (l'Île)", "BV", "BVT", "074"},
	{"Brazil", "Brésil (le)", "BR", "BRA", "076"},
	{"Belize", "Belize (le)", "BZ", "BLZ", "084"},
	{"British Indian Ocean Territory (the)", "Indien (le Territoire britannique de l'océan)", "IO", "IOT", "086"},
	{"Solomon Islands", "Salomon (Îles)", "SB", "SLB", "090"},
	{"Virgin Islands (British)", "Vierges britanniques (les Îles)", "VG", "VGB", "092"},
	{"Brunei Darussalam", "Brunéi Darussalam (le)", "BN", "BRN", "096"},
	{"Bulgaria", "Bulgarie (la)", "BG", "BGR", "100"},
	{"Myanmar", "Myanmar (le)", "MM", "MMR", "104"},
	{"Burundi", "Burundi (le)", "BI", "BDI", "108"},
	{"Belarus", "Bélarus (le)", "BY", "BLR", "112"},
	{"Cambodia", "Cambodge (le)", "KH", "KHM", "116"},
	{"Cameroon", "Cameroun (le)", "CM", "CMR", "120"},
	{"Canada", "Canada (le)", "CA", "CAN", "124"},
	{"Cabo Verde", "Cabo Verde", "CV", "CPV", "132"},
	{"Cayman Islands (the)", "Caïmans (les Îles)", "KY", "CYM", "136"},
	{"Central African Republic (the)", "République centrafricaine (la)", "CF", "CAF", "140"},
	{"Sri Lanka", "Sri Lanka", "LK", "LKA", "144"},
	{"Chad", "Tchad (le)", "TD", "TCD", "148"},
	{"Chile", "Chili (le)", "CL", "CHL", "152"},
	{"China", "Chine (la)", "CN", "CHN", "156"},
	{"Taiwan (Province of China)", "Taïwan (Province de Chine)", "TW", "TWN", "158"},
	{"Christmas Island", "Christmas (l'Île)", "CX", "CXR", "162"},
	{"Cocos (Keeling) Islands (the)", "Cocos (les Îles)/ Keeling (les Îles)", "CC", "CCK", "166"},
	{"Colombia", "Colombie (la)", "CO", "COL", "170"},
	{"Comoros (the)", "Comores (les)", "KM", "COM", "174"},
	{"Mayotte", "Mayotte", "YT", "MYT", "175"},
	{"Congo (the)", "Congo (le)", "CG", "COG", "178"},
	{"Congo (the Democratic Republic of the)", "Congo (la République démocratique du)", "CD", "COD", "180"},
	{"Cook Islands (the)", "Cook (les Îles)", "CK", "COK", "184"},
	{"Costa Rica", "Costa Rica (le)", "CR", "CRI", "188"},
	{"Croatia", "Croatie (la)", "HR", "HRV", "191"},
	{"Cuba", "Cuba", "CU", "CUB", "192"},
	{"Cyprus", "Chypre", "CY", "CYP", "196"},
	{"Czech Republic (the)", "tchèque (la République)", "CZ", "CZE", "203"},
	{"Benin", "Bénin (le)", "BJ", "BEN", "204"},
	{"Denmark", "Danemark (le)", "DK", "DNK", "208"},
	{"Dominica", "Dominique (la)", "DM", "DMA", "212"},
	{"Dominican Republic (the)", "dominicaine (la République)", "DO", "DOM", "214"},
	{"Ecuador", "Équateur (l')", "EC", "ECU", "218"},
	{"El Salvador", "El Salvador", "SV", "SLV", "222"},
	{"Equatorial Guinea", "Guinée équatoriale (la)", "GQ", "GNQ", "226"},
	{"Ethiopia", "Éthiopie (l')", "ET", "ETH", "231"},
	{"Eritrea", "Érythrée (l')", "ER", "ERI", "232"},
	{"Estonia", "Estonie (l')", "EE", "EST", "233"},
	{"Faroe Islands (the)", "Féroé (les Îles)", "FO", "FRO", "234"},
	{"Falkland Islands (the) [Malvinas]", "Falkland (les Îles)/Malouines (les Îles)", "FK", "FLK", "238"},
	{"South Georgia and the South Sandwich Islands", "Géorgie du Sud-et-les Îles Sandwich du Sud (la)", "GS", "SGS", "239"},
	{"Fiji", "Fidji (les)", "FJ", "FJI", "242"},
	{"Finland", "Finlande (la)", "FI", "FIN", "246"},
	{"Åland Islands", "Åland(les Îles)", "AX", "ALA", "248"},
	{"France", "France (la)", "FR", "FRA", "250"},
	{"French Guiana", "Guyane française (la )", "GF", "GUF", "254"},
	{"French Polynesia", "Polynésie française (la)", "PF", "PYF", "258"},
	{"French Southern Territories (the)", "Terres australes françaises (les)", "TF", "ATF", "260"},
	{"Djibouti", "Djibouti", "DJ", "DJI", "262"},
	{"Gabon", "Gabon (le)", "GA", "GAB", "266"},
	{"Georgia", "Géorgie (la)", "GE", "GEO", "268"},
	{"Gambia (the)", "Gambie (la)", "GM", "GMB", "270"},
	{"Palestine, State of", "Palestine, État de", "PS", "PSE", "275"},
	{"Germany", "Allemagne (l')", "DE", "DEU", "276"},
	{"Ghana", "Ghana (le)", "GH", "GHA", "288"},
	{"Gibraltar", "Gibraltar", "GI", "GIB", "292"},
	{"Kiribati", "Kiribati", "KI", "KIR", "296"},
	{"Greece", "Grèce (la)", "GR", "GRC", "300"},
	{"Greenland", "Groenland (le)", "GL", "GRL", "304"},
	{"Grenada", "Grenade (la)", "GD", "GRD", "308"},
	{"Guadeloupe", "Guadeloupe (la)", "GP", "GLP", "312"},
	{"Guam", "Guam", "GU", "GUM", "316"},
	{"Guatemala", "Guatemala (le)", "GT", "GTM", "320"},
	{"Guinea", "Guinée (la)", "GN", "GIN", "324"},
	{"Guyana", "Guyana (le)", "GY", "GUY", "328"},
	{"Haiti", "Haïti", "HT", "HTI", "332"},
	{"Heard Island and McDonald Islands", "Heard-et-Îles MacDonald (l'Île)", "HM", "HMD", "334"},
	{"Holy See (the)", "Saint-Siège (le)", "VA", "VAT", "336"},
	{"Honduras", "Honduras (le)", "HN", "HND", "340"},
	{"Hong Kong", "Hong Kong", "HK", "HKG", "344"},
	{"Hungary", "Hongrie (la)", "HU", "HUN", "348"},
	{"Iceland", "Islande (l')", "IS", "ISL", "352"},
	{"India", "Inde (l')", "IN", "IND", "356"},
	{"Indonesia", "Indonésie (l')", "ID", "IDN", "360"},
	{"Iran (Islamic Republic of)", "Iran (République Islamique d')", "IR", "IRN", "364"},
	{"Iraq", "Iraq (l')", "IQ", "IRQ", "368"},
	{"Ireland", "Irlande (l')", "IE", "IRL", "372"},
	{"Israel", "Israël", "IL", "ISR", "376"},
	{"Italy", "Italie (l')", "IT", "ITA", "380"},
	{"Côte d'Ivoire", "Côte d'Ivoire (la)", "CI", "CIV", "384"},
	{"Jamaica", "Jamaïque (la)", "JM", "JAM", "388"},
	{"Japan", "Japon (le)", "JP", "JPN", "392"},
	{"Kazakhstan", "Kazakhstan (le)", "KZ", "KAZ", "398"},
	{"Jordan", "Jordanie (la)", "JO", "JOR", "400"},
	{"Kenya", "Kenya (le)", "KE", "KEN", "404"},
	{"Korea (the Democratic People's Republic of)", "Corée (la République populaire démocratique de)", "KP", "PRK", "408"},
	{"Korea (the Republic of)", "Corée (la République de)", "KR", "KOR", "410"},
	{"Kuwait", "Koweït (le)", "KW", "KWT", "414"},
	{"Kyrgyzstan", "Kirghizistan (le)", "KG", "KGZ", "417"},
	{"Lao People's Democratic Republic (the)", "Lao, République démocratique populaire", "LA", "LAO", "418"},
	{"Lebanon", "Liban (le)", "LB", "LBN", "422"},
	{"Lesotho", "Lesotho (le)", "LS", "LSO", "426"},
	{"Latvia", "Lettonie (la)", "LV", "LVA", "428"},
	{"Liberia", "Libéria (le)", "LR", "LBR", "430"},
	{"Libya", "Libye (la)", "LY", "LBY", "434"},
	{"Liechtenstein", "Liechtenstein (le)", "LI", "LIE", "438"},
	{"Lithuania", "Lituanie (la)", "LT", "LTU", "440"},
	{"Luxembourg", "Luxembourg (le)", "LU", "LUX", "442"},
	{"Macao", "Macao", "MO", "MAC", "446"},
	{"Madagascar", "Madagascar", "MG", "MDG", "450"},
	{"Malawi", "Malawi (le)", "MW", "MWI", "454"},
	{"Malaysia", "Malaisie (la)", "MY", "MYS", "458"},
	{"Maldives", "Maldives (les)", "MV", "MDV", "462"},
	{"Mali", "Mali (le)", "ML", "MLI", "466"},
	{"Malta", "Malte", "MT", "MLT", "470"},
	{"Martinique", "Martinique (la)", "MQ", "MTQ", "474"},
	{"Mauritania", "Mauritanie (la)", "MR", "MRT", "478"},
	{"Mauritius", "Maurice", "MU", "MUS", "480"},
	{"Mexico", "Mexique (le)", "MX", "MEX", "484"},
	{"Monaco", "Monaco", "MC", "MCO", "492"},
	{"Mongolia", "Mongolie (la)", "MN", "MNG", "496"},
	{"Moldova (the Republic of)", "Moldova , République de", "MD", "MDA", "498"},
	{"Montenegro", "Monténégro (le)", "ME", "MNE", "499"},
	{"Montserrat", "Montserrat", "MS", "MSR", "500"},
	{"Morocco", "Maroc (le)", "MA", "MAR", "504"},
	{"Mozambique", "Mozambique (le)", "MZ", "MOZ", "508"},
	{"Oman", "Oman", "OM", "OMN", "512"},
	{"Namibia", "Namibie (la)", "NA", "NAM", "516"},
	{"Nauru", "Nauru", "NR", "NRU", "520"},
	{"Nepal", "Népal (le)", "NP", "NPL", "524"},
	{"Netherlands (the)", "Pays-Bas (les)", "NL", "NLD", "528"},
	{"Curaçao", "Curaçao", "CW", "CUW", "531"},
	{"Aruba", "Aruba", "AW", "ABW", "533"},
	{"Sint Maarten (Dutch part)", "Saint-Martin (partie néerlandaise)", "SX", "SXM", "534"},
	{"Bonaire, Sint Eustatius and Saba", "Bonaire, Saint-Eustache et Saba", "BQ", "BES", "535"},
	{"New Caledonia", "Nouvelle-Calédonie (la)", "NC", "NCL", "540"},
	{"Vanuatu", "Vanuatu (le)", "VU", "VUT", "548"},
	{"New Zealand", "Nouvelle-Zélande (la)", "NZ", "NZL", "554"},
	{"Nicaragua", "Nicaragua (le)", "NI", "NIC", "558"},
	{"Niger (the)", "Niger (le)", "NE", "NER", "562"},
	{"Nigeria", "Nigéria (le)", "NG", "NGA", "566"},
	{"Niue", "Niue", "NU", "NIU", "570"},
	{"Norfolk Island", "Norfolk (l'Île)", "NF", "NFK", "574"},
	{"Norway", "Norvège (la)", "NO", "NOR", "578"},
	{"Northern Mariana Islands (the)", "Mariannes du Nord (les Îles)", "MP", "MNP", "580"},
	{"United States Minor Outlying Islands (the)", "Îles mineures éloignées des États-Unis (les)", "UM", "UMI", "581"},
	{"Micronesia (Federated States of)", "Micronésie (États fédérés de)", "FM", "FSM", "583"},
	{"Marshall Islands (the)", "Marshall (Îles)", "MH", "MHL", "584"},
	{"Palau", "Palaos (les)", "PW", "PLW", "585"},
	{"Pakistan", "Pakistan (le)", "PK", "PAK", "586"},
	{"Panama", "Panama (le)", "PA", "PAN", "591"},
	{"Papua New Guinea", "Papouasie-Nouvelle-Guinée (la)", "PG", "PNG", "598"},
	{"Paraguay", "Paraguay (le)", "PY", "PRY", "600"},
	{"Peru", "Pérou (le)", "PE", "PER", "604"},
	{"Philippines (the)", "Philippines (les)", "PH", "PHL", "608"},
	{"Pitcairn", "Pitcairn", "PN", "PCN", "612"},
	{"Poland", "Pologne (la)", "PL", "POL", "616"},
	{"Portugal", "Portugal (le)", "PT", "PRT", "620"},
	{"Guinea-Bissau", "Guinée-Bissau (la)", "GW", "GNB", "624"},
	{"Timor-Leste", "Timor-Leste (le)", "TL", "TLS", "626"},
	{"Puerto Rico", "Porto Rico", "PR", "PRI", "630"},
	{"Qatar", "Qatar (le)", "QA", "QAT", "634"},
	{"Réunion", "Réunion (La)", "RE", "REU", "638"},
	{"Romania", "Roumanie (la)", "RO", "ROU", "642"},
	{"Russian Federation (the)", "Russie (la Fédération de)", "RU", "RUS", "643"},
	{"Rwanda", "Rwanda (le)", "RW", "RWA", "646"},
	{"Saint Barthélemy", "Saint-Barthélemy", "BL", "BLM", "652"},
	{"Saint Helena, Ascension and Tristan da Cunha", "Sainte-Hélène, Ascension et Tristan da Cunha", "SH", "SHN", "654"},
	{"Saint Kitts and Nevis", "Saint-Kitts-et-Nevis", "KN", "KNA", "659"},
	{"Anguilla", "Anguilla", "AI", "AIA", "660"},
	{"Saint Lucia", "Sainte-Lucie", "LC", "LCA", "662"},
	{"Saint Martin (French part)", "Saint-Martin (partie française)", "MF", "MAF", "663"},
	{"Saint Pierre and Miquelon", "Saint-Pierre-et-Miquelon", "PM", "SPM", "666"},
	{"Saint Vincent and the Grenadines", "Saint-Vincent-et-les Grenadines", "VC", "VCT", "670"},
	{"San Marino", "Saint-Marin", "SM", "SMR", "674"},
	{"Sao Tome and Principe", "Sao Tomé-et-Principe", "ST", "STP", "678"},
	{"Saudi Arabia", "Arabie saoudite (l')", "SA", "SAU", "682"},
	{"Senegal", "Sénégal (le)", "SN", "SEN", "686"},
	{"Serbia", "Serbie (la)", "RS", "SRB", "688"},
	{"Seychelles", "Seychelles (les)", "SC", "SYC", "690"},
	{"Sierra Leone", "Sierra Leone (la)", "SL", "SLE", "694"},
	{"Singapore", "Singapour", "SG", "SGP", "702"},
	{"Slovakia", "Slovaquie (la)", "SK", "SVK", "703"},
	{"Viet Nam", "Viet Nam (le)", "VN", "VNM", "704"},
	{"Slovenia", "Slovénie (la)", "SI", "SVN", "705"},
	{"Somalia", "Somalie (la)", "SO", "SOM", "706"},
	{"South Africa", "Afrique du Sud (l')", "ZA", "ZAF", "710"},
	{"Zimbabwe", "Zimbabwe (le)", "ZW", "ZWE", "716"},
	{"Spain", "Espagne (l')", "ES", "ESP", "724"},
	{"South Sudan", "Soudan du Sud (le)", "SS", "SSD", "728"},
	{"Sudan (the)", "Soudan (le)", "SD", "SDN", "729"},
	{"Western Sahara*", "Sahara occidental (le)*", "EH", "ESH", "732"},
	{"Suriname", "Suriname (le)", "SR", "SUR", "740"},
	{"Svalbard and Jan Mayen", "Svalbard et l'Île Jan Mayen (le)", "SJ", "SJM", "744"},
	{"Swaziland", "Swaziland (le)", "SZ", "SWZ", "748"},
	{"Sweden", "Suède (la)", "SE", "SWE", "752"},
	{"Switzerland", "Suisse (la)", "CH", "CHE", "756"},
	{"Syrian Arab Republic", "République arabe syrienne (la)", "SY", "SYR", "760"},
	{"Tajikistan", "Tadjikistan (le)", "TJ", "TJK", "762"},
	{"Thailand", "Thaïlande (la)", "TH", "THA", "764"},
	{"Togo", "Togo (le)", "TG", "TGO", "768"},
	{"Tokelau", "Tokelau (les)", "TK", "TKL", "772"},
	{"Tonga", "Tonga (les)", "TO", "TON", "776"},
	{"Trinidad and Tobago", "Trinité-et-Tobago (la)", "TT", "TTO", "780"},
	{"United Arab Emirates (the)", "Émirats arabes unis (les)", "AE", "ARE", "784"},
	{"Tunisia", "Tunisie (la)", "TN", "TUN", "788"},
	{"Turkey", "Turquie (la)", "TR", "TUR", "792"},
	{"Turkmenistan", "Turkménistan (le)", "TM", "TKM", "795"},
	{"Turks and Caicos Islands (the)", "Turks-et-Caïcos (les Îles)", "TC", "TCA", "796"},
	{"Tuvalu", "Tuvalu (les)", "TV", "TUV", "798"},
	{"Uganda", "Ouganda (l')", "UG", "UGA", "800"},
	{"Ukraine", "Ukraine (l')", "UA", "UKR", "804"},
	{"Macedonia (the former Yugoslav Republic of)", "Macédoine (l'ex‑République yougoslave de)", "MK", "MKD", "807"},
	{"Egypt", "Égypte (l')", "EG", "EGY", "818"},
	{"United Kingdom of Great Britain and Northern Ireland (the)", "Royaume-Uni de Grande-Bretagne et d'Irlande du Nord (le)", "GB", "GBR", "826"},
	{"Guernsey", "Guernesey", "GG", "GGY", "831"},
	{"Jersey", "Jersey", "JE", "JEY", "832"},
	{"Isle of Man", "Île de Man", "IM", "IMN", "833"},
	{"Tanzania, United Republic of", "Tanzanie, République-Unie de", "TZ", "TZA", "834"},
	{"United States of America (the)", "États-Unis d'Amérique (les)", "US", "USA", "840"},
	{"Virgin Islands (U.S.)", "Vierges des États-Unis (les Îles)", "VI", "VIR", "850"},
	{"Burkina Faso", "Burkina Faso (le)", "BF", "BFA", "854"},
	{"Uruguay", "Uruguay (l')", "UY", "URY", "858"},
	{"Uzbekistan", "Ouzbékistan (l')", "UZ", "UZB", "860"},
	{"Venezuela (Bolivarian Republic of)", "Venezuela (République bolivarienne du)", "VE", "VEN", "862"},
	{"Wallis and Futuna", "Wallis-et-Futuna", "WF", "WLF", "876"},
	{"Samoa", "Samoa (le)", "WS", "WSM", "882"},
	{"Yemen", "Yémen (le)", "YE", "YEM", "887"},
	{"Zambia", "Zambie (la)", "ZM", "ZMB", "894"},
}
