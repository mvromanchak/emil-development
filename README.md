
### What is this repository for? ###
Implementation of a gRPC service that receives data from GPS devices. Service
 get data and save it to the database. GPS device uses REST API.

Targers in Makefile:
=
1. `make local-down`
2. `make local-up`
3. `make local-build`

## Future improvements:
- Logging can be added easily with go-kit internal lib.
- Tracing with Zipkin or DataDog also can be added in go-kit middleware.
- Hosting must be based on  Docker, orchestration can add load balancing and scalability.
- Error custom handling is very important in a language where is good practice to return an error from every function.
- Caching and cold/hot data are very important for good performance.
- My solution based on clean architecture and go-kit can be easily extended. Data, transport, endpoints layers separated.
- Replace sdk folder in specific repo.
# API Documentation:

### JWT struct
```json
{
  "sub": "1234567890",
  "devise_id": "T-800"
}
```

### GPS Exchange Format 
body for request.
```xml
curl --location --request POST 'http://localhost:5686/api/v1/gps' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZGV2aXNlX2lkIjoiVC04MDAifQ.ykWR7C-lIL_tnauX16IxoJzD0pQ3rmGg_wKj3_Q7ReE' \
--header 'Content-Type: application/xml' \
--data-raw '<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" creator="T-800" version="1.1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">
<metadata>
<link href="http://www.garmin.com">
<text>Garmin International</text>
</link>
<time>2009-10-17T22:58:43Z</time>
</metadata>
<trk>
<name>Example GPX Document</name>
<trkseg>
<trkpt lat="47.644548" lon="-122.326897">
<ele>4.46</ele>
<time>2009-10-17T18:37:26Z</time>
</trkpt>
<trkpt lat="47.644548" lon="-122.326897">
<ele>4.94</ele>
<time>2009-10-17T18:37:31Z</time>
</trkpt>
<trkpt lat="47.644548" lon="-122.326897">
<ele>6.87</ele>
<time>2009-10-17T18:37:34Z</time>
</trkpt>
</trkseg>
</trk>
</gpx>'

```
## Built With:
- "github.com/go-kit" - best tool-kit.
- "github.com/pkg/errors" - error wrapper.