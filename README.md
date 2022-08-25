# GOLANG HTTP STATUS CHECKER SERVER
Golang based multithreaded server to check the website status as provided by a client.

## Compilation and Execution
1. Open Terminal at home directory of project.
2. Begin server using 
    go run main.go
3. Set the execution URL at main.go, line 12. (set to localhost:8080/websites by default).
4. Send list of websites to the server using post request of the following format.
    curl -X POST URL/website -d {"websites": {"hostname1", "hostname2"....}}
    curl -X POST localhost:8080/websites -d '{"websites": ["google.com", "yahoo.com", "abcd.com"]}'
5. Call get request with name to get status of a particular list, else to get list of all websites.
    curl URL/website                                          // For listing out all websites
    curl URL/website?name=hostname1&name=hostname2            // For particular websites
6. Update the list with another post request if necessary else continue,

## Testing bash
1. Open Terminal at home directory of project.
2. Run tests using 
    sh tests.sh