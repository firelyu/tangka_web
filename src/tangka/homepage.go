package tangka

type Homepage struct {
	WebName string
}

func NewHomepage(webname string) *Homepage {
	return &Homepage{WebName:webname}
}
