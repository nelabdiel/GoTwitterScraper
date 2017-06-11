package main
import (
    "fmt"
    "os"
    "encoding/json"
    "github.com/ChimeraCoder/anaconda"
)

var secrets struct {
    ApiKey  string  `json:"apiKey"`
    ApiSecret  string `json:"apiSecret"`
    AccessToken  string `json:"accessToken"`
    AccessTokenSecret string `json:"accessTokenSecret"`
}



func main() {
    configFile, err := os.Open("../variables.json")
    if err != nil {
        fmt.Println(err.Error())
    }

    jsonParser := json.NewDecoder(configFile)
    if err = jsonParser.Decode(&secrets); err != nil {
        fmt.Println(err.Error())
    }
    anaconda.SetConsumerKey(secrets.ApiKey)
    anaconda.SetConsumerSecret(secrets.ApiSecret)
    api := anaconda.NewTwitterApi(secrets.AccessToken,secrets.AccessTokenSecret)
    searchResult, _ := api.GetSearch("charleston AND event", nil) 
    for _ , tweet := range searchResult.Statuses {fmt.Println(tweet.Text)}
}

