package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var aiNumberPattern = regexp.MustCompile(`-?\d+(\.\d+)?`)

func decodeAIJSONContent(content string) (json.RawMessage, error) {
	text := strings.TrimSpace(strings.TrimPrefix(content, "\ufeff"))
	if text == "" {
		return nil, ErrAIInvalidResponse
	}

	candidates := []string{text}
	if stripped := stripMarkdownJSONFence(text); stripped != text {
		candidates = append(candidates, stripped)
	}

	for _, candidate := range candidates {
		if raw, err := decodeFirstJSONRaw(candidate); err == nil && isJSONObject(raw) {
			return raw, nil
		}
		for i := 0; i < len(candidate); i++ {
			if candidate[i] != '{' {
				continue
			}
			raw, err := decodeFirstJSONRaw(candidate[i:])
			if err == nil && isJSONObject(raw) {
				return raw, nil
			}
		}
	}

	return nil, ErrAIInvalidResponse
}

func decodeFirstJSONRaw(text string) (json.RawMessage, error) {
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	var raw json.RawMessage
	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func stripMarkdownJSONFence(text string) string {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "```") {
		return text
	}
	firstLineEnd := strings.IndexByte(text, '\n')
	if firstLineEnd < 0 {
		return text
	}
	text = text[firstLineEnd+1:]
	if end := strings.LastIndex(text, "```"); end >= 0 {
		text = text[:end]
	}
	return strings.TrimSpace(text)
}

func isJSONObject(raw json.RawMessage) bool {
	raw = bytes.TrimSpace(raw)
	return len(raw) > 0 && raw[0] == '{'
}

func isJSONArray(raw json.RawMessage) bool {
	raw = bytes.TrimSpace(raw)
	return len(raw) > 0 && raw[0] == '['
}

func unwrapAIObject(raw json.RawMessage, arrayKey string, coreKeys ...string) json.RawMessage {
	for depth := 0; depth < 4; depth++ {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil {
			return raw
		}
		if hasAnyRawKey(obj, coreKeys...) {
			return raw
		}

		for _, key := range []string{"data", "result", "output", "response", "content", "message", "menu", "recommendation"} {
			value := rawValue(obj, key)
			if len(value) == 0 {
				continue
			}
			value = bytes.TrimSpace(value)
			switch {
			case isJSONObject(value):
				raw = value
				goto nextDepth
			case isJSONArray(value):
				return wrapJSONArray(arrayKey, value)
			case len(value) > 0 && value[0] == '"':
				var nested string
				if err := json.Unmarshal(value, &nested); err == nil {
					if nestedRaw, err := decodeAIJSONContent(nested); err == nil {
						raw = nestedRaw
						goto nextDepth
					}
				}
			}
		}
		return raw
	nextDepth:
	}
	return raw
}

func wrapJSONArray(key string, value json.RawMessage) json.RawMessage {
	encodedKey, _ := json.Marshal(key)
	return json.RawMessage(fmt.Sprintf("{%s:%s}", encodedKey, bytes.TrimSpace(value)))
}

func hasAnyRawKey(obj map[string]json.RawMessage, keys ...string) bool {
	for _, key := range keys {
		if len(rawValue(obj, key)) > 0 {
			return true
		}
	}
	return false
}

func rawValue(obj map[string]json.RawMessage, keys ...string) json.RawMessage {
	for _, key := range keys {
		if value, ok := obj[key]; ok {
			return value
		}
	}
	for objKey, value := range obj {
		for _, key := range keys {
			if strings.EqualFold(objKey, key) {
				return value
			}
		}
	}
	return nil
}

func rawString(obj map[string]json.RawMessage, keys ...string) string {
	return rawToString(rawValue(obj, keys...))
}

func rawToString(raw json.RawMessage) string {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		return number.String()
	}
	var boolean bool
	if err := json.Unmarshal(raw, &boolean); err == nil {
		if boolean {
			return "true"
		}
		return "false"
	}
	if isJSONObject(raw) {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err == nil {
			return firstNonEmpty(rawString(obj, "name", "title", "text", "content", "description", "名称", "标题"))
		}
	}
	return strings.TrimSpace(string(raw))
}

func rawInt(obj map[string]json.RawMessage, keys ...string) int {
	return rawToInt(rawValue(obj, keys...))
}

func rawToInt(raw json.RawMessage) int {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return 0
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		if value, err := strconv.Atoi(number.String()); err == nil {
			return value
		}
		if value, err := strconv.ParseFloat(number.String(), 64); err == nil {
			return int(value)
		}
	}
	text := rawToString(raw)
	match := aiNumberPattern.FindString(text)
	if match == "" {
		return 0
	}
	value, _ := strconv.ParseFloat(match, 64)
	return int(value)
}

func rawFloat(obj map[string]json.RawMessage, keys ...string) float64 {
	raw := bytes.TrimSpace(rawValue(obj, keys...))
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return 0
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		value, _ := strconv.ParseFloat(number.String(), 64)
		return value
	}
	text := rawToString(raw)
	match := aiNumberPattern.FindString(text)
	if match == "" {
		return 0
	}
	value, _ := strconv.ParseFloat(match, 64)
	return value
}

func rawStringList(obj map[string]json.RawMessage, limit int, keys ...string) []string {
	return rawToStringList(rawValue(obj, keys...), limit)
}

func rawToStringList(raw json.RawMessage, limit int) []string {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil
	}
	var values []json.RawMessage
	if err := json.Unmarshal(raw, &values); err == nil {
		result := make([]string, 0, len(values))
		for _, value := range values {
			text := rawToString(value)
			if text != "" {
				result = append(result, text)
			}
		}
		return normalizeStringList(result, limit)
	}
	text := rawToString(raw)
	if text == "" {
		return nil
	}
	return normalizeStringList(splitAITextList(text), limit)
}

func splitAITextList(text string) []string {
	text = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(text)
	parts := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n' || r == ',' || r == '，' || r == '、' || r == ';' || r == '；' || r == '/'
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = cleanAIListItem(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 && strings.TrimSpace(text) != "" {
		result = append(result, strings.TrimSpace(text))
	}
	return result
}

func cleanAIListItem(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimLeft(text, "-*•· ")
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`^\d+[\.\)、\)]\s*`).ReplaceAllString(text, "")
	return strings.TrimSpace(text)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
