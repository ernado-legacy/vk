package vk

const (
	methodVideoGet = "video.get"
)

type Video struct {
	Resource
}

type VideoGetFields struct {
	Offset   int    `url:"offset,omitempty"`
	Count    int    `url:"count,omitempty"`
	Extended Bool   `url:"extended,omitempty"`
	Videos   string `url:"videos,omitempty"`
}

type VideoImage struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

type VideoItem struct {
	Duration int          `json:"duration"`
	Player   string       `json:"player"`
	Files    VideoFiles   `json:"files"`
	Images   []VideoImage `json:"image"`
}

type VideoFiles struct {
	External string `json:"external"`
}

type VideoGetResult struct {
	Count int         `json:"count"`
	Items []VideoItem `json:"items"`
}

func (v Video) Get(fields VideoGetFields) (result VideoGetResult, err error) {
	return result, v.Decode(v.Request(methodVideoGet, fields), &result)
}
