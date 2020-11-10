
# datetime.link

Readable links to a point in time

## Credits

- [GeoNames](https://www.geonames.org/) and its contributors

## Reporting Problems

File an [issue](https://github.com/serverwentdown/datetime.link/issues/new) on 
GitHub. 

## Rationale

Sometimes, you want to refer to a point in time, and also want to provide a 
timezone converter for them. datetime.link provides links to points in time and 
presents them in a set of selected timezones and/or local time.

## Compatibility

To ensure compatibility with non-JavaScript clients, the server responds with a 
pure HTML page without local time. JavaScript provides the rendering of local 
time and local time format using the native [Intl](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DateTimeFormat/DateTimeFormat)
and [Date](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date)
APIs, eventually the [Temporal API](https://github.com/tc39/proposal-temporal). 
JavaScript also provides the visual editing interface for the date, time and 
list of timezones to show in.

## Timezone Specifiers

[GeoNames](https://www.geonames.org/) data is used to generate 
[`data/cities.json`](https://datetime.link/data/cities.json), a huge 5MB JSON 
blob containing cities with a population greater than 15000. In the event this 
is not enough, more cities can be included by using the alternate 
`cities5000.txt` file. 

The timezone specifiers are generated by the code in `scripts/data.go`, which 
assemble an identifier replacing non-alphanumeric and quote characters with 
underscores. Dashes then are used to join the city name, administrative division
level 1 names and country names. 

Alternatively, a fixed timezone can be specified as an offset like `+08:00`. 
This caters for scenarios where the local DST and other local time differences
should not be accounted for. These are guarenteed to be stable and accurate. 

## Timezone Data

Go [relies on local tzdata](https://golang.org/pkg/time/#LoadLocation), and thus
`datetime.link` relies on local tzdata. 

## Upcoming Improvements

See [Issues](https://github.com/serverwentdown/datetime.link/issues).

