package locales

// DefaultLanguage represents the default language
var DefaultLanguage = en

// AllLanguages represents all the supported language
// To add new languages, please refer to https://ezbookkeeping.mayswind.net/translating
var AllLanguages = map[string]*LocaleInfo{
	"de": {
		Content: de,
	},
	"en": {
		Content: en,
	},
	"es": {
		Content: es,
	},
	"fr": {
		Content: fr,
	},
	"it": {
		Content: it,
	},
	"ja": {
		Content: ja,
	},
	"nl": {
		Content: nl,
	},
	"pt-BR": {
		Content: ptBR,
	},
	"ru": {
		Content: ru,
	},
	"uk": {
		Content: uk,
	},
	"vi": {
		Content: vi,
	},
	"zh-Hans": {
		Content: zhHans,
	},
	"zh-Hant": {
		Content: zhHant,
	},
}
