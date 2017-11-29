Тестовое задание “чат"
======================

## Задача

Написать на go клиент и сервер для чата, которые будут общаться друг с другом через web socket. 

Клиент должен подключаться к серверу и отправлять текстовые сообщения. Формат сообщений:

    command_code[::[msg]]\r\n 
    
Возможные значения для command_code:

    auth — авторизация, msg — имя клиента,
    end — завершить сеанс, msg и “::" не должно быть,
    some_key — произвольный ключ (может быть пустым), msg — произвольное сообщение.

msg — произвольный текст, может быть пустым.


1. Сервер должен получать эти сообщения и выводить в формате
`[client_name]: some_key | msg`. 
`client_name` — имя клиента, которое было в `msg` при авторизации, `some_key` — ключи сообщений, которые приходили после авторизации `msg` — сообщения для соответствующих ключей.
2. Для отображение надо сделать страничку (веб интерфейс), на которой будет следующее:
список сообщений, которые обновляются в реальном времени, список клиентов, который обновляется в реальном времени, текущее количество клиентов,
текущее количество сообщений.
3. Все полученные сообщения должны храниться в вашей любимой базе данных и отображаться при подключению к веб интерфейсу.
4. Клиент должен иметь настройки для генерации сообщений (количество сообщений за сеанс, время сеанса и любые произвольные на ваш выбор). Настройки должны храниться в формате `JSON` или приниматься в командной строке при запуске клиента.
              
## Пример

Клиент отправляет:

    auth::Джеки Чан
    k1::Привет!
    k2::Меня зовут Джеки Чан
    k3::А тебя?
    end
    k4::Пока!

Результат, который отображается на сервере:

Сообщения: 2 | Клиенты: 1
-------------|-----------
[Джеки Чан]: k1 \| Привет! |   Джеки Чан
[Джеки Чан]: k2 \| Меня зовут Джеки Чан |



Сообщения: 3 | Клиенты: 0
-------------|-----------
[Джеки Чан]: k1 \| Привет! | Empty
[Джеки Чан]: k2 \| Меня зовут Джеки Чан |
[Джеки Чан]: k3 \| А тебя? |