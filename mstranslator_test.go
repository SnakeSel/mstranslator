package mstranslator

import (
	"os"
	"testing"
)

func TestTranslate(t *testing.T) {
	translate := New(Config{
		Key:    MSkey,
		Region: MSregion,
		Debug:  os.Stdout,
	})

	text, err := translate.Translate("Hello, World!", "auto", "ru")

	if err == nil {
		if text != "Всем привет!" {
			t.Error("[TestTranslate] failed: Translate not valid.")
		}
	} else {
		t.Errorf("[TestTranslate] failed: %s", err.Error())
	}

}

func TestDetect(t *testing.T) {
	translate := New(Config{
		Key:    MSkey,
		Region: MSregion,
		Debug:  os.Stdout,
	})

	_, lang, err := translate.Detect("Nächster Stil")
	if err == nil {
		if lang != "de" {
			t.Error("[TestDetect] failed: language not valid.")
		}

	} else {
		t.Errorf("[TestDetect] failed: %s", err.Error())
	}

}
