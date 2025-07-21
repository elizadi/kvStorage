box.schema.space.create('kv_storage', {
    if_not_exists = true,
    format = {
        {name = 'key', type = 'string'},
        {name = 'value', type = 'string'}
    }
})

box.space.kv_storage:create_index('primary', {
    if_not_exists = true,
    parts = {'key'}})

