# XTechnologies_test_case_golang
## Условие: 
*Требуется сделать сервис, который следит за изменением курсов BTC-USDT и фиатных валют к рублю, сохраняет изменения в БД и отдает клиентам по REST API.*
*Подробные условия описаны в файле **README_TERMS.md***
### Данное задание является тестовым, реализованы следующие endpoint'ы:
 - GET .. api/btcusdt - передает последнее значение btc/usdt
 - POST .. api/btcusdt - передает историю значений btc/usdt в JSON формате
 - GET .. api/currencies - передает последний курс рубля к валютам в JSON формате
 - POST .. api/currencies - передает историю курса рубля к валютам в JSON формате
 - GET .. api/latest - передает последний курс валют и btc в JSON формате
 - POST .. api/latest - передает историю курсов валют и btc в JSON формате
 - GET .. api/latest/{CHAR}* - передает курс необходимой валюты и btc
 - POST .. api/latest/{CHAR}* - передает историю курсов необходимой валюты и btc
 
*Например GET .. api/latest/USD - даст курс доллара и btc
### Описание реализации некоторых условий в ***README_TERMS.md***
- Реализация "мягкого выключения" БД реализована путем отслеживания системных сигналов о выключении/перезагрузки системы и закрытия БД перед завершением программы
- Реализация первичного запуска БД, а затем самого сервиса реализована в ***docker-compose.yaml***

### Описание модулей и их работы:
#### В данной программе реализована "плоская" структура с модулями, описанными ниже:
- ***apihandler.go*** - содержит функции, выводящие информацию в зависимости от метода http запроса
- ***calc.go*** - содержит функции, пересчитывающие валютные соотношения
- ***datahanler.go*** - содержит функции, обрабатывающие информацию от запрошенных http адресов
- ***errorlogger.go*** - макет обработчика ошибок
- ***main.go*** - основной модуль, содержит основные функции, объединяющие модули
- ***pgdbconnect.go*** - модуль для соединения с БД
- ***queries.go*** - содержит все функции, связанные с запросами к БД
- ***screen.go*** - макет для вывода информации в консоль (пока нерабочий, закомментирован)
- ***server.go*** - модуль для инициализации сервера
- ***source.go*** - содержит функции, получающие информацию от запрошенный http адресов
#### Общий алгоритм работы программы:
- Программа не грузит БД при многочисленных запросах по endpoint'ам. Все данные обновляются в переменных с циклом в 10сек.
- Если какой-либо запрашиваемый http адрес не дает код 200 - программа выдает ошибку, без предоставления информации при запросе на любой endpoint (для исключения вывода неактульной информации)

### Запуск сервиса:
- Для запуска сервиса с помощью Docker, используйте файлы в папке ***For Docker***
