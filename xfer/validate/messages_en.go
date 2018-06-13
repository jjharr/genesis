package validate

var messages_en = map[string]string{
	`email.messagefmt`:        `%s must be a valid address`,
	`email.negatedmessagefmt`: `%s must not be an email address`,

	`url.messagefmt`:        `%s must be a full URL`,
	`url.negatedmessagefmt`: `%s must not be a URL`,

	`dialstring.messagefmt`:        `%s must be a port, IP address, or DNS address`,
	`dialstring.negatedmessagefmt`: `%s must not be a port, IP address, or DNS address`,

	`requrl.messagefmt`:        `%s must be a full URL`,
	`requrl.negatedmessagefmt`: `%s must not be a full URL`,

	`requri.messagefmt`:        `%s must be a full URI`,
	`requri.negatedmessagefmt`: `%s must not be a full URI`,

	`alpha.messagefmt`:        `%s must only contain letters`,
	`alpha.negatedmessagefmt`: `%s must not contain letters`,

	`utfletter.messagefmt`:        `%s must only contain letters`,
	`utfletter.negatedmessagefmt`: `%s must not contain letters`,

	`alphanum.messagefmt`:        `%s must only contain letters and numbers`,
	`alphanum.negatedmessagefmt`: `%s must not contain letters or numbers`,

	`utfletternum.messagefmt`:        `%s must only contain letters and numbers`,
	`utfletternum.negatedmessagefmt`: `%s must not contain letters or numbers`,

	`utfnumeric.messagefmt`:        `%s must only contain numbers`,
	`utfnumeric.negatedmessagefmt`: `%s must not contain numbers`,

	`utfdigit.messagefmt`:        `%s must only contain numbers`,
	`utfdigit.negatedmessagefmt`: `%s must not contain numbers`,

	`numeric.messagefmt`:        `%s must only contain numbers`,
	`numeric.negatedmessagefmt`: `%s must not contain numbers`,

	`hexidecimal.messagefmt`:        `%s must be a hex value`,
	`hexidecimal.negatedmessagefmt`: `%s must not be a hex value`,

	`hexcolor.messagefmt`:        `%s must be a hex color`,
	`hexcolor.negatedmessagefmt`: `%s must not be a hex color`,

	`rgbcolor.messagefmt`:        `%s must be an RGB color`,
	`rgbcolor.negatedmessagefmt`: `%s must not be an RGB color`,

	`lowercase.messagefmt`:        `%s must be all lower case`,
	`lowercase.negatedmessagefmt`: `%s must not have lower case letters`,

	`uppercase.messagefmt`:        `%s must all be all upper case`,
	`uppercase.negatedmessagefmt`: `%s must not have upper case letters`,

	`float.messagefmt`:        `%s must be a floating point (decimal) number`,
	`float.negatedmessagefmt`: `%s must not be a floating point (decimal) number`,

	`null.messagefmt`:        `%s must be null`,
	`null.negatedmessagefmt`: `%s must not be null`,

	`uuid.messagefmt`:        `%s must be a UUID`,
	`uuid.negatedmessagefmt`: `%s must not be a UUID`,

	`uuid3.messagefmt`:        `%s must be a UUID (v3)`,
	`uuid3.negatedmessagefmt`: `%s must not be a UUID (v3)`,

	`uuid4.messagefmt`:        `%s must be a UUID (v4)`,
	`uuid4.negatedmessagefmt`: `%s must not be a UUID (v4)`,

	`uuid5.messagefmt`:        `%s must be a UUID (v5)`,
	`uuid5.negatedmessagefmt`: `%s must not be a UUID (v5)`,

	`creditcard.messagefmt`:        `%s must be a valid credit card number`,
	`creditcard.negatedmessagefmt`: `%s must not be a credit card number`,

	`json.messagefmt`:        `%s must be a valid JSON`,
	`json.negatedmessagefmt`: `%s must not be JSON`,

	`multibyte.messagefmt`:        `%s must be multibyte text`,
	`multibyte.negatedmessagefmt`: `%s must not be multibyte text`,

	`ascii.messagefmt`:        `%s must be ascii text`,
	`ascii.negatedmessagefmt`: `%s must not be ascii text`,

	`printableascii.messagefmt`:        `%s must be printable ascii text`,
	`printableascii.negatedmessagefmt`: `%s must not be printable ascii text`,

	`fullwidth.messagefmt`:        `%s must contain fullwidth UTF characters`,
	`fullwidth.negatedmessagefmt`: `%s must not contain fullwidth UTF characters`,

	`halfwidth.messagefmt`:        `%s must contain halfwidth UTF characters`,
	`halfwidth.negatedmessagefmt`: `%s must not contain halfwidth UTF characters`,

	`variablewidth.messagefmt`:        `%s must contain variable-width UTF characters`,
	`variablewidth.negatedmessagefmt`: `%s must not contain variable-width UTF characters`,

	`base64.messagefmt`:        `%s must be valid Base64`,
	`base64.negatedmessagefmt`: `%s must not be Base64`,

	`ip.messagefmt`:        `%s must be a valid IP address`,
	`ip.negatedmessagefmt`: `%s must not be an IP address`,

	`port.messagefmt`:        `%s must be a valid port`,
	`port.negatedmessagefmt`: `%s must not be a port number`,

	`ipv4.messagefmt`:        `%s must be a valid IPv4 address`,
	`ipv4.negatedmessagefmt`: `%s must not be an IPv4 address`,

	`dns.messagefmt`:        `%s must be a valid DNS name`,
	`dns.negatedmessagefmt`: `%s must not be a DNS name`,

	`host.messagefmt`:        `%s must be a valid host name`,
	`host.negatedmessagefmt`: `%s must not be a host name`,

	`mac.messagefmt`:        `%s must be a valid MAC address`,
	`mac.negatedmessagefmt`: `%s must not be a MAC address`,

	`latitude.messagefmt`:        `%s must be a valid latitude`,
	`latitude.negatedmessagefmt`: `%s must not be a latitude`,

	`longitude.messagefmt`:        `%s must be a valid longitude`,
	`longitude.negatedmessagefmt`: `%s must not be a longitude`,

	`ssn.messagefmt`:        `%s must be a valid SSN`,
	`ssn.negatedmessagefmt`: `%s must not be a SSN`,

	`semver.messagefmt`:        `%s must be a valid semantic version`,
	`semver.negatedmessagefmt`: `%s must not be a semantic version`,
}
