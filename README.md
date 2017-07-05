# An event finder script that uses the Twitter API to find events in the Charleston Area, 
parses the data and adds the events to your Google Calendar with the use of the Google API.

### Built in Go, the idea was for me to familiarized myself with how to write Go code. 

The Binary was built for darwin64. 

Make sure to add _tweeter.json_ and _google.json_ to /var/ 

Make sure to register your app with Twitter in order to get your credentials.

_tweeter.json_ should look like this:

{ "apiKey": "XXXX", "apiSecret": "XXXX", "accessToken": "XXXX", "accessTokenSecret": "XXXX" }

For _google.json_: 

Follow the instructions on step 1 from https://developers.google.com/google-apps/calendar/quickstart/go ,
rename the file google.json instead of client_secret.json and move it to /var/  
