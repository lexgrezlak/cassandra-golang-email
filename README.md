# Features
```
curl -X POST localhost:8080/api/message -d '{"email":"john@gmail.com","title":"Hello World","content":"simple text","magic_number":101}'
```
Add message to the database


```
curl -X POST localhost:8080/api/messages/john@gmail.com?limit=3&cursor=encoded-cursor
```
Pull all messages from the database with given email address as well as the encoded cursor.

You can optionally specify the amount of emails you want as well as the cursor for pagination purposes


```
curl -X POST localhost:8080/api/send -d '{"magic_number":101}'
```
Send emails with the specified magic number using smtp. After that the messages are deleted from the database. 

The smtp credentials are hard-coded into the code to follow the specification.
So if this feature doesn't work, you have to change the email credentials that are located at service/util.go

## Tech stack
- Go
- Cassandra
- Docker

## How to run

```
docker-compose up
```
After Docker installs all the images you have to wait for a short while until Cassandra fully loads.
