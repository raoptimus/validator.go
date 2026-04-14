# validator.go

[![CI](https://github.com/raoptimus/validator.go/actions/workflows/ci.yml/badge.svg)](https://github.com/raoptimus/validator.go/actions/workflows/ci.yml)
[![Coverage](https://github.com/raoptimus/validator.go/wiki/coverage.svg)](https://raw.githack.com/wiki/raoptimus/validator.go/coverage.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/raoptimus/validator.go.svg)](https://pkg.go.dev/github.com/raoptimus/validator.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/raoptimus/validator.go)](https://goreportcard.com/report/github.com/raoptimus/validator.go)
[![GitHub Release](https://img.shields.io/github/release/raoptimus/validator.go.svg)](https://github.com/raoptimus/validator.go/releases)
[![License](https://img.shields.io/github/license/raoptimus/validator.go.svg)](https://github.com/raoptimus/validator.go/blob/main/LICENSE)


A Go validation library with composable rules, nested struct support, and detailed error reporting.

## Installation

```bash
go get github.com/raoptimus/validator.go/v2
```

## Usage

### Validate a single value

```go
import validator "github.com/raoptimus/validator.go/v2"

ctx := context.Background()

err := validator.ValidateValue(ctx, "user@example.com",
    validator.NewRequired(),
    validator.NewEmail(),
)
// err == nil
```

### Validate a struct

Field names in `RuleSet` must match struct field names. Error paths use `json` tags as aliases.

```go
type SignUpRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
    Age   int    `json:"age"`
}

req := SignUpRequest{
    Email: "bad-email",
    Name:  "",
    Age:   200,
}

err := validator.Validate(ctx, &req, validator.RuleSet{
    "Email": {validator.NewRequired(), validator.NewEmail()},
    "Name":  {validator.NewRequired(), validator.NewStringLength(1, 100)},
    "Age":   {validator.NewRequired(), validator.NewNumber(1, 150)},
})

// err.Error() == "email: Email is not a valid email. name: Value cannot be blank. age: Value must be no greater than 150."
```

### Validate a map

```go
data := map[string]any{
    "count": 0,
}

err := validator.Validate(ctx, data, validator.RuleSet{
    "count": {
        validator.NewRequired(),
        validator.NewNumber(1, 100),
    },
})
```

### Nested structs

```go
type Address struct {
    City   string `json:"city"`
    Street string `json:"street"`
}

type User struct {
    Name    string  `json:"name"`
    Address Address `json:"address"`
}

err := validator.Validate(ctx, &User{Name: "Alice"}, validator.RuleSet{
    "Name": {validator.NewRequired()},
    "Address": {
        validator.NewNested(validator.RuleSet{
            "City":   {validator.NewRequired()},
            "Street": {validator.NewRequired(), validator.NewStringLength(1, 255)},
        }),
    },
})

// Errors have dotted paths: "address.city", "address.street"
```

### Validate each element in a slice

```go
err := validator.ValidateValue(ctx, []string{"ok", "", "fine"},
    validator.NewEach(validator.NewStringLength(1, 255)),
)

// Error path: "1" (index of the empty string)
```

### Conditional validation

```go
validator.NewEmail().When(func(ctx context.Context, value any) bool {
    return value != nil && value.(string) != ""
})
```

### Skip on empty / skip on error

```go
// Skip email format check if the value is empty
validator.NewEmail().SkipOnEmpty()

// Skip this rule if a previous rule already failed
validator.NewStringLength(1, 100).SkipOnError()
```

### OR rule

At least one of the rules must pass:

```go
validator.NewOR("Must be a valid email or URL",
    validator.NewEmail(),
    validator.NewURL(),
)
```

### Custom callback

```go
validator.NewCallback(func(ctx context.Context, value string) error {
    if strings.Contains(value, "forbidden") {
        return validator.NewValidationError("Value contains forbidden word.")
    }
    return nil
})
```

### Error handling

```go
err := validator.Validate(ctx, &req, rules)
if err != nil {
    var result validator.Result
    if errors.As(err, &result) {
        // Map of field path -> error messages
        errorMap := result.ErrorMessagesIndexedByPath()
        // e.g. {"email": ["Email is not a valid email."], "age": ["Value must be no greater than 150."]}
    }
}
```

### Multilanguage support

Validation messages are translated automatically based on the language in context.
Built-in languages: English (`en`, default) and Russian (`ru`).

```go
// Set language in context (e.g. from HTTP middleware)
ctx := validator.WithLanguage(context.Background(), validator.LanguageRU)

err := validator.ValidateValue(ctx, "", validator.NewRequired())
// err.Error() == "Значение не должно быть пустым."
```

#### Register a custom language

```go
func init() {
    validator.Translations.Register(validator.Language("es"), map[string]string{
        validator.MessageRequired: "El valor no puede estar vacío.",
        validator.MessageTooShort: "El valor debe contener al menos {min} caracteres.",
        // ...
    })
}
```

#### Check for missing translations

```go
missing := validator.Translations.Missing(validator.Language("es"))
// returns a list of message IDs without translation for the given language
```

#### Use DummyTranslator (no translation, placeholder replacement only)

```go
validator.SetTranslator(&validator.DummyTranslator{})
```

## Available Rules

| Rule | Description |
|------|-------------|
| `Required` | Value must not be empty |
| `Number` | Integer within min/max range |
| `StringLength` | UTF-8 string length within min/max |
| `Email` | Valid email format |
| `URL` | Valid URL format |
| `IP` | Valid IP address |
| `UUID` | Valid UUID format |
| `MAC` | Valid MAC address |
| `MSISDN` | Valid mobile phone number |
| `OGRN` | Valid Russian OGRN number |
| `JSON` | Valid JSON structure |
| `SQL` | Valid SQL statement |
| `MatchRegularExpression` | Matches a custom regex |
| `In` | Value is in a predefined set |
| `InRange` | Numeric value within a range |
| `Compare` | Compare against another value (==, !=, <, >, <=, >=) |
| `Nested` | Validate nested struct fields |
| `Each` | Validate each element in a slice/array |
| `OR` | At least one rule must pass |
| `Callback` | Custom validation function |
| `UniqueValues` | All elements in a slice are unique |
| `NumericString` | String is a valid numeric representation |
| `HumanText` | Valid human-readable text |
| `ImageMeta` | Image metadata validation |
| `Time` | Time format and range validation |

## License

[BSD 3-Clause](LICENSE)
