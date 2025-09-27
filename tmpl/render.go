package tmpl

import (
	"bytes"
	"maps"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Render processes a text template with the given name, content, and data context,
func Render(name string, content string, data any) []byte {
	tmpl := template.New(name)
	tmpl.Funcs(template.FuncMap{
		"array":  func(v ...any) []any { return v },
		"append": func(slice []any, v ...any) []any { return append(slice, v...) },
		"capitalize": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"fromYaml": func(v string) (any, error) {
			var out any
			err := yaml.Unmarshal([]byte(v), &out)
			return out, err
		},
		"from": func(m any, k string) any {
			val := reflect.ValueOf(m)
			switch val.Kind() {
			case reflect.Map:
				keyVal := reflect.ValueOf(k)
				v := val.MapIndex(keyVal)
				if v.IsValid() {
					return v.Interface()
				}
				return nil
			case reflect.Struct:
				field := val.FieldByName(k)
				if field.IsValid() {
					return field.Interface()
				}
				return nil
			case reflect.Pointer:
				if val.IsNil() {
					return nil
				}
				elem := val.Elem()
				if elem.Kind() == reflect.Struct {
					field := elem.FieldByName(k)
					if field.IsValid() {
						return field.Interface()
					}
				}
				return nil
			default:
				return nil
			}
		},
		"upper": strings.ToUpper,
		"has": func(m any, k string) bool {
			val := reflect.ValueOf(m)
			switch val.Kind() {
			case reflect.Map:
				keyVal := reflect.ValueOf(k)
				v := val.MapIndex(keyVal)
				return v.IsValid()
			case reflect.Struct:
				field := val.FieldByName(k)
				return field.IsValid()
			default:
				return false
			}
		},
		"ifNil": func(def any, v any) any {
			if v == nil || (reflect.ValueOf(v).Kind() == reflect.Pointer && reflect.ValueOf(v).IsNil()) {
				return def
			}
			return v
		},
		"include": include(tmpl),
		"indent": func(spaces int, v string) string {
			pad := bytes.Repeat([]byte(" "), spaces)
			return string(pad) + string(bytes.ReplaceAll([]byte(v), []byte("\n"), append([]byte("\n"), pad...)))
		},
		"join": strings.Join,
		"map": func(v ...any) map[string]any {
			if len(v)%2 != 0 {
				panic("map function requires an even number of arguments")
			}
			m := make(map[string]any, len(v)/2)
			for i := 0; i < len(v); i += 2 {
				key, ok := v[i].(string)
				if !ok {
					panic("map keys must be strings")
				}
				m[key] = v[i+1]
			}
			return m
		},
		"merge": func(ms ...map[string]any) map[string]any {
			result := make(map[string]any)
			for _, m := range ms {
				maps.Copy(result, m)
			}
			return result
		},
		"pascal": func(s string) string {
			if s == "" {
				return s
			}
			parts := strings.FieldsFunc(s, func(r rune) bool {
				return r == '_' || r == '-' || r == ' '
			})
			for i, part := range parts {
				if part == "" {
					continue
				}
				parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			}
			return strings.Join(parts, "")
		},
		"lower": strings.ToLower,
		"set": func(m any, k string, v any) map[string]any {
			val := reflect.ValueOf(m)
			if val.Kind() != reflect.Map {
				panic("set function requires a map as the first argument")
			}
			keyVal := reflect.ValueOf(k)
			val.SetMapIndex(keyVal, reflect.ValueOf(v))
			result := make(map[string]any)
			for _, key := range val.MapKeys() {
				result[key.String()] = val.MapIndex(key).Interface()
			}
			return result
		},
		"strings": func(v ...string) []string { return v },
		"toYaml": func(v any) string {
			data, err := yaml.Marshal(v)
			if err != nil {
				panic(err)
			}
			return string(data)
		},
		"trim": strings.Trim,
	})
	tmpl = template.Must(tmpl.Parse(string(content)))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// RenderH is like Render but also takes helper templates to be included.
func RenderH(name string, content string, helpers string, data any) []byte {
	return Render(name, content, data)
}

func include(tmpl *template.Template) func(name string, data any) (string, error) {
	return func(name string, data any) (string, error) {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}
