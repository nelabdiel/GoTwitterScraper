package main
import (
    "fmt"
    "os"
    "encoding/json"
    "regexp"
	"time"
    "net/url"
    "strings"
	"net/http"

	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

    "github.com/ChimeraCoder/anaconda"
	"github.com/yhat/scrape"
    "golang.org/x/net/html/atom"
	"golang.org/x/net/html"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

//Twitter credentials
var secrets struct {
    ApiKey  string  `json:"apiKey"`
    ApiSecret  string `json:"apiSecret"`
    AccessToken  string `json:"accessToken"`
    AccessTokenSecret string `json:"accessTokenSecret"`
}



// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("calendar-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}





//Main function that extracts links from tweets.
func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("/var/google.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	//srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	calendarService, err := calendar.New(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Opening json with twitter credentials.
    twitterSecrets, err := os.Open("/var/tweeter.json")
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
    v.Set("count", "20")
	//Search for tweets that contain charleston and event.
    result, err := api.GetSearch("charleston AND event -RT", v)
	//Go through tweets
    for _, tweet := range result.Statuses {
	//Check if there is a link in tweet.
	for _, link := range strings.Split(tweet.Text, " ") {
            hasLink := strings.Index(link, "http") == 0
            if hasLink {

				resp, err := http.Get(link)

                // Parse the page
                root, err := html.Parse(resp.Body)
                if err != nil {
                    // handle error
                }
                // Search for the title
                title, ok := scrape.Find(root, scrape.ByTag(atom.Title))
                if ok {

					//Lets get the details of the event.
					//Title
					scrapedTitle := scrape.Text(title)

					//Regex for parsing date
					months1 := "(Jan(uary)?|Feb(ruary)?|Mar(ch)?|Apr(il)?|May|Jun(e)?"
					months2 := "|Jul(y)?|Aug(ust)?|Sep(tember)?|Oct(ober)?|Nov(ember)?|Dec(ember)?)"
					dayYear := `\s+\d{1,2},\s+\d{4}`
					dateFormat := months1 + months2 + dayYear
					//Find date
					r, _ := regexp.Compile(dateFormat)
					date := r.FindString(scrapedTitle)

					//if there is a date print title, date and link.
					if len(date) > 0 {
						//Print title.
						fmt.Println(scrapedTitle)
						const longForm = "January 2, 2006"
						t, _ := time.Parse(longForm, date)
						//Print time
						fmt.Println(t.Format(time.RFC3339))
						//Print link
						fmt.Println(link)

						//Details of our new event.
						newEvent := calendar.Event{
							Summary: scrapedTitle + " " + link,
							Start: &calendar.EventDateTime{DateTime: t.AddDate(0,
								0, 0).Format(time.RFC3339)},
							End: &calendar.EventDateTime{DateTime: t.AddDate(0,
								0, 1).Format(time.RFC3339)},
						}
						//Inserting event to Google Calendar.
						//calendarService.Events.Insert("primary", &newEvent)

						createdEvent, err := calendarService.Events.Insert("primary", &newEvent).Do()
						if err != nil {
							fmt.Println(err)
							return
						}
						//Print confirmation that it has been added.
						fmt.Println("New event in your calendar: \"", createdEvent.Summary, "\" starting at ",
							createdEvent.Start.DateTime)


				}
                }


            }
        }
    }
}

