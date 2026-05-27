




# 1: Prime Time

To keep costs down, a hot new government department is contracting out its mission-critical primality testing to the lowest bidder. (That's you).

Officials have devised a JSON-based request-response protocol. Each request is a single line containing a JSON object, terminated by a newline character ('\n', or ASCII 10). Each request begets a response, which is also a single line containing a JSON object, terminated by a newline character.

After connecting, a client may send multiple requests in a single session. Each request should be handled in order.

A conforming request object has the required field method, which must always contain the string "isPrime", and the required field number, which must contain a number. Any JSON number is a valid number, including floating-point values.

Example request:

{"method":"isPrime","number":123}
A request is malformed if it is not a well-formed JSON object, if any required field is missing, if the method name is not "isPrime", or if the number value is not a number.

Extraneous fields are to be ignored.

A conforming response object has the required field method, which must always contain the string "isPrime", and the required field prime, which must contain a boolean value: true if the number in the request was prime, false if it was not.

Example response:

{"method":"isPrime","prime":false}
A response is malformed if it is not a well-formed JSON object, if any required field is missing, if the method name is not "isPrime", or if the prime value is not a boolean.

A response object is considered incorrect if it is well-formed but has an incorrect prime value. Note that non-integers can not be prime.

Accept TCP connections.

Whenever you receive a conforming request, send back a correct response, and wait for another request.

Whenever you receive a malformed request, send back a single malformed response, and disconnect the client.

Make sure you can handle at least 5 simultaneous clients.



## Notes

### request
`{"method":"isPrime","number":123}`

- one json object
- termiated by newline
- must contain field "method" which always contains isPrime
- must contain a number in "isPrime"
- any json number is a valid number *also floating point* (non-integers can not be prime)
#### request is malformed if
- not well formed json
- any requed fild is missing
- method name is not "isPrime"
- number value is not a number

> ignore all other fields if present

### response
`{"method":"isPrime","prime":false}`

- one json object
- termiated by newline
- must contain field "method" always contains isPrime
- most contain filed "prime" which contains boolean 
  if the number in request was prime or not
#### response is malformed if 
- not well-formed JSON
- any required field is missing
- method name is not "isPrime"
- prime value is not boolean
  
client connects
- can have multiple requests in single session
- handle in order
- if conforming send response and wait for another request
- if malfomrmed send back single malformed response and close socket

5 clients




