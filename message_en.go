/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

func init() {
	Translations.Register(LanguageEN, map[string]string{
		MessageRequired:                MessageRequired,
		MessageStringType:              MessageStringType,
		MessageTooShort:                MessageTooShort,
		MessageTooLong:                 MessageTooLong,
		MessageCompareEqual:            MessageCompareEqual,
		MessageCompareNotEqual:         MessageCompareNotEqual,
		MessageCompareGreaterThan:      MessageCompareGreaterThan,
		MessageCompareGreaterThanEqual: MessageCompareGreaterThanEqual,
		MessageCompareLessThan:         MessageCompareLessThan,
		MessageCompareLessThanEqual:    MessageCompareLessThanEqual,
		MessageNotNumber:               MessageNotNumber,
		MessageNumberTooBig:            MessageNumberTooBig,
		MessageNumberTooSmall:          MessageNumberTooSmall,
		MessageNotNumeric:              MessageNotNumeric,
		MessageInvalidURL:              MessageInvalidURL,
		MessageInvalidDeepLink:         MessageInvalidDeepLink,
		MessageInvalidIP:               MessageInvalidIP,
		MessageIPv4NotAllowed:          MessageIPv4NotAllowed,
		MessageIPv6NotAllowed:          MessageIPv6NotAllowed,
		MessageInvalidUUID:             MessageInvalidUUID,
		MessageInvalidUUIDVersion:      MessageInvalidUUIDVersion,
		MessageInvalidJSON:             MessageInvalidJSON,
		MessageInvalidMAC:              MessageInvalidMAC,
		MessageValueInvalid:            MessageValueInvalid,
		MessageInvalidEmail:            MessageInvalidEmail,
		MessageInvalid:                 MessageInvalid,
		MessageTimeFormat:              MessageTimeFormat,
		MessageTimeTooBig:              MessageTimeTooBig,
		MessageTimeTooSmall:            MessageTimeTooSmall,
		MessageInvalidImageMeta:        MessageInvalidImageMeta,
		MessageInvalidMimeType:         MessageInvalidMimeType,
		MessageFileTooLarge:            MessageFileTooLarge,
		MessageInRangeInvalid:          MessageInRangeInvalid,
		MessageEachIncorrectInput:      MessageEachIncorrectInput,
		MessageUniqueValues:            MessageUniqueValues,
		MessageInvalidOGRN:             MessageInvalidOGRN,
		MessageInvalidOGRNIP:           MessageInvalidOGRNIP,
		MessageInvalidOGRNLength:       MessageInvalidOGRNLength,
		MessageInvalidSQL:              MessageInvalidSQL,
		MessageInvalidMSISDN:           MessageInvalidMSISDN,
		MessageInvalidHumanText:        MessageInvalidHumanText,
		MessageNestedNotStruct:         MessageNestedNotStruct,
	})
}
