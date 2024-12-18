![Static Badge](https://img.shields.io/badge/%D1%81%D1%82%D0%B0%D1%82%D1%83%D1%81-%D0%B3%D0%BE%D1%82%D0%BE%D0%B2%D0%BE-blue)
![Static Badge](https://img.shields.io/badge/GO-1.23+-blue)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/zagart47/rsssf)
![GitHub last commit (by committer)](https://img.shields.io/github/last-commit/zagart47/rsssf)
![GitHub forks](https://img.shields.io/github/forks/zagart47/rsssf)

# RSSSF
RSS сервис

## Содержание
- [Технологии](#технологии)
- [Использование](#использование)
- [Разработка](#разработка)
- [Contributing](#contributing)
- [FAQ](#faq)
- [To do](#to-do)
- [Команда проекта](#команда-проекта)

## Технологии
- [Golang](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)

## Использование
Внести настройки бд в ```config/config.json```.

Склонировать репозиторий
```powershell
git clone https://github.com/zagart47/rsssf.git
```
```powershell
cd cmd/rssf
```
```powershell
go run main.go
```


## Разработка

### Требования
Для установки и запуска проекта необходимы golang и прямые руки.

## Contributing
Если у вас есть предложения или идеи по дополнению проекта или вы нашли ошибку, то пишите мне в tg: @zagart47

## FAQ
### Зачем вы разработали этот проект?
Это тестовое задание.

## To do
- [x] Веб-приложение пользователя работоспособно. Отображаются последние новости из источников, указанных в конфигурации
- [x] Структура пакетов логична и отражает структуру приложения
- [x] Для всех пакетов (кроме, возможно, исполняемого) написаны тесты с достаточным покрытием
- [x] Модель данных соответствует требованиям и условиям задачи:
- [x] структура БД логична,
- [x] XML с потоком RSS декодируется верно.
- [x] В каталоге сервера существует корректный файл конфигурации config.json.	2
- [x] Присутствует файл со схемой БД (в случае использования реляционной СУБД).	3
- [x] Экспортируемые методы снабжены комментариями.
- [x] Для обхода лент RSS используются отдельные горутины.
- [x] Для обработки результатов обхода RSS и ошибок используются каналы.

## Команда проекта
- [Артур Загиров](https://t.me/zagart47) — Golang Developer

