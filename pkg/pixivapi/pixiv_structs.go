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
	// Id: an identification string and/or number Pixiv assigns to an illustration
	Id string `json:"id"`
	// Title: the title a Pixiv user provides to an illustration on upload
	Title string `json:"title"`
	// IllustrationType: unknown
	IllustrationType int `json:"illustType"`
	// XRestrict: unknown, probably indicates NSFW/R18 content
	XRestrict int `json:"xRestrict"`
	// Restrict: unknown
	Restrict int `json:"restrict"`
	// Sl: unknown
	Sl int `json:"sl"`
	// URL: a URL for a smaller version of the image, for use in rendering the webpage
	URL string `json:"url"`
	// Tags: the tags a Pixiv user provides to an illustration on upload
	Tags []string `json:"tags"`
	// UserId: an identification string and/or number Pixiv assigns to the user who uploaded
	// the illustration
	UserId string `json:"userId"`
	// UserName: the name a Pixiv user used when they uploaded the illustration
	UserName string `json:"userName"`
	// Width: the width of the original image, presumably in pixels
	Width int `json:"width"`
	// Height: the height of the original image, presumably in pixels
	Height int `json:"height"`
	// PageCount: unknown
	PageCount int `json:"pageCount"`
	// IsBookmarkable: unknown, likely relates to the bookmarking feature
	IsBookmarkable bool `json:"isBookmarkable"`
	// BookmarkData: unknown, likely relates to the bookmarking feature
	BookmarkData string `json:"bookmarkData"`
	// Alt: the alternate name / subtitle that a user enters when uploading an illustration
	Alt string `json:"alt"`
	// TitleCaptionTranslation: unknown, presumably provides an alternate title for the illustration
	// if the user's browser's language differs from the langauge the author uploaded in
	TitleCaptionTranslation map[string]string `json:"titleCaptionTranslation"`
	// CreateDate: a string representing the original creation date of the illustration
	CreateDate string `json:"createDate"`
	// UpdateDate: a string representing the date the illustration was last updated
	UpdateDate string `json:"updateDate"`
	// IsUnlisted: unknown
	IsUnlisted bool `json:"isUnlisted"`
	// IsMasked: unknown
	IsMasked bool `json:"isMasked"`
	// URLs: a map of URLs which provide the image in various sizes, similar to the URL field
	URLs map[string]string `json:"urls"`
	// ProfileImageURL: a URL for a smaller version of the uploader's profile image, for rendering on the webpage
	ProfileImageURL string `json:"profileImageUrl"`
}
