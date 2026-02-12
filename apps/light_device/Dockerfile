FROM ruby:3.2-alpine

WORKDIR /app

RUN apk add --no-cache build-base tzdata

# Устанавливаем нужные гемы
RUN gem install sinatra puma json

COPY light_device.rb .
COPY config.ru .

EXPOSE 8085

# Запуск через puma и rackup
CMD ["puma", "-b", "tcp://0.0.0.0:8085", "config.ru"]
