package api

import (
	"PapagoGo/powershell"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getNaverClientId() string {
	return os.Getenv("X-Naver-Client-Id")
}

func getNaverClientSecret() string {
	return os.Getenv("X-Naver-Client-Secret")
}

func translate(src LangCode, tgt LangCode, text string) string {
	clientId := getNaverClientId()
	clientSecret := getNaverClientSecret()

	params := url.Values{}
	params.Add("source", string(src))
	params.Add("target", string(tgt))
	params.Add("text", text)

	requestURL := translationUrl + "?" + params.Encode()

	req, _ := http.NewRequest("POST", requestURL, nil)
	req.Header.Set("User-Agent", "curl/7.49.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Naver-Client-Id", clientId)
	req.Header.Set("X-Naver-Client-Secret", clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting response from Papago translation API: %+v", err)
		return ""
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %+v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading the response body from the API: %v", err)
		return ""
	}

	var apiResp TranslationResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return ""
	}

	return apiResp.Message.Result.TranslatedText
}

func detectLang(text string) LangCode {
	clientId := getNaverClientId()
	clientSecret := getNaverClientSecret()

	params := url.Values{}
	params.Add("query", text)

	requestURL := languageDetectionUrl + "?" + params.Encode()

	req, _ := http.NewRequest("POST", requestURL, nil)
	req.Header.Set("User-Agent", "curl/7.49.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Naver-Client-Id", clientId)
	req.Header.Set("X-Naver-Client-Secret", clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting response from Papago translation API: %+v", err)
		return ""
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %+v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading the response body from the API: %v", err)
		return ""
	}

	var apiResp DetectResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return ""
	}

	return apiResp.LangCode
}

func copyLatestTranslationToClipboard(lastTranslation string) {
	if lastTranslation == "" {
		fmt.Printf("%sThere is no latest tranlation yet%s\n", powershell.ColorRed, powershell.ColorReset)
		return
	}
	err := clipboard.WriteAll(lastTranslation)
	if err != nil {
		log.Printf("Error writing the latest translation into clipboard: %+v", err)
		return
	}
	fmt.Printf("%sCopied the latest translation into clipboard%s\n", powershell.ColorMagenta, powershell.ColorReset)
}

func GetUserInputAndTranslate() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%sEnter text to translate%s\n", powershell.ColorCyan, powershell.ColorReset)
	fmt.Printf("%spress Ctrl+C to exit%s\n", powershell.ColorYellow, powershell.ColorReset)
	fmt.Printf("%senter Ctrl+A to copy the latest translation into clipboard%s\n", powershell.ColorYellow, powershell.ColorReset)

	var lastTranslation string
	for {
		fmt.Println("")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		} else if input == "\x01" {
			copyLatestTranslationToClipboard(lastTranslation)
			continue
		}

		src, tgt := handleLangCodes(input)

		lastTranslation = translate(src, tgt, input)
		fmt.Printf("%s%s%s\n", powershell.ColorGreen, lastTranslation, powershell.ColorReset)
	}
}

func handleLangCodes(input string) (LangCode, LangCode) {
	src := detectLang(input)
	if src == Unknown {
		return Unknown, Unknown
	}

	if src == Korean {
		return src, English
	}
	return src, Korean
}
