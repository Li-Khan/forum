FROM golang:1.17
ENV TZ="Asia/Almaty" 
RUN mkdir /app 
COPY . /app
WORKDIR /app 
RUN go build -o server ./cmd
EXPOSE 27969
CMD ["/app/server"]