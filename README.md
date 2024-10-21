# GO SURF FORECAST

Define your favorite surf spots.\
Get the surf spots condition for next 7 days based on stormglass API.\
Backend implementation in Go.\

**This project is a work in progress.**

## Prerequisites
- Go 1.23 installed
- Docker installed

## Surf spots configuration
Example with surf spots around La Rochelle, France.
You need to provide an id, a name, gps coordinates and the direction (angle) for the spot.

```yaml
spots:
  - id: 1
    name : "Plage de Gros Joncs - Ile de Ré"
    latitude: 46.1740867
    longitude: -1.3853837
    direction : 220
  - id: 2
    name : "Pointe du Lizay - Ile de Ré"
    latitude: 46.257935
    longitude: -1.518474
    direction : 320
  - id: 3
    name : "Plage de Vert Bois - Ile d'Oléron"
    latitude: 45.874214
    longitude: -1.263475
    direction : 260
```

## Start
In the root directory of the project, run the following commands:

```
export POSTGRES_HOST=postgres
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=surf_forecast
export WEATHER_DATA_SOURCE=file # other possible value is stormglass but you will need a valid stormglass api key (see "Stormglass" section in this README)
export STORMGLASS_API_KEY=xxx-yyy-zzz

cd docker
docker compose up --build -d
```

## API endpoints

### /spots
/spots return the forecast for surf spots

Available query parameters :
- start=2024-10-12T08:00:00Z (iso dateTime between 11/10/2024 and 20/10/2024 if you use static data)
- duration=2 (from 1 to 7)

```sh
curl -X GET http://localhost:8080/api/spots/start=2024-10-12T08:00:00Z&duration=2
```

The response contains each surf spots and the rating by hour with a score from 0 to 5

```json
{
    "spots": [
        {
            "id": 1,
            "name": "Plage de Gros Joncs - Ile de Ré",
            "ratings": [
                {
                    "rating": 2.221791666666667,
                    "time": "2024-10-12T09:00:00Z"
                },
                {
                    "rating": 2.3784027777777776,
                    "time": "2024-10-12T10:00:00Z"
                },
                ...
            ]
        },
        {
            "id": 2,
            "name": "Pointe du Lizay - Ile de Ré",
            "ratings": [
                {
                    "rating": 0.6341527777777778,
                    "time": "2024-10-12T09:00:00Z"
                },
                {
                    "rating": 0.9013472222222221,
                    "time": "2024-10-12T10:00:00Z"
                },
                ...
            ]
        },
        ...     
    ]
}
```

### /spots/best
/spots/best return the best surf spot and the best time to go there in the next X days from a start date

Available query parameters :
- start=2024-10-17T08:00:00Z (iso dateTime between 11/10/2024 and 20/10/2024 if you use static data)
- duration=4 (from 1 to 7)

```sh
curl -X GET http://localhost:8080/api/spots/best/start=2024-10-17T08:00:00Z&duration=4
```

The response contains only one surf spot. The one with the best rating and the best time to go there.

```json
{
    "id": 1,
    "name": "Plage de Gros Joncs - Ile de Ré",
    "ratings": [
        {
            "rating": 4.609416666666666,
            "time": "2024-10-20T19:00:00Z"
        }
    ]
}
```


## To do list
- [x] api endpoint returning the best surf spot and the best time to go there
- [x] querying stormglass at startup and store weather data in a db
- [ ] using the db instead of static json files
- [x] docker for api server an db
- [ ] add tests

## Stormglass
If you need a stormglass api key, create a free account on https://stormglass.io/


## Clean
To purge your docker env, in the root directory of the project, run the following commands:

```sh
cd docker
docker compose down -v --rmi all
```