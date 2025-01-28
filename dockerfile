# Используем легковесный образ Alpine Linux
FROM alpine:latest

# Устанавливаем curl
RUN apk --no-cache add curl

# Команда по умолчанию: выполняем запрос к API и выводим результат
CMD ["sh", "-c", "curl -s https://ifconfig.me"]
