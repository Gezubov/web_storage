FROM ubuntu:latest
LABEL authors="georgijzubov"

ENTRYPOINT ["top", "-b"]