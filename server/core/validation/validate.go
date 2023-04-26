package validation

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/gofiber/fiber/v2"
)

type ValidationRuleFunc func(args []string, target any, value string) bool

var Rules map[string]ValidationRuleFunc

func init() {
	Rules = map[string]ValidationRuleFunc{
		// String rules
		"min": func(args []string, target any, value string) bool {
			min, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				log.Fatalln(err)
			}
			if len(value) < int(min) {
				return false
			}
			return true
		},
		"max": func(args []string, target any, value string) bool {
			max, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				log.Fatalln(err)
			}
			if len(value) > int(max) {
				return false
			}
			return true
		},

		// Type rules
		"boolean": func(args []string, target any, value string) bool {
			return value == "true" || value == "false"
		},
		"integer": func(args []string, target any, value string) bool {
			if _, err := strconv.ParseInt(value, 10, 64); err != nil {
				return false
			}
			return true
		},
		"float": func(args []string, target any, value string) bool {
			if _, err := strconv.ParseFloat(value, 64); err != nil {
				return false
			}
			return true
		},
		"date": func(args []string, target any, value string) bool {
			re := regexp.MustCompile(`^20\d\d-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\dZ$`)
			return re.MatchString(value)
		},
		"uuid": func(args []string, target any, value string) bool {
			return uuid.IsValid(value)
		},
		"email": func(args []string, target any, value string) bool {
			re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
			return re.MatchString(value)
		},
		"enum": func(args []string, target any, value string) bool {
			for _, arg := range args {
				if arg == value {
					return true
				}
			}
			return false
		},

		// Database rules
		"exists": func(args []string, target any, value string) bool {
			selectQuery := "SELECT COUNT(`id`) FROM `" + args[0] + "` WHERE `" + args[1] + "` = "
			if strings.HasSuffix(args[1], "id") {
				selectQuery += "UUID_TO_BIN(?)"
			} else {
				selectQuery += "?"
			}
			query := database.Query(selectQuery, value)
			defer query.Close()
			query.Next()
			var count int64
			_ = query.Scan(&count)
			return count == 1
		},
		"unique": func(args []string, target any, value string) bool {
			if len(args) > 2 {
				if reflect.ValueOf(target).Elem().FieldByName(args[2]).Interface().(string) == value {
					return true
				}
			}

			selectQuery := "SELECT COUNT(`id`) FROM `" + args[0] + "` WHERE `" + args[1] + "` = "
			if strings.HasSuffix(args[1], "id") {
				selectQuery += "UUID_TO_BIN(?)"
			} else {
				selectQuery += "?"
			}
			query := database.Query(selectQuery, value)
			defer query.Close()
			query.Next()
			var count int64
			_ = query.Scan(&count)
			return count == 0
		},
	}
}

type ValidationFailedError struct{}

func (e *ValidationFailedError) Error() string {
	return "The validation failed"
}

func ValidateStruct(c *fiber.Ctx, data any) error {
	return ValidateStructUpdates(c, nil, data)
}

func ValidateStructUpdates(c *fiber.Ctx, target any, data any) error {
	errors := map[string][]string{}

	dataValue := reflect.ValueOf(data).Elem()
	for i := 0; i < dataValue.NumField(); i++ {
		field := dataValue.Type().Field(i)
		valueField := dataValue.Field(i)
		var value *string
		if !valueField.IsZero() {
			if valueField.Kind() == reflect.Ptr {
				value = valueField.Interface().(*string)
			} else {
				value = valueField.Addr().Interface().(*string)
			}
		}

		formName := field.Tag.Get("form")
		tag := field.Tag.Get("validate")
		if tag != "" {
			for _, rule := range strings.Split(tag, "|") {
				parts := strings.Split(rule, ":")
				ruleName := parts[0]
				var ruleArgs []string = []string{}
				if len(parts) > 1 {
					ruleArgs = strings.Split(parts[1], ",")
				}

				if ruleName == "required" {
					if value == nil {
						errors[formName] = append(errors[formName], ruleName)
					}
				} else if value != nil {
					if ruleName == "nullable" {
						if *value == "" {
							return nil
						} else {
							continue
						}
					}
					if _, ok := Rules[ruleName]; ok {
						if !Rules[ruleName](ruleArgs, target, *value) {
							errors[formName] = append(errors[formName], rule)
						}
					} else {
						log.Fatalln("Validate: rule '" + ruleName + "' doesn't exists")
					}
				}
			}
		}
	}

	if len(errors) > 0 {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errors,
		})
		return &ValidationFailedError{}
	}
	return nil
}
