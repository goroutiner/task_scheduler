FROM ubuntu:latest
WORKDIR /app
COPY . ./
EXPOSE 7540
RUN apt-get update
CMD [ "./cmd/scheduler" ]