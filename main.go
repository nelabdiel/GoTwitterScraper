package main
import (
    "fmt"
    "os"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/ChimeraCoder/anaconda"
)

var secrets struct {
    ApiKey  string  `json:"apiKey"`
    ApiSecret  string `json:"apiSecret"`
    AccessToken  string `json:"accessToken"`
    AccessTokenSecret string `json:"accessTokenSecret"`
}

func main() {
    twitterSecrets, err := os.Open("../variables.json")
    if err != nil {
        fmt.Println(err.Error())
    }

    jsonParser := json.NewDecoder(twitterSecrets)
    if err = jsonParser.Decode(&secrets); err != nil {
        fmt.Println(err.Error())
    }
    anaconda.SetConsumerKey(secrets.ApiKey)
    anaconda.SetConsumerSecret(secrets.ApiSecret)
    api := anaconda.NewTwitterApi(secrets.AccessToken,secrets.AccessTokenSecret)

    v := url.Values{}
    v.Set("count", "30")
    result, err := api.GetSearch("charleston AND event -RT", v)
    for _, tweet := range result.Statuses {
	for _, link := range strings.Split(tweet.Text, " ") {
            hasLink := strings.Index(link, "http") == 0
            if hasLink {
                fmt.Println(link)
            }
        }
    }
}

