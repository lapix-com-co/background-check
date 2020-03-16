package reply

import (
	"regexp"
	"strings"
)

type DocumentCharts struct{}

var (
	documentPattern = regexp.MustCompile(`escriba (los|el) ?(.*) (ultimos?|primero?s?) digitos? del documento a consultar`)
	namePattern     = regexp.MustCompile(`escriba l([ao])s? ([\w]+)? ?(primer[as]{0,2}|ultimas?) letras? del ([\w\s]+) de la persona`)
)

func (d DocumentCharts) Is(u User, i string) bool {
	li := strings.ToLower(i)
	return documentPattern.MatchString(li) || namePattern.MatchString(li)
}

func (d DocumentCharts) Answer(u User, i string) (string, bool, error) {
	matched := documentPattern.FindAllStringSubmatch(strings.ToLower(i), 1)

	if len(matched) == 0 {
		matched = namePattern.FindAllStringSubmatch(strings.ToLower(i), 4)

		if len(matched) == 0 {
			return "", false, ErrInvalidQuestion
		}

		switch matched[0][4] {
		case "primer nombre":
			if matched[0][3] == "primeras" || matched[0][3] == "primer" {
				switch matched[0][2] {
				case "":
					return u.Name.FirstName[:1], true, nil
				case "dos":
					return u.Name.FirstName[:2], true, nil
				case "tres":
					return u.Name.FirstName[:3], true, nil
				}

			} else {
				switch matched[0][2] {
				case "":
					return u.Name.FirstName[len(u.Name.FirstName)-1:], true, nil
				case "dos":
					return u.Name.FirstName[len(u.Name.FirstName)-2:], true, nil
				case "tres":
					return u.Name.FirstName[len(u.Name.FirstName)-3:], true, nil
				}
			}
		}
		return "", false, nil
	}

	if matched[0][3] == "primeras" || matched[0][3] == "primeros" {
		switch matched[0][2] {
		case "dos":
			return u.Document.Number[:2], true, nil
		case "tres":
			return u.Document.Number[:3], true, nil
		}
	} else {
		switch matched[0][2] {
		case "dos":
			return u.Document.Number[len(u.Document.Number)-2:], true, nil
		case "tres":
			return u.Document.Number[len(u.Document.Number)-3:], true, nil
		}
	}

	return "", false, nil
}
