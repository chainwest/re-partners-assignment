# Миграция на sqlx

## Что изменилось

Проект был переписан с использования `database/sql` на `github.com/jmoiron/sqlx` для упрощения работы с базой данных PostgreSQL.

## Основные преимущества sqlx

### 1. Автоматическое сканирование в структуры

**До (database/sql):**
```go
var model PackSetModel
err := r.db.QueryRowContext(ctx, query, id).Scan(
    &model.ID,
    &model.Name,
    &model.Sizes,
    &model.CreatedAt,
    &model.UpdatedAt,
)
```

**После (sqlx):**
```go
var model PackSetModel
err := r.db.GetContext(ctx, &model, query, id)
```

### 2. Именованные параметры

**До (database/sql):**
```go
query := `
    UPDATE pack_sets
    SET name = $1, sizes = $2, updated_at = $3
    WHERE id = $4
`
result, err := r.db.ExecContext(ctx, query,
    model.Name,
    model.Sizes,
    model.UpdatedAt,
    model.ID,
)
```

**После (sqlx):**
```go
query := `
    UPDATE pack_sets
    SET name = :name, sizes = :sizes, updated_at = :updated_at
    WHERE id = :id
`
result, err := r.db.NamedExecContext(ctx, query, model)
```

### 3. Упрощённая работа со списками

**До (database/sql):**
```go
rows, err := r.db.QueryContext(ctx, query, limit, offset)
if err != nil {
    return nil, err
}
defer rows.Close()

var packSets []*domain.PackSizeSet
for rows.Next() {
    var model PackSetModel
    err := rows.Scan(
        &model.ID,
        &model.Name,
        &model.Sizes,
        &model.CreatedAt,
        &model.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    packSets = append(packSets, model.ToPackSizeSet())
}

if err := rows.Err(); err != nil {
    return nil, err
}
```

**После (sqlx):**
```go
var models []PackSetModel
err := r.db.SelectContext(ctx, &models, query, limit, offset)
if err != nil {
    return nil, err
}

packSets := make([]*domain.PackSizeSet, 0, len(models))
for i := range models {
    packSets = append(packSets, models[i].ToPackSizeSet())
}
```

### 4. Использование тегов структур

Все модели используют теги `db` для маппинга полей:

```go
type PackSetModel struct {
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    Sizes     IntArray  `db:"sizes"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```

## Изменённые файлы

1. **go.mod** - добавлена зависимость `github.com/jmoiron/sqlx v1.4.0`
2. **internal/infra/postgres/db.go** - изменён тип возвращаемого значения с `*sql.DB` на `*sqlx.DB`
3. **internal/infra/postgres/repository.go** - все методы переписаны с использованием sqlx API
4. **cmd/api/main.go** - обновлён импорт и тип переменной `db`

## Обратная совместимость

- Все публичные API остались без изменений
- Тесты проходят успешно
- Функциональность не изменилась

## Производительность

sqlx не добавляет значительных накладных расходов по сравнению с `database/sql`, так как использует те же самые низкоуровневые механизмы. Основное преимущество - удобство разработки и меньше шаблонного кода.

## Дополнительные возможности sqlx

- `sqlx.In()` - для работы с IN-запросами
- `sqlx.Named()` - для преобразования именованных запросов
- `DB.Rebind()` - для автоматической замены плейсхолдеров под разные БД
- Поддержка транзакций с теми же удобными методами
- Расширенная поддержка NULL-значений

## Ссылки

- [Документация sqlx](https://jmoiron.github.io/sqlx/)
- [GitHub репозиторий](https://github.com/jmoiron/sqlx)

