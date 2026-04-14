/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

const (
	// Required
	MessageRequired = "Value cannot be blank."

	// StringLength
	MessageStringType = "This value must be a string."
	MessageTooShort   = "This value should contain at least {min}."
	MessageTooLong    = "This value should contain at most {max}."

	// Compare
	MessageCompareEqual            = "Value must be equal to '{targetValueOrAttribute}'."
	MessageCompareNotEqual         = "Value must not be equal to '{targetValueOrAttribute}'."
	MessageCompareGreaterThan      = "Value must be greater than '{targetValueOrAttribute}'"
	MessageCompareGreaterThanEqual = "Value must be greater than or equal to '{targetValueOrAttribute}'"
	MessageCompareLessThan         = "Value must be less than '{targetValueOrAttribute}'"
	MessageCompareLessThanEqual    = "Value must be less than or equal to '{targetValueOrAttribute}'"

	// Number
	MessageNotNumber      = "Value must be a number."
	MessageNumberTooBig   = "Value must be no greater than {max}."
	MessageNumberTooSmall = "Value must be no less than {min}."

	// Numeric
	MessageNotNumeric = "Value must be a numeric."

	// URL
	MessageInvalidURL      = "This value is not a valid URL."
	MessageInvalidDeepLink = "This value is not a valid deep link url."

	// IP
	MessageInvalidIP      = "Must be a valid IP address."
	MessageIPv4NotAllowed = "Must not be an IPv4 address."
	MessageIPv6NotAllowed = "Must not be an IPv6 address."

	// UUID
	MessageInvalidUUID        = "Invalid UUID format."
	MessageInvalidUUIDVersion = "UUID version must be equal to {version}."

	// JSON
	MessageInvalidJSON = "Must be a valid JSON"

	// MAC
	MessageInvalidMAC = "Must be a valid MAC address."

	// MatchRegularExpression
	MessageValueInvalid = "Value is invalid."

	// Generic (used by Time, Each)
	MessageInvalid = "Value is invalid"

	// Email
	MessageInvalidEmail = "Email is not a valid email."

	// Time
	MessageTimeFormat   = "Format of the time value must be equal {format}"
	MessageTimeTooBig   = "Time must be no greater than {max}."
	MessageTimeTooSmall = "Time must be no less than {min}."

	// ImageMeta
	MessageInvalidImageMeta = "This value must be a ImageMetaData struct."
	MessageInvalidMimeType  = "MimeType must be a equal one of [{mimeTypes}]."
	MessageFileTooLarge     = "File size must be greater than {maxFileSizeBytes} bytes."

	// InRange
	MessageInRangeInvalid = "This value is invalid"

	// Each
	MessageEachIncorrectInput = "Value must be array"

	// UniqueValues
	MessageUniqueValues = "The list of values must be unique."

	// OGRN
	MessageInvalidOGRN       = "This value is not a valid OGRN."
	MessageInvalidOGRNIP     = "This value is not a valid OGRNIP."
	MessageInvalidOGRNLength = "This value should contain either 13 or 15 characters."

	// SQL
	MessageInvalidSQL = "Value is invalid sql."

	// MSISDN
	MessageInvalidMSISDN = "MSISDN format is invalid."

	// HumanText
	MessageInvalidHumanText = "This value must be a normal text."

	// Nested
	MessageNestedNotStruct = "value should be a struct. %T given."
)
