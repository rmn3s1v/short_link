package model

import "time"

type URL struct{
	ID int64
	OriginURL string
	ShortURL string
	DateCreate time.Time
}
