package normalizetext

import "strings"

func normalizeText(textBytes []byte) string {
	var out []string
	var currentLine []rune
	text := string(textBytes)
	for _, b := range text {
		switch {
		case isRunePunctuation(b):
			if len(currentLine) > 0 {
				out = append(out, string(currentLine))
			}
			currentLine = []rune{}
		case b == 10 || b == 13 || b == 32:
			if len(currentLine) > 0 && currentLine[len(currentLine)-1] != 32 {
				currentLine = append(currentLine, 32)
			}
		case isRuneLowercaseCharacter(b):
			currentLine = append(currentLine, b)
		case isRuneUppercaseCharacter(b):
			currentLine = append(currentLine, b+32)
		default:
			// no-op
		}
	}
	if len(currentLine) > 0 {
		out = append(out, string(currentLine))
	}
	return strings.Join(out, "\n")
}

func isRunePunctuation(b rune) bool {
	return b == 46 || b == 33 || b == 63
}

func isRuneLowercaseCharacter(b rune) bool {
	isEnglishLowercase := b >= 97 && b <= 122
	isExtendedLowercase := b >= 224 && b <= 246 || b >= 248 && b <= 253
	return isEnglishLowercase || isExtendedLowercase
}

func isRuneUppercaseCharacter(b rune) bool {
	isEnglishUppercase := b >= 65 && b <= 90
	isExtendedUppercase := b >= 192 && b <= 214 || b >= 216 && b <= 221
	return isEnglishUppercase || isExtendedUppercase
}
