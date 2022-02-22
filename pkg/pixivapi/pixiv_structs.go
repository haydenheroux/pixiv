package pixivapi

type PixivResponseBody struct {
	Thumbnails PixivThumbnails `json:"body"`
}

type PixivThumbnails struct {
	Illusts PixivIllustrations `json:"thumbnails"`
}

type PixivIllustrations struct {
	IllustrationList []PixivIllustration `json:"illust"`
}

type PixivIllustration struct {
	Id                      string            `json:"id"`
	Title                   string            `json:"title"`
	IllustrationType        int               `json:"illustType"`
	XRestrict               int               `json:"xRestrict"`
	Restrict                int               `json:"restrict"`
	Sl                      int               `json:"sl"`
	URL                     string            `json:"url"`
	Tags                    []string          `json:"tags"`
	UserId                  string            `json:"userId"`
	UserName                string            `json:"userName"`
	Width                   int               `json:"width"`
	Height                  int               `json:"height"`
	PageCount               int               `json:"pageCount"`
	IsBookmarkable          bool              `json:"isBookmarkable"`
	BookmarkData            string            `json:"bookmarkData"`
	Alt                     string            `json:"alt"`
	TitleCaptionTranslation map[string]string `json:"titleCaptionTranslation"`
	CreateDate              string            `json:"createDate"`
	UpdateDate              string            `json:"updateDate"`
	IsUnlisted              bool              `json:"isUnlisted"`
	IsMasked                bool              `json:"isMasked"`
	URLs                    map[string]string `json:"urls"`
	ProfileImageURL         string            `json:"profileImageUrl"`
}
