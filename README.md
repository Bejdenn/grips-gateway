# grips-gateway

grips-gateway is a web handler that acts as a gateway for incoming GRIPS course searches. This service is useful when
combined with applications like Alfred or others, that support setting web query templates.

## How is it working?

The service processes incoming HTTP requests that contain a so-called 'hint'. A hint is a keyword that is associated
with a course on GRIPS. This makes it easy to open a course without typing in a full name.

For example, a course 'Biophysics' could have hints like 'bio', 'physics', or an abbreviation like 'BP'.

## Browsers and HTTP status code 301

grips-gateway responds to users with the HTTP status code '301: Moved Permanently'. Browsers cache the destination
URL that the server has responded with. If you now would change a hint and wouldn't refresh/delete your browser's
cache, the browser would still try to open the original course page associated with the hint.
