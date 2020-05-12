package vk

import (
	"bytes"
	"encoding/json"
)

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

type videoImage struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

type VideoImage struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

func (v *VideoImage) UnmarshalJSON(b []byte) error {
	if bytes.HasPrefix(b, []byte(`[`)) {
		// Blank image.
		return nil
	}
	var im videoImage
	if err := json.Unmarshal(b, &im); err != nil {
		return err
	}
	*v = VideoImage(im)
	return nil
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
