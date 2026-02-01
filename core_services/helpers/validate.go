package helpers

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func MakeValidate(fiedValidate map[string]interface{}) (map[string]string, error) {
	mergedErrors := make(map[string]interface{})
	for keyField, fiedValidated := range fiedValidate {
		if err, ok := fiedValidated.(error); ok {
			fieldError := fmt.Sprintf("%s: %s", keyField, err.Error())
			parsedErrors := ConverstError(errors.New(fieldError))
			mergedErrors[keyField] = parsedErrors[keyField]
			continue
		}

		if IsZero(reflect.ValueOf(fiedValidated)) {
			fieldError := fmt.Sprintf("%s: zero value", keyField)
			parsedErrors := ConverstError(errors.New(fieldError))
			if existing, exists := mergedErrors[keyField]; exists {
				if slice, ok := existing.([]string); ok {
					mergedErrors[keyField] = append(slice, parsedErrors[keyField])
				} else {
					mergedErrors[keyField] = []string{existing.(string), parsedErrors[keyField]}
				}
				continue
			} else {
				mergedErrors[keyField] = parsedErrors[keyField]
			}
		}

		val := reflect.ValueOf(fiedValidated)
		fmt.Println("Kind:", val.Kind()) // พิมพ์ชนิดข้อมูลของ fiedValidated
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			if val.IsNil() || val.Len() == 0 {
				fieldError := fmt.Sprintf("%s: zero value", keyField)
				parsedErrors := ConverstError(errors.New(fieldError))
				if existing, exists := mergedErrors[keyField]; exists {
					if slice, ok := existing.([]string); ok {
						mergedErrors[keyField] = append(slice, parsedErrors[keyField])
					} else {
						mergedErrors[keyField] = []string{existing.(string), parsedErrors[keyField]}
					}
				} else {
					mergedErrors[keyField] = parsedErrors[keyField]
				}
				continue
			}
		}

		if data, ok := fiedValidated.([]map[string]interface{}); ok {
			for index, fieldValues := range data {
				for keyFieldName, value := range fieldValues {
					if err, ok := value.(error); ok {
						fieldKey := fmt.Sprintf("[%d].%s", index, keyFieldName)
						fieldError := fmt.Sprintf("[%d].%s: %s", index, keyFieldName, err.Error())
						parsedErrors := ConverstError(errors.New(fieldError))
						if existing, exists := mergedErrors[keyField]; exists {
							switch existing := existing.(type) {
							case string:
								mergedErrors[keyField] = []string{existing, parsedErrors[fieldKey]}
							case []string:
								mergedErrors[keyField] = append(existing, parsedErrors[fieldKey])
							case map[string]string:
								existing[fieldKey] = parsedErrors[fieldKey]
								mergedErrors[keyField] = existing
							default:
								mergedErrors[keyField] = []string{parsedErrors[fieldKey]}
							}
							continue
						} else {
							mergedErrors[keyField] = map[string]string{
								fieldKey: parsedErrors[fieldKey],
							}
						}
					}

					if IsZero(reflect.ValueOf(value)) {
						fieldKey := fmt.Sprintf("[%d].%s", index, keyFieldName)
						fieldError := fmt.Sprintf("[%d].%s: zero value", index, keyFieldName)
						parsedErrors := ConverstError(errors.New(fieldError))
						if existing, exists := mergedErrors[keyField]; exists {
							switch existing := existing.(type) {
							case string:
								mergedErrors[keyField] = []string{existing, parsedErrors[fieldKey]}
							case []string:
								mergedErrors[keyField] = append(existing, parsedErrors[fieldKey])
							case map[string]string:
								existing[fieldKey] = parsedErrors[fieldKey]
								mergedErrors[keyField] = existing
							default:
								mergedErrors[keyField] = []string{parsedErrors[fieldKey]}
							}
						} else {
							mergedErrors[keyField] = map[string]string{
								fieldKey: parsedErrors[fieldKey],
							}
						}
						continue
					}

					val := reflect.ValueOf(value)
					if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
						if val.IsNil() || val.Len() == 0 {
							fieldKey := fmt.Sprintf("[%d].%s", index, keyFieldName)
							fieldError := fmt.Sprintf("[%d].%s: zero value", index, keyFieldName)
							parsedErrors := ConverstError(errors.New(fieldError))
							if existing, exists := mergedErrors[keyField]; exists {
								switch existing := existing.(type) {
								case string:
									mergedErrors[keyField] = []string{existing, parsedErrors[fieldKey]}
								case []string:
									mergedErrors[keyField] = append(existing, parsedErrors[fieldKey]) // ใช้ append แทนที่จะเขียนทับเป็น map
								case map[string]string:
									existing[fieldKey] = parsedErrors[fieldKey]
									mergedErrors[keyField] = existing
								default:
									mergedErrors[keyField] = []string{parsedErrors[fieldKey]}
								}
							} else {
								mergedErrors[keyField] = map[string]string{
									fieldKey: parsedErrors[fieldKey],
								}
							}
							continue
						}
					}
				}
			}
		}

		if data, ok := fiedValidated.(map[string]interface{}); ok {
			for key, value := range data {
				if err, ok := value.(error); ok {
					fieldError := fmt.Sprintf("%s: %s", key, err.Error())
					parsedErrors := ConverstError(errors.New(fieldError))
					if existing, exists := mergedErrors[keyField]; exists {
						switch existing := existing.(type) {
						case string:
							mergedErrors[keyField] = []string{existing, parsedErrors[key]}
						case []string:
							mergedErrors[keyField] = append(existing, parsedErrors[key]) // ใช้ append แทนที่จะเขียนทับเป็น map
						case map[string]string:
							existing[key] = parsedErrors[key]
							mergedErrors[keyField] = existing
						default:
							mergedErrors[keyField] = []string{parsedErrors[key]}
						}
					} else {
						mergedErrors[keyField] = map[string]string{
							key: parsedErrors[key],
						}
					}
					continue
				}

				if IsZero(reflect.ValueOf(value)) {
					fieldError := fmt.Sprintf("%s: zero value", key)
					parsedErrors := ConverstError(errors.New(fieldError))
					if existing, exists := mergedErrors[keyField]; exists {
						switch existing := existing.(type) {
						case string:
							mergedErrors[keyField] = []string{existing, parsedErrors[key]}
						case []string:
							mergedErrors[keyField] = append(existing, parsedErrors[key]) // ใช้ append แทนที่จะเขียนทับเป็น map
						case map[string]string:
							existing[key] = parsedErrors[key]
							mergedErrors[keyField] = existing
						default:
							mergedErrors[keyField] = []string{parsedErrors[key]}
						}
					} else {
						mergedErrors[keyField] = map[string]string{
							key: parsedErrors[key],
						}
					}
				}

				val := reflect.ValueOf(value)
				if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
					if val.IsNil() || val.Len() == 0 {
						fieldError := fmt.Sprintf("%s: zero value", key)
						parsedErrors := ConverstError(errors.New(fieldError))
						if existing, exists := mergedErrors[keyField]; exists {
							switch existing := existing.(type) {
							case string:
								mergedErrors[keyField] = []string{existing, parsedErrors[key]}
							case []string:
								mergedErrors[keyField] = append(existing, parsedErrors[key]) // ใช้ append แทนที่จะเขียนทับเป็น map
							case map[string]string:
								existing[key] = parsedErrors[key]
								mergedErrors[keyField] = existing
							default:
								mergedErrors[keyField] = []string{parsedErrors[key]}
							}
						} else {
							mergedErrors[keyField] = map[string]string{
								key: parsedErrors[key],
							}
						}
						continue
					}
				}
			}
		}
	}

	if len(mergedErrors) > 0 {
		allErrors := make(map[string]string)
		FlattenErrors("", mergedErrors, allErrors)
		return allErrors, nil
	}
	return nil, nil
}

