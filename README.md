# newsletter2rss
convert email newsletter to rss

# how to test

create a new newsletter by post:

```
curl -X POST -d "title=microservice weekly" localhost:3000/feed
```

send testing mail:

```
nc localhost 2525 < parser/sample_email_tests/microservice-weekly.txt
```
