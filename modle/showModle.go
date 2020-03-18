package modle

type Video struct {
	Id        int    `gorm:"AUTO_INCREMENT"`
	Vid       string `gorm:"size:255;index"`
	VideoName string `gorm:"size:255"`
	Img       string `gorm:"size255"`
	Language  string `gorm:"size:255"`

	Describe string `gorm:"size:255"`
	Tag      string `gorm:"size:255"`
	Mark     float64
	Episode  []Episode `gorm:"ForeignKey:VideoId;AssociationForeignKey:VId"`
}

type Episode struct {
	Id            int    `gorm:"AUTO_INCREMENT"`
	Vid           string `gorm:"size:255"`
	EpisodeName   string `gorm:"size:255"`
	EpisodeNumber int
	PlayImg       string `gorm:"size:255"`
	Img           string `gorm:"size:255"`
	Content       string
	Describe      string
	Read          int
	UrlPath       string `gorm:"size:255"`
	Path          string `gorm:"size:255"`
	Tag           string `gorm:"size:255"`
	Language      string `gorm:"size:255"`
	Season        string `gorm:"size:255"`
}
type Anime struct {
	Id            int `json:"id" orm:"pk"`
	MovieId       string
	EpisodeNumber string
	EpisodeName   string
	Download      string
	EpisodeUrl    string
	EpisodeImg    string
	EpisodeInfo   string
	DeEpisodeInfo string
}
