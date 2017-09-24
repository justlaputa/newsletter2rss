FROM centurylink/ca-certs
MAINTAINER laputa <justlaputa@gmail.com>

COPY templates/ /templates/
COPY news2rss /news2rss

ENTRYPOINT ["/news2rss"]
