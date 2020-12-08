Web services allow different parts of an application to communicate over a network,
that serves as the searation of concerns, allowing each part of an application to focus on its own responsibility

REST APIs are the most common
- work over HTTP
- client = website, another web service, mobile app
- uses Uniform Resource Identifiers (URI) (end points)
 - in order to map to the resources 
- has a response and a request
- stateless : the requests contain all the information that is needed to get a response, that way the service doesn't need to preserve any information or know anything about the client, this also allows for scalibilityalibility

Each http address/URI/pattern has to have a handler that handles the request on that pattern and decides what to do with it.
ServeMux matches up the request with the handlers that most closely represent the pattern that the request was made on.
To listen for incoming requests: http.ListenAndServe() - accepts the port to be listened on and an instance of ServeMux (nil represents the default ServeMux)

ListenAndServeTLS - if you need to handle traffic over encrypted http, allows to configure a certificate

Request object
 - all request have an http method which we can access (Request.Method (string))
 - Request.Header(map[string][]string)
 - Request.Body of type io.ReadCloser (implements Reader{} interface), if the body is empty - EOF is returned

URL path parameters
 - 1. Match the url to their closest static path from the handler
 - 2. Parse out the path parameters from the Url in order to use them 
 - Request.URL (Scheme, Opaque, User, Host, Path (!), RawPath, ForceQuery, RawQuery, Fragment

scheme://host[:port#]/path/.../[?query-string][#anchor]
scheme         assign underlying protocol (such as HTTP, HTTPS, FTP)
host           IP or domain name of HTTP server
port#          default port is 80, and it can be omitted in this case.
        If you want to use other ports, you must specify which port. For example,
        http://www.cnblogs.com:8080/
path           resources path
query-string   data are sent to server
anchor         anchor

Middleware
 - includes authentication, logging, session management
 - executed either before or after the handlers are called
 - one way to do this is to wrap our handler in one of the special adapter functions
func middlewareHandler (handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//do stuff before intended handler here
		handler.ServeHttp(w, r)
		//do stuff after intended handler here
	})
}

func intendedFunction(w http.ResponseWriter, r *http.Request) {
	//business logic here
}

func main() {
intendedHandler := http.HandlerFunc(intendedFunction)
http.Handle("/foo", middlewareHandler(intendedHandler))
http.ListenAndServe(":5000", nil)
}

Teacher advices the middleware func to return and accept as a paramater handler of type http.Handler
and convert HandlerFunc to Handle in main, to be able to pass it to the middleware function, then do http.Handle(route, middlewarefunc with the handler param)

but instead I left the main as it was (http.HandleFunc(route, middlewarefunc) and in the middleware function the type I accept and return is http.HandlerFunc, seemed like less of a headache to make it that way...


CORS (Cross-origin Resource Sharing)
 - built-in all the browsers, protects from cross-origin attacks
 - prevents the website to access resources at another origin
http://globomantics.com
+http://globomantic.com/products, can retrieve products from every path on that domain
+http://globomantic.com/api/products/123
-http://globomantic.com:8080/products - different port on the same server, wouldn't work
-https://globomantic.com/products - same with the change of protocol
-https://dev.globomantic.com/dashboard - or subdomains (by default)

checkout mozilla's resource page on CORS

Websockets
- are needed in case the data on the server changes, the client needs to know about that too - they need to somehow communicate that
- uses include: video games, chats, stock trading platforms

