/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

func init() {
	Translations.Register(LanguageRU, map[string]string{
		MessageRequired:                "Значение не должно быть пустым.",
		MessageStringType:              "Значение должно быть строкой.",
		MessageTooShort:                "Значение должно содержать минимум {min} символов.",
		MessageTooLong:                 "Значение должно содержать максимум {max} символов.",
		MessageCompareEqual:            "Значение должно быть равно «{targetValueOrAttribute}».",
		MessageCompareNotEqual:         "Значение не должно быть равно «{targetValueOrAttribute}».",
		MessageCompareGreaterThan:      "Значение должно быть больше «{targetValueOrAttribute}».",
		MessageCompareGreaterThanEqual: "Значение должно быть больше или равно «{targetValueOrAttribute}».",
		MessageCompareLessThan:         "Значение должно быть меньше «{targetValueOrAttribute}».",
		MessageCompareLessThanEqual:    "Значение должно быть меньше или равно «{targetValueOrAttribute}».",
		MessageNotNumber:               "Значение должно быть числом.",
		MessageNumberTooBig:            "Значение не должно превышать {max}.",
		MessageNumberTooSmall:          "Значение не должно быть меньше {min}.",
		MessageNotNumeric:              "Значение должно быть числовым.",
		MessageInvalidURL:              "Значение не является допустимым URL.",
		MessageInvalidDeepLink:         "Значение не является допустимым deep link URL.",
		MessageInvalidIP:               "Значение должно быть допустимым IP-адресом.",
		MessageIPv4NotAllowed:          "Значение не должно быть IPv4-адресом.",
		MessageIPv6NotAllowed:          "Значение не должно быть IPv6-адресом.",
		MessageInvalidUUID:             "Недопустимый формат UUID.",
		MessageInvalidUUIDVersion:      "Версия UUID должна быть равна {version}.",
		MessageInvalidJSON:             "Значение должно быть допустимым JSON.",
		MessageInvalidMAC:              "Значение должно быть допустимым MAC-адресом.",
		MessageValueInvalid:            "Значение недопустимо.",
		MessageInvalidEmail:            "Значение не является допустимым email-адресом.",
		MessageInvalid:                 "Значение недопустимо.",
		MessageTimeFormat:              "Формат времени должен соответствовать {format}.",
		MessageTimeTooBig:              "Время не должно превышать {max}.",
		MessageTimeTooSmall:            "Время не должно быть меньше {min}.",
		MessageInvalidImageMeta:        "Значение должно быть структурой ImageMetaData.",
		MessageInvalidMimeType:         "MimeType должен соответствовать одному из [{mimeTypes}].",
		MessageFileTooLarge:            "Размер файла не должен превышать {maxFileSizeBytes} байт.",
		MessageInRangeInvalid:          "Значение недопустимо.",
		MessageEachIncorrectInput:      "Значение должно быть массивом.",
		MessageUniqueValues:            "Список значений должен быть уникальным.",
		MessageInvalidOGRN:             "Значение не является допустимым ОГРН.",
		MessageInvalidOGRNIP:           "Значение не является допустимым ОГРНИП.",
		MessageInvalidOGRNLength:       "Значение должно содержать 13 или 15 символов.",
		MessageInvalidSQL:              "Значение не является допустимым SQL.",
		MessageInvalidMSISDN:           "Недопустимый формат MSISDN.",
		MessageInvalidHumanText:        "Значение должно быть обычным текстом.",
		MessageNestedNotStruct:         "Значение должно быть структурой. Получен %T.",
	})
}
