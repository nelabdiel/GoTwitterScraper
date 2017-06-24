# Building an events Twitter scraper to familiarize myself with Go.

Done: As it is right now it looks for Tweets that mentions the words _Charleston_ and _Event_ 
and that are not retweets and it extracts the links.

Done: The next step will be to parse the links and extract information about the events.

_It extracts the title of the page which often contains the information needed._

To Do: The last step will be to update a calendar to add the events.



The Binary was built for darwin64. Make sure to add a json file _tweeter.json_ to /var/ 

_tweeter.json_ should look like this:

{ "apiKey": "XXXX", "apiSecret": "XXXX", "accessToken": "XXXX", "accessTokenSecret": "XXXX" }
