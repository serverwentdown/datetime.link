
# WIP: datetime.link

Clean links to a point in time

## Rationale

Sometimes, you want to refer to a point in time, but want to provide a timezone converter for them. datetime.link provides links to points in time and presents them in localtime. You can also optionally select which timezones to present the time in.

## Compatibility

To achieve backwards compatibility and usability without JavaScript, the dates must be converted and rendered into pure HTML. This makes cURL a working target. 

Progressive enhancement extends on the UI by providing a visual editor for timezones, dates and times that change the URL. This is done with JavaScript and native input elements. In the future custom input elements could replace the native ones.
