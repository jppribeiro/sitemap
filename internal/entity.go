package internal

type Url struct {
	Loc string `xml:"loc"`
}

type Urlset struct {
	Url []Url `xml:"url"`
}