func ConverstError(err error) map[string]string {
	stringSlice := strings.Split(err.Error(), ",")
	strErrors := make(map[string]string)
	for _, err := range stringSlice {
		if !strings.Contains(err, ":") {
			continue
		}
		stringSlice2 := strings.Split(err, ":")
		errField := ToSnakeCase(strings.ReplaceAll(stringSlice2[0], " ", ""))
		errMessage := strings.TrimSpace(stringSlice2[1])
		strErrors[errField] = HandleErrMesssage(errField, errMessage)
	}
	return strErrors
}

func HandleErrMesssage(errField, err string) string {
	switch err {
	case "zero value":
		switch errField {
		default:
			return "โปรดระบุ"
		}
	case "nonzero":
		return "โปรดระบุ"
	default:
		return "เกิดข้อผิดพลาด โปรดลองอีกครั้ง"
	}
}

func IsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	default:
		return false
	}
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func FlattenErrors(prefix string, input interface{}, output map[string]string) {
	// ตรวจสอบชนิดของข้อมูลที่เป็น map
	switch v := input.(type) {
	case map[string]interface{}:
		// ถ้าเป็น map ให้วนลูปเพื่อเรียกฟังก์ชันซ้ำ
		for key, value := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}
			FlattenErrors(newPrefix, value, output)
		}
	case map[string]string:
		// ถ้าเป็น map[string]string ให้เก็บข้อมูล key-value ลงใน output
		for key, value := range v {
			newPrefix := prefix + "." + key
			output[newPrefix] = value
		}
	case []interface{}:
		// ถ้าเป็น slice ให้ใช้ index แสดงค่าของมัน
		for index, item := range v {
			newPrefix := fmt.Sprintf("%s[%d]", prefix, index)
			FlattenErrors(newPrefix, item, output)
		}
	case string:
		// ถ้าค่าคือ string เก็บข้อมูลลงใน output
		output[prefix] = v
	default:
		// สำหรับค่าอื่นๆ (ที่ไม่ใช่ map หรือ slice)
		output[prefix] = fmt.Sprintf("%v", v)
	}
}
