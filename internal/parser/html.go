package parser

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
)

const(
	REGEX_BASE_URL = `<base.*?href="(.*?)"`
	REGEX_ANCHOR_URL = `<a.*?href="(.*?)"`
)

func GetLinks(page http.Response) []*url.URL {
	body := page.Body

	defer body.Close()

	s, err := io.ReadAll(body)

	if err != nil {
		return nil
	}

	anchorUrls := extractRegexp(string(s), REGEX_ANCHOR_URL)
	baseUrl := extractRegexp(string(s),REGEX_BASE_URL)

	var valid []*url.URL

	for _, a := range anchorUrls {
		var u = &url.URL{}
		var err error

		if isFullUrl(a) {
			u, err = url.Parse(a)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			valid = append(valid, u)
			continue
		}

		if len(baseUrl) > 0 {
			b, err := url.Parse(baseUrl[0])

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			u.Scheme = b.Scheme
			u.Host = b.Host
			u.Path = path.Join(b.Path, a)

			valid = append(valid, u)
			continue
		}

		u.Scheme = page.Request.URL.Scheme
		u.Host = page.Request.URL.Host
		u.Path = path.Join(page.Request.URL.Path, a)

		valid = append(valid, u)
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