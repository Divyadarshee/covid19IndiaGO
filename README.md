# COVID19 Cases in India in GOlang

## Deployment
Deployment in Heroku

## Sources
1. covid19-data: https://data.covid19india.org/v4/min/data.min.json
2. Free reverse Geocoding: https://locationiq.com/

## API endpoint
1. base url:  https://covid19-go-deploy.herokuapp.com/
2. get cases and vaccinations complete details using:
   1. Path "/cases"
   2. query parameters:
      1. Latitude (only within India)
      2. Longitude (only within India)


## Example
https://covid19-go-deploy.herokuapp.com/cases?latitude=26.92&longitude=75.82