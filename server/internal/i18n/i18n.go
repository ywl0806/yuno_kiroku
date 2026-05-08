package i18n

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"sync"

	_i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const DefaultLocale = "ko"

var (
	mu     sync.RWMutex
	bundle *_i18n.Bundle
)

//go:embed locales
var locales embed.FS

// Init 메시지 맵 초기화 (앱 시작 시 호출)
func Init() {
	mu.Lock()
	defer mu.Unlock()
	bundle = _i18n.NewBundle(language.Korean)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	langs, err := fs.Glob(locales, "locales/*.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, lang := range langs {
		bundle.LoadMessageFileFS(locales, lang)
	}

}

// SupportedLocales 지원 로케일 목록
func SupportedLocales() []language.Tag {
	return bundle.LanguageTags()
}

// T 번역 함수
func T(locale []string, key string, templateData map[string]string) string {
	localizer := _i18n.NewLocalizer(bundle, locale...)

	// templateData 치환
	for k, v := range templateData {
		newVal, _ := localizer.Localize(&_i18n.LocalizeConfig{
			MessageID:    v,
			TemplateData: nil,
		})
		templateData[k] = newVal
	}

	// 메시지 번역
	msg, err := localizer.Localize(&_i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: templateData,
	})

	if err != nil {
		return key
	}
	return msg
}
