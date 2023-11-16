package api

type TranslationResponse struct {
	Message struct {
		Type    string            `json:"@type"`
		Service string            `json:"@service"`
		Version string            `json:"@version"`
		Result  TranslationResult `json:"result"`
	} `json:"message"`
}

type TranslationResult struct {
	SrcLangType    string `json:"srcLangType"`
	TarLangType    string `json:"tarLangType"`
	TranslatedText string `json:"translatedText"`
	EngineType     string `json:"engineType"`
}

type DetectResponse struct {
	LangCode `json:"langCode"`
}

type LangCode string
