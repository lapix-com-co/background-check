package reply

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	capitalsPattern = regexp.MustCompile(`c.?al es la capital del? ([\w\sáéíóú]+) ?(\(sin tilde\))?\??`)
	capitals        = map[string]string{
		"antioquia":        "Medellín",
		"colombia":         "Bogotá",
		"boyacá":           "Tunja",
		"valle del cauca":  "Cali",
		"vallle del cauca": "Cali",
		"atlantico":        "barranquilla",
		"amazonas":         "leticia",
	}
)

var _ Strategy = &Capitals{}

// NoAccentsCapitals returns the capital and removes the accent.
type NoAccentsCapitals struct {
	capitals *Capitals
}

// NewNoAccentsCapitals constructor, it add the default capitals.
func NewNoAccentsCapitals() *NoAccentsCapitals {
	return &NoAccentsCapitals{NewCapitals()}
}

// Is will match is the question match the capitalsPattern c.?al es la capital de ([a-zA-Z]+).
func (c NoAccentsCapitals) Is(u User, i string) bool {
	return c.capitals.Is(u, i)
}

// Answer will find the capital for the given location and will removes the accents.
func (c NoAccentsCapitals) Answer(u User, i string) (string, bool, error) {
	answer, ok, err := c.capitals.Answer(u, i)
	if err != nil || !ok {
		return answer, ok, err
	}

	return removeAccents(answer), true, nil
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func removeAccents(i string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, i)
	return result
}

// Capitals will find any capital for the given country or city.
type Capitals struct {
	capitals map[string]string
}

// NewCapitals add the default capitals to the new struct
func NewCapitals() *Capitals {
	return &Capitals{capitals}
}

// Is will match is the question match the capitalsPattern c.?al es la capital de ([a-zA-Z]+)
func (c Capitals) Is(u User, i string) bool {
	// don't know how to make the capitalsPattern case-insensitive
	return capitalsPattern.MatchString(strings.ToLower(i))
}

// Answer will find the capital for the given location.
func (c Capitals) Answer(u User, i string) (string, bool, error) {
	matches := capitalsPattern.FindAllStringSubmatch(strings.ToLower(i), 1)
	if len(matches) == 0 {
		return "", false, ErrInvalidQuestion
	}

	place := matches[0][1]
	if v, ok := c.capitals[sanitized(place)]; ok {
		return v, true, nil
	}

	return "", false, nil
}

func sanitized(i string) string {
	return strings.TrimSpace(strings.ToLower(i))
}
