![Go](https://github.com/jasondborneman/solar3/workflows/Go/badge.svg?branch=main)

# solar3
Newest version of code that combines various APIs to gather data about my home solar array behind a SolarEdge inverter. Tweets a subset of the data along with a graph of Power Generated and Sun Altitude every 15 minutes. Runs as a GCP Cloud Function Data gathered includes:

* Power Generated: in 15 minute chunks
* DateTime: the time of the SolarEdge power reading
* DateTimeStored: the time that the data was stored in the database
* CloudCover: % Cloud Coverage as determined by Lat/Long via OpenWeather API
* Humidity: % Humidity as determined by Lat/Long via OpenWeather API
* Pressure: Air pressure as determined by Lat/Long via OpenWeather API
* RainPastHour: Rain fall (if any) in teh past hour as determined by Lat/Long via Openweather API
* SnowPastHour: Snow fall (if any) in the past hour as determined by Lat/Long via OpenWeather API
* Temp: (F) as determined by Lat/Long via OpenWeather API
* WeatherID: The ID of the current weather conditions as determined by Lat/Long via OpenWeather API
* SunAltitude: As determined by Lat/Long via the IPGeolocation API
* SunAzimuth: As determiend by Lat/Long via the IPGeolocation API

## APIs Used
* Twitter: https://github.com/jasondborneman/go-twitter (a fork of https://github.com/drswork/go-twitter, in turn a fork of https://github.com/dghubble/go-twitter. Forks include some handling for posting media to tweets)
* SolarEdge API: https://www.solaredge.com/sites/default/files//se_monitoring_api.pdf
* OpenWeather API: https://openweathermap.org/current
* IPGeolocation Astronomy API: https://ipgeolocation.io/documentation/astronomy-api.html

## Environment Variables Needed
* SOLAREDGE_APIKEY : SolarEdge API Key can be found in your SolarEdge Dashboard Admin>Access page.
* SOLAREDGE_SITEID : Part of your SolarEdsge Dashboard URL
* IPGEOLOCATION_APIKEY : API Key for IPGeolocation.io
* SOLAR3_LATITUDE : Latitude of your solar array
* SOLAR3_LONGITUDE : Longitude of your solar array
* OPENWEATHER_APIKEY : API Key for OpenWeather
* GCP_PROJECT : The name of your GCP Project
* TWITTER_CONSUMERKEY : Twitter Consumer Key
* TWITTER_CONSUMERSECRET : Twitter Consumer Secret
* TWITTER_ACCESSTOKEN : Twitter Access Token
* TWITTER_ACCESSSECRET : Twitter Access Secret
* STUPID_AUTH : Make up your own GUID or any string really.

### Optional Environment Variables
* DO_TWEET : [true|false] to turn on and off tweeting. Defaults to false
* DO_SAVEGRAPH : [true|false] to turn on and off saving a local copy of the graph

## Auth
Set a local environment variable wherever your function is running called STUPID_AUTH with a value of some string key (such as a GUID). When calling the function endpoint, the body should be:

```
{
  "stupidAuth":"<your key here>"
}
```

As the name suggests, auth is really stupid. It just checks that the body key matches the Environment Variable.

## Tech Stack

* Function is written in Go.
* Data is stored in a GCP Firestore database
* Runs as a GCP Cloud Function
* Triggered by a GCP Cloud Scheduler Job
* Stored and CI/CD in Github 
* Also stored in a linked GCP Cloud Source Repository

## Deployment

It's currently set up using a GitHub CI/CD Action. This will be fleshed out further to include builds of all PRs to `main` in addition to the Build/Deploy of all commits to `main`.

It's also set up to deploy from a local command line. Set up a linked GCP Cloud Source Repository for the GitHub repo (or just store all your code there and skip GitHub altogether). The following shell command will deploy from the Cloud Source Repository to the Cloud Function:

```
gcloud functions deploy solar3_github \
  --source https://source.developers.google.com/projects/<project id>repos/<cloud source repo name>/moveable-aliases/<branch>/paths/ \
  --runtime go113 \
  --trigger-http \
  --allow-unauthenticated
```

it is recommended, if you're going to use the GitHub Action to deploy, to set up the Function first time on the command line, as GitHub Action for deploying to GCP Cloud Functions doesn't allow unauthenticated functions to be deployed from what I can tell, but it won't mess with the setting if it's already there.

## Running locally

Rename function.go -> function.go.old
Rename main.go.old -> main.go
```
$ go build .
$ ./solar3
```

## Future Ideas:
* Fix the graph to show the secondary Y axis for Sun Altitude.
* Tie into machine learning to find what data best correlates with high power generation. It's probably pretty obvious (sun altitude & cloud cover) but it'd be interesting to try.

## TODO:
* Unit tests - still need to figure out mocking in Go to do this. This will probably mean some restructuring of the code to make it more testable. I was writing/learning as I went writing this.
