package actpub

import "time"

type JSONTime time.Time

const TimeFmtStr = "2006-01-02T15:04:05Z"

func (t JSONTime) MarshalText() ([]byte, error) {
	tim := time.Time(t)
	return []byte(tim.Format(TimeFmtStr)), nil
}

func (t *JSONTime) UnmarshalText(text []byte) error {
	val, err := time.Parse(TimeFmtStr, string(text))
	if err != nil {
		return err
	} else {
		*t = JSONTime(val)
		return nil
	}
}

// AsObject represents an Activity Streams Object
type AsObject struct {
	AtContext []any `json:"@context,omitempty"`

	Type         string    `json:"type,omitempty"`
	Id           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Summary      string    `json:"summary,omitempty"`
	Content      string    `json:"content,omitempty"`
	InReplyTo    string    `json:"inReplyTo,omitempty"`
	AttributedTo string    `json:"attributedTo,omitempty"`
	Published    *JSONTime `json:"published,omitempty"`
	To           string    `json:"to,omitempty"`
	Bto          string    `json:"bto,omitempty"`
	Cc           string    `json:"cc,omitempty"`
	Bcc          string    `json:"bcc,omitempty"`
	URL          string    `json:"url,omitempty"`

	Tag        []any `json:"tag,omitempty"`
	Icon       any   `json:"icon,omitempty"`
	Image      any   `json:"image,omitempty"`
	Audience   any   `json:"audience,omitempty"`
	Attachment []any `json:"attachment,omitempty"`
	Replies    any   `json:"replies,omitempty"`
}
