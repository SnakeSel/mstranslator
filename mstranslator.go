package mstranslator

import (
	"bytes"
	"encoding/json"

	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

// Config is the configuration struct you should pass to New().
type Config struct {
	Url    string
	Key    string
	Region string
	// Debug is an optional writer which will be used for debug output.
	Debug io.Writer
}
type Translation struct {
	log *log.Logger
	Config
}

// New returns a new Translation.
func New(conf Config) *Translation {

	tr := new(Translation)

	if conf.Url != "" {
		tr.Url = conf.Url
	} else {
		tr.Url = "https://api.cognitive.microsofttranslator.com"
	}

	tr.Key = conf.Key

	tr.Region = conf.Region

	if conf.Debug == nil {
		conf.Debug = ioutil.Discard
	}

	tr.log = log.New(conf.Debug, "[MSTrans]\t", log.LstdFlags)

	return tr
}

// Translate text from a language to another
func (tr *Translation) Translate(source, sourceLang, targetLang string) (string, error) {

	// Build the request URL. See: https://golang.org/pkg/net/url/#example_URL_Parse
	uri, err := url.Parse(tr.Url)
	if err != nil {
		tr.log.Println("Error parse url")
		return "", err
	}
	uri.Path = path.Join(uri.Path, "/translate")

	q := uri.Query()
	if len(sourceLang) > 0 {
		tr.log.Println("add \"from\" to param")
		q.Add("from", sourceLang)
	}
	q.Add("to", targetLang)
	q.Add("api-version", "3.0")
	uri.RawQuery = q.Encode()

	// Create an anonymous struct for your request body and encode it to JSON
	body := []struct {
		Text string
	}{
		{Text: source},
	}
	b, _ := json.Marshal(body)

	tr.log.Println(uri.String())

	// Build the HTTP POST request
	req, err := http.NewRequest("POST", uri.String(), bytes.NewBuffer(b))
	if err != nil {
		tr.log.Println("Error NewRequest")
		return "", err
	}
	// Add required headers to the request
	req.Header.Add("Ocp-Apim-Subscription-Key", tr.Key)
	req.Header.Add("Ocp-Apim-Subscription-Region", tr.Region)
	req.Header.Add("Content-Type", "application/json")

	// Call the Translator API
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		tr.log.Println("Post error")
		return "", err
	}

	tr.log.Printf("Response code: %d", res.StatusCode)

	switch res.StatusCode {
	case http.StatusOK: //200
		// Decode the JSON response
		type ResponseJSON struct {
			DetectedLanguage struct {
				Language string  `json:"language"`
				Score    float32 `json:"score"`
			} `json:"detectedLanguage"`
			Translations []struct {
				Text string `json:"text"`
				To   string `json:"to"`
			} `json:"translations"`
		}

		tr.log.Println("Decode json")
		var result []ResponseJSON
		if err := json.NewDecoder(res.Body).Decode(&result); err == nil {

			if len(result) == 1 {
				tr.log.Println(result[0])
				return result[0].Translations[0].Text, nil
			} else {
				return "", fmt.Errorf("unknown number of responses ")
			}
		}
		//tr.log.Println("fallthrough")
		//fallthrough //go next case
	//case StatusBadRequest: //400 	Invalid request
	//case StatusTooManyRequests: //429 Slow down
	//case StatusInternalServerError: //500 Detection error
	default:
		// Decode the JSON response
		// var result2 interface{}
		// if err := json.NewDecoder(res.Body).Decode(&result2); err != nil {
		// 	return "", err
		// }

		// // Format and print the response to terminal
		// prettyJSON, _ := json.MarshalIndent(result2, "", "  ")
		// fmt.Printf("%s\n", prettyJSON)
		type ResponseJSON struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		var result ResponseJSON
		if err := json.NewDecoder(res.Body).Decode(&result); err == nil {
			return "", fmt.Errorf("%d: %s", result.Error.Code, result.Error.Message)

		}
		// m := result2.(map[string]interface{})
		// if val, ok := m["error"]; ok {
		// 	respErr := val.(map[string]interface{})
		// 	return "", fmt.Errorf("%v: %v", respErr["code"], respErr["code"])
		// }

	}

	return "", errors.New("unknown answer")

}

// Detect the language of the text
// Return: confidence, language, error
func (tr *Translation) Detect(text string) (float32, string, error) {
	// Build the request URL. See: https://golang.org/pkg/net/url/#example_URL_Parse
	uri, err := url.Parse(tr.Url)
	if err != nil {
		tr.log.Println("Error parse url")
		return -1, "", err
	}
	uri.Path = path.Join(uri.Path, "/detect")

	q := uri.Query()
	q.Add("api-version", "3.0")
	uri.RawQuery = q.Encode()

	// Create an anonymous struct for your request body and encode it to JSON
	body := []struct {
		Text string
	}{
		{Text: text},
	}
	b, _ := json.Marshal(body)

	tr.log.Println(uri.String())

	// Build the HTTP POST request
	req, err := http.NewRequest("POST", uri.String(), bytes.NewBuffer(b))
	if err != nil {
		tr.log.Println("Error NewRequest")
		return -1, "", err
	}
	// Add required headers to the request
	req.Header.Add("Ocp-Apim-Subscription-Key", tr.Key)
	req.Header.Add("Ocp-Apim-Subscription-Region", tr.Region)
	req.Header.Add("Content-Type", "application/json")

	// Call the Translator API
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		tr.log.Println("Post error")
		return -1, "", err
	}

	tr.log.Printf("Response code: %d", res.StatusCode)

	switch res.StatusCode {
	case http.StatusOK: //200
		// Decode the JSON response
		type Alternative struct {
			IsTranslationSupported     bool    `json:"isTranslationSupported"`
			IsTransliterationSupported bool    `json:"isTransliterationSupported"`
			Language                   string  `json:"language"`
			Score                      float32 `json:"score"`
		}
		type ResponseJSON struct {
			Alternatives               []Alternative `json:"alternatives"`
			IsTranslationSupported     bool          `json:"isTranslationSupported"`
			IsTransliterationSupported bool          `json:"isTransliterationSupported"`
			Language                   string        `json:"language"`
			Score                      float32       `json:"score"`
		}

		tr.log.Println("Decode json")
		var result []ResponseJSON
		if err := json.NewDecoder(res.Body).Decode(&result); err == nil {

			if len(result) == 1 {
				tr.log.Println(result[0])
				return result[0].Score, result[0].Language, nil
			} else {
				return -1, "", fmt.Errorf("unknown number of responses ")
			}
		}
	default:
		type ResponseJSON struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		var result ResponseJSON
		if err := json.NewDecoder(res.Body).Decode(&result); err == nil {
			return -1, "", fmt.Errorf("%d: %s", result.Error.Code, result.Error.Message)

		}
	}

	return -1, "", errors.New("unknown answer")
}
