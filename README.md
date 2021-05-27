# Microsoft Text Translation
[API Docs](https://docs.microsoft.com/en-US/azure/cognitive-services/translator/)  
Register for Microsoft Text Translation API ([see instructions](https://docs.microsoft.com/en-us/azure/cognitive-services/translator/translator-how-to-signup))  
Use the obtained subscription key to instantiate a translator as shown below.  

### Install:
```
go get -u github.com/snakesel/mstranslator
```

### Example usage:

```go
package main

import (
    "fmt"
    ms "github.com/snakesel/mstranslator"
)

func main() {
    translate := ms.New(ms.Config{
        Key:    "YOUR-SUBSCRIPTION-KEY",
        Region: "YOUR_RESOURCE_LOCATION",
    })

    // you can use "" for source language
    // so, translator will detect language
    trtext, err := translate.Translate("Hello, World!", "", "ru")
    if err == nil {
        fmt.Println(trtext)
    } else {
        fmt.Println(err.Error())
    }

    // Detect the language of the text
    conf, lang, err := translate.Detect("Nächster Stil")
    if err == nil {
        fmt.Printf("%s (%f)\n", lang, conf)
    } else {
        fmt.Println(err.Error())
    }
}
```
