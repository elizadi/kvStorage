## kv-хранилище с использованием tarantool
kv-хранилище с доступом по http, реализованное на tarantool

### Доступные методы:
- `/key_value_storage/:key` Метод Get, возвращает значение, хранящееся по указанному ключу
- `/key_value_storage/new_value` Метод Post, сохраняет пару ключ-значение. Принимает на вход json: 
``` json 
{
    "key":"string",
    "value":"string"
}
```
- `/key_value_storage/value` Метод Post, редактирует уже существующую пару ключ-значение. Принимает на вход json:
``` json 
{
    "key":"string",
    "value":"string"
}
```
- `/key_value_storage/:key` Метод Delete, удаляет значение по указанному ключу

Публичный адрес для запросов: `https://immutably-learning-squirrel.cloudpub.ru`
