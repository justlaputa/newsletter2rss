FROM centurylink/ca-certs
MAINTAINER laputa <justlaputa@gmail.com>

COPY templates/ /templates/
COPY news2rss /news2rss
COPY web-app/build/ /public

ENTRYPOINT ["/news2rss"]
