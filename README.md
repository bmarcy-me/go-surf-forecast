# LR GO SURF FORECAST

Find the best surf spots for next 7 days around La Rochelle, France.\
Backend implementation in Go.\
This project is a work in progress.

## Prerequisites
- Go 1.23 installed

## Start
```
go run cmd/main.go
```

## API endpoints

### /spots
/spots return the forecast for surf spots around La Rochelle\
Available query parameters :
- start=2024-10-12T08:00:00Z (iso dateTime between 11/10/2024 and 20/10/2024 because we use static data for now)
- duration=2 (from 1 to 7)

```
curl http://localhost:8080/api/spots/start=2024-10-12T08:00:00Z&duration=2
```

## Available surf spots
3 spots available for now 

|  Id   | Name                              |
| ----- | --------------------------------- |
| 1     | Plage de Gros Joncs - Ile de Ré   |
| 2     | Pointe du Lizay - Ile de Ré       |
| 3     | Plage de Vert Bois - Ile d'Oléron |



## To do list
- [ ] api endpoint returning the best surf spot and the best time to go there
- [ ] api endpoint returning the best surf spot for a given time
- [ ] querying stormglass at startup and store weather data in a db
- [ ] using the db instead of static json files
- [ ] docker for api server an db
- [ ] add surf spots around La Rochelle
- [ ] add tests