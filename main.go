package main
import (
    "fmt"
    "os"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/ChimeraCoder/anaconda"
)

//Twitter credentials
var secrets struct {
    ApiKey  string  `json:"apiKey"`
    ApiSecret  string `json:"apiSecret"`
    AccessToken  string `json:"accessToken"`
    AccessTokenSecret string `json:"accessTokenSecret"`
}

//Main function that extracts links from tweets.
func main() {
	//Opening json with twitter credentials.
    twitterSecrets, err := os.Open("/Users/Nel/Desktop/Code/go/variables.json")
    if err != nil {
        fmt.Println(err.Error())
    }

    jsonParser := json.NewDecoder(twitterSecrets)
    if err = jsonParser.Decode(&secrets); err != nil {
        fmt.Println(err.Error())
    }
	//Loading credentials and setting up credentials.
    anaconda.SetConsumerKey(secrets.ApiKey)
    anaconda.SetConsumerSecret(secrets.ApiSecret)
    api := anaconda.NewTwitterApi(secrets.AccessToken,secrets.AccessTokenSecret)

    v := url.Values{}
    v.Set("count", "30")
	//Search for tweets that contain charleston and event.
    result, err := api.GetSearch("charleston AND event -RT", v)
	//Go through tweets
    for _, tweet := range result.Statuses {
		//Check if there is a link in tweet.
		for _, link := range strings.Split(tweet.Text, " ") {
            hasLink := strings.Index(link, "http") == 0
            if hasLink {
				//Print link.
                fmt.Println(link)
            }
        }
    }
}

