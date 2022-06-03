package parser

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const(
	REGEX_BASE_URL = `<base.*?href="(.*?)"`
	REGEX_ANCHOR_URL = `<a.*?href="(.*?)"`
)

func GetLinks(page http.Response) []string {
	body := page.Body

	defer body.Close()

	s, err := io.ReadAll(body)

	if err != nil {
		return nil
	}

	baseUrl := extractRegexp(string(s),REGEX_BASE_URL)

	anchorUrls := extractRegexp(string(s), REGEX_ANCHOR_URL)

	var valid []string

	for _, a := range anchorUrls {
		if !isFullUrl(a) {
			if len(baseUrl) > 0 {
				valid = append(valid, fmt.Sprintf("%s%s", baseUrl[0], a))
			}

			continue
		}

		valid = append(valid, a)
	}

	return valid
}

func extractRegexp(s string, r string) []string {
	ex := regexp.MustCompile(r)

	m := ex.FindAllStringSubmatch(s, -1)

	res := make([]string, len(m))

	for i, match := range m {
		res[i] = match[1]
	}

	return res
}

func isFullUrl(u string) bool {
	r := regexp.MustCompile(`^http`)

	return r.MatchString(u)
}