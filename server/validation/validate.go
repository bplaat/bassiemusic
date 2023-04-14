package validation

import (
	"github.com/gofiber/fiber/v2"
)

// type ValidationRuleFunc func(args []string, value string) bool

func Validate(c *fiber.Ctx, data any) error {
	// TODO

	// rules := map[string]ValidationRuleFunc{
	// 	// String rules
	// 	"min": func(args []string, value string) bool {
	// 		min, err := strconv.ParseInt(args[0], 10, 32)
	// 		if err != nil {
	// 			log.Fatalln(err)
	// 		}
	// 		if len(*value) >= int(min) {
	// 			return true
	// 		}
	// 		return false
	// 	},
	// 	"max": func(args []string, value string) bool {
	// 		max, err := strconv.ParseInt(args[0], 10, 32)
	// 		if err != nil {
	// 			log.Fatalln(err)
	// 		}
	// 		if len(*value) <= int(max) {
	// 			return true
	// 		}
	// 		return false
	// 	},

	// 	// Type rules
	// 	// if name == "boolean" {
	// 	// 	if *value != "true" && *value != "false" {
	// 	// 		return false
	// 	// 	}
	// 	// }
	// 	// if name == "integer" {
	// 	// 	if _, err := strconv.ParseInt(*value, 10, 64); err != nil {
	// 	// 		return false
	// 	// 	}
	// 	// }
	// 	// if name == "uuid" {

	// 	// }
	// 	// if name == "email" {

	// 	// }
	// 	// if name == "enum" {

	// 	// }

	// 	// Database rules
	// 	// if name == "exists" {

	// 	// }
	// 	// if name == "unique" {

	// 	// }
	// }

	// errors := map[string][]string{}

	// value := reflect.ValueOf(data).Elem()
	// for i := 0; i < value.NumField(); i++ {
	// 	field := value.Type().Field(i)
	// 	valueField := value.Field(i)
	// 	var value *string
	// 	if !valueField.IsZero() {
	// 		if valueField.Kind() == reflect.Ptr {
	// 			value = valueField.Interface().(*string)
	// 		} else {
	// 			value = valueField.Addr().Interface().(*string)
	// 		}
	// 	}

	// 	formName := field.Tag.Get("form")
	// 	tag := field.Tag.Get("validate")
	// 	if tag != "" {
	// 		for _, rule := range strings.Split(tag, "|") {
	// 			parts := strings.Split(rule, ":")
	// 			ruleName := parts[0]
	// 			ruleArgs := strings.Split(parts[1], ",")
	// 			if value == nil &&  {
	// 				continue
	// 			}
	// 			if !rules[ruleName](ruleArgs, value) {
	// 				errors[formName] = append(errors[formName], ruleName)
	// 			}
	// 		}
	// 	}
	// }

	// if len(errors) > 0 {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"success": false,
	// 		"errors":  errors,
	// 	})
	// }
	return nil
}
