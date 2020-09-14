# ip2location

Detect Geo data by IP with ip2location database file

## Configurations

| Env Variable | Default | Description |
|--------------|---------|-------------| 
| IP2LOCATION_DATABASE | IP2LOCATION.BIN | Path to DB file |
| IP2LOCATION_HTTPLISTENER | 0.0.0.0:8001 | Listen address |
| IP2LOCATION_VERBOSE | false | Verbose mode for logging all requests to stderr |

## Usage

Send IP to the service with one way of:

- Query argument `ip`. Example: `GET http://domain.com?ip=1.2.3.4`
- HTTP header `X-IP2LOCATION-IP`
- In POST request body  

Result - json from official library `github.com/ip2location/ip2location-go`:

```
type IP2Locationrecord struct {
	Country_short      string
	Country_long       string
	Region             string
	City               string
	Isp                string
	Latitude           float32
	Longitude          float32
	Domain             string
	Zipcode            string
	Timezone           string
	Netspeed           string
	Iddcode            string
	Areacode           string
	Weatherstationcode string
	Weatherstationname string
	Mcc                string
	Mnc                string
	Mobilebrand        string
	Elevation          float32
	Usagetype          string
}
```
