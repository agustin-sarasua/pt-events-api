package main

type Event struct {
	Name                  string    `json:"Name,omitempty"`
	ID                    string    `json:"Id"`
	EventType             string    `json:"EventType"`
	StartTime             string    `json:"StartTime,omitempty"`
	EndTime               string    `json:"EndTime,omitempty"`
	Sub                   string    `json:"Sub,omitempty"`
	PlaceID               string    `json:"PlaceId,omitempty"`
	AllowedPersons        []*Person `json:"AllowedPersons,omitempty"`
	TrackUnknownPersons   bool      `json:"TrackUnknownPersons,omitempty"`
	NotifyWhenUnknown     bool      `json:"NotifyWhenUnknown,omitempty"`
	NotifyWhenAllowed     bool      `json:"NotifyWhenAllowed,omitempty"`
	UploadPicturesUnknown bool      `json:"UploadPicturesUnknown,omitempty"`
	UploadPicturesAllowed bool      `json:"UploadPicturesAllowed,omitempty"`
}

type Person struct {
	ID       string   `json:"Id"`
	Name     string   `json:"Name,omitempty"`
	Pictures []string `json:"Pictures,omitempty"`
	Sub      string   `json:"Sub"`
}
