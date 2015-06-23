package newsGetter

type jsonNewsBody struct {
	By string
	Id int
	//Kids 			[]int
	Score          int
	Text           string
	Time           int
	Title          string
	Type           string
	Url            string
	ProviderName   string
	ProviderUrl    string
	CreatedAt      string
	RelatedStories []RelatedStories
	Category       TopicIdentity
	Image          interface{}
	NewsPageView   int `bson:"news_page_view"`
}

// Topics topics list map holder
type Topics map[string]TopicIdentity

//TopicIdentity topic identifier
type TopicIdentity struct {
	Initial string `json:"initial"`
	Name    string `json:"name"`
}

// type GoogleNewsResults struct {
// 	GsearchResultClass string
// 	ClusterUrl         string
// 	Content            string
// 	UnescapedUrl       string
// 	Url                string
// 	Title              string
// 	TitleNoFormatting  string
// 	Publish            string
// 	PublishedDate      string
// 	Language           string
// 	RelatedStories     []RelatedStories
// 	Image              Image
// 	Category           TopicIdentity
// }

// RelatedStories google related stories
type RelatedStories struct {
	Language          string `json:"language"`
	Location          string `json:"location"`
	PublishedDate     string `json:"publishedDate"`
	Publisher         string `json:"publisher"`
	SignedRedirectURL string `json:"signedRedirectUrl"`
	Title             string `json:"title"`
	TitleNoFormatting string `json:"titleNoFormatting"`
	UnescapedURL      string `json:"unescapedUrl"`
	URL               string `json:"url"`
}

// Image google news item top image
type Image struct {
	Publisher string `json:"publisher"`
	URL       string `json:"url"`
}

// GoogleNewsResults google news result struct
type GoogleNewsResults struct {
	GsearchResultClass string `json:"GsearchResultClass"`
	ClusterURL         string `json:"clusterUrl"`
	Content            string `json:"content"`
	Image              struct {
		OriginalContextURL string `json:"originalContextUrl"`
		Publisher          string `json:"publisher"`
		TbHeight           int    `json:"tbHeight"`
		TbURL              string `json:"tbUrl"`
		TbWidth            int    `json:"tbWidth"`
		URL                string `json:"url"`
	} `json:"image"`
	Language          string           `json:"language"`
	Location          string           `json:"location"`
	PublishedDate     string           `json:"publishedDate"`
	Publisher         string           `json:"publisher"`
	RelatedStories    []RelatedStories `json:"relatedStories"`
	SignedRedirectURL string           `json:"signedRedirectUrl"`
	Title             string           `json:"title"`
	TitleNoFormatting string           `json:"titleNoFormatting"`
	UnescapedURL      string           `json:"unescapedUrl"`
	URL               string           `json:"url"`
	Category          TopicIdentity
}
