package modle

type Video struct {
	Id        int `gorm:"AUTO_INCREMENT"`
	Vid       string
	VideoName string `gorm:"size:255"`
	Img       string
	Language  string
	Describe  string
	Tag       string
	Mark      float64
}

type Episode struct {
	Id            int `gorm:"AUTO_INCREMENT"`
	Vid           string
	EpisodeName   string `gorm:"size:255"`
	EpisodeNumber int
	PlayImg       string
	Img           string
	Content       string
	Describe      string
	Read          int
	UrlPath       string
	Path          string
	Tag           string
	Language      string
	Season        string
}
