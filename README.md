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

## Swagger endpoint
<app_url>/swagger/index.html

## Example
- https://covid19-go-deploy.herokuapp.com/cases?latitude=26.92&longitude=75.82
- https://covid19-go-deploy.herokuapp.com/swagger/index.html


## Useful Links while developing

### Golang echo mongodb sample:
- [Starting point](https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-echo-version-2gdg)

### Swagger:
- [echo-swagger](https://github.com/swaggo/echo-swagger)
- [echo-swagger tutorial](https://medium.com/geekculture/tutorial-generate-swagger-specification-and-swaggerui-for-echo-go-web-framework-3ac33afc77e2)

### Redis:
- [redis commands](https://redis.io/commands/)
- [Redigo helpers — AddFlat and ScanStruct](https://itnext.io/storing-go-structs-in-redis-using-rejson-dab7f8fc0053)
- [Getting started with redigo](https://developer.redis.com/develop/golang/)