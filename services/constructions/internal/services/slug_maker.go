package services

import "github.com/gosimple/slug"

func MakeSlug(s string) string {
	slug.MaxLength = 64
	return slug.MakeLang(s, "ru")
}
