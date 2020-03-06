# Hogger
The purpose of this API is to aid in the exploitation of SSRF vulnerabilities and XSS vulnerabilities, but it 
can likely be used to assist in various exploits. 

The main features are the following:
+ Open redirect
+ Echo input back in response
+ Set headers to be returned
+ Set HTTP status to be returned 
+ Store and retrieve payloads 

All endpoints can be provided a `status` and `headers` parameter that will indicate headers and HTTP status code that should be sent 
in the response. These can be sent in either the Body or the URL.

The `status` parameter should be a valid HTTP status code.
The `headers` parameter should be a comma separated list of headers and values to be in the response. The headers 
themselves will be separated by a `:`. For example, you can send a request of 

```
GET /echo?echo=Hello+World!&status=404&headers=Content-Type:text/html,Access-Control-Allow-Origin:*
```

There are 4 endpoints. Parameters for the endpoints can be sent in the URL or the Body.

+ /echo

This endpoint takes a paramater of `echo`, which will simply echo back whatever is provided. 

+ /redirect

This endpoint takes a parameter of `url`, which will send a redirect to the specified URL. If you provide a 
status code, it should be in the 3XX range, otherwise 302 will be chosen by default. 

+ /store and /p

The `/store` endpoint takes a parameter of `k` and `v`, representing a key and value pair. The key and value 
will be stored for 2 minutes, at which point they will be deleted to clean up. 

The `/p` endpoint takes a parameter of `k`, which will lookup a stored payload and return it if it is present. 

The purpose of this would likely be to quickly exploit a XSS vulnerability, or something similar where you don't 
have enough room for your payload and would like to load it from somewhere else. For example, you could send 
your payload to `/store?k=mypayload&v=AAAAAAAAAAAAAAAA<script>alert(1)</script>`, and then access the payload 
by requesting `/p?k=mypayload`. 

Notably, all of these endpoints will respond with an open CORS policy to `*` origins, although this can 
overwritten with user provided headers. 