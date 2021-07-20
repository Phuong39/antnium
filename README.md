# antanium 

```
Anti Tanium
```

There are two components: 
* Client: The actual trojan
* Server: C2 infrastructure 


## How to use

Configure your campaign in: 
* model/campaign.go 

Compile client and server: 
* make 

And deploy them. 


## Client 

Commands: 
* exec: Execute a file
* fileupload: upload a file 
* filedownload: download a file 
* dir: directory content
* test: test (e.g. unit- or integration tests)
* ping: announce we are still alive, send info, get new commands (sent without request from server)

## Server

* Runs on a specific port
* uploads files from client via REST to `./upload/`
* serves directory `./static/`


## Security 

The client and server share a static encryption key, and a API key. 

If the blue team manages to extract the API key from a HTTP proxy or client binary, they
gain access to the server API, which enables them to:
* flood the server with fake clients 
* observe public executed commands easily 

If the blue team manages to extract the encryption key from a client binary, they can: 
* decrypt all past communications of all client instances (if they have proxy log)
* Issue new commands to existing clients (if they can perform HTTP MITM on proxy)

This is intentional. The campgain is only protected against outsiders, not a motivated blue team. 


## Testing

```
go test ./...
go test ./server
```