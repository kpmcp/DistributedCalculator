# Distributed Calculator

## [Запуск](#Запуск)
## [Структура](#Структура)
## [Как работает](#Взаимодействие)

### Запуск
В терминале ввести `go run main.go`.\
Далее перейти по адресу `http://localhost:8080`\
На сайте есть 2 страницы - с вводом выражения и со списком всех считающихся выражений, после ввода выражения и нажатия кнопки посчитать, операция будет находится во вкладке `Выражения` (чтобы увидеть результат обновляйте страницу)\
Чтобы изменить время выполнения операций или количество воркеров, в файле `pkg/env/env.go` нужно изменить переменные Plus, Minus, Mul, Div, Workers

### Cтруктура
В папке `frontend` лежат html страницы.\
Папка `pkg` - база данных, логер, парсер выражений\
Папка `internal` - агент и воркер\
Папка `http` - оркестратор

### Взаимодействие
Полученное с сайта выражение записывается в базу данных, далее его обрабатывает парсер.\
Затем делается запрос к агенту, выражение поступает на обработку (если все воркеры заняты, то выражение добавляется в список Waiting)\
Вычисленное выражение записывается в базу данных
