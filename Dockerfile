FROM ubuntu:latest
LABEL authors="tomislavlovrencic"

ENTRYPOINT ["top", "-b"]