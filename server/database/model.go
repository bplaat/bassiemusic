package database

import (
	"reflect"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type Model[T any] struct {
	TableName     string
	PrimaryKey    string
	Process       QueryBuilderProcess[T]
	Relationships map[string]QueryBuilderProcess[T]
	Columns       []QueryBuilderColumn
}

func (m Model[T]) Init() Model[T] {
	if m.PrimaryKey == "" {
		m.PrimaryKey = "id"
	}

	// Get columns info from struct tags
	var model T
	modelType := reflect.TypeOf(model)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("column")
		if tag != "" {
			parts := strings.Split(tag, ",")
			m.Columns = append(m.Columns, QueryBuilderColumn{
				Name:   field.Name,
				Column: parts[0],
				Type:   parts[1],
			})
		}
	}

	return m
}

func (m Model[T]) query() QueryBuilder[T] {
	return QueryBuilder[T]{
		Model:      &m,
		OrderByStr: "`created_at`",
	}
}

func (m Model[T]) Create(model *T) T {
	insertQuery := "INSERT INTO `" + m.TableName + "` ("
	index := 0
	for _, column := range m.Columns {
		insertQuery += "`" + column.Column + "`"
		if index != len(m.Columns)-1 {
			insertQuery += ", "
		}
		index++
	}

	modelType := reflect.ValueOf(model)
	modelID := uuid.NewV4()
	insertValues := []any{modelID}
	insertQuery += " VALUES (UUID_TO_BIN(?), "
	index = 0
	for _, column := range m.Columns {
		if column.Type == "uuid" {
			insertQuery += "UUID_TO_BIN(?)"
		} else {
			insertQuery += "?"
		}
		insertValues = append(insertValues, modelType.FieldByName(column.Name).Elem())
		if index != len(m.Columns)-1 {
			insertQuery += ", "
		}
		index++
	}
	insertQuery += ")"

	Exec(insertQuery, insertValues)

	return *m.Find(modelID.String())
}

func (m Model[T]) Update(model *T) {
	// TODO
}

func (m Model[T]) Join(join string) QueryBuilder[T] {
	return m.query().Join(join)
}

func (m Model[T]) With(relationships ...string) QueryBuilder[T] {
	return m.query().With(relationships...)
}

func (m Model[T]) Where(column string, value any) QueryBuilder[T] {
	return m.query().Where(column, value)
}

func (m Model[T]) WhereOr(column string, value any) QueryBuilder[T] {
	return m.query().WhereOr(column, value)
}

func (m Model[T]) WhereRaw(whereRaw string, value any) QueryBuilder[T] {
	return m.query().WhereRaw(whereRaw, value)
}

func (m Model[T]) WhereOrRaw(whereRaw string, value any) QueryBuilder[T] {
	return m.query().WhereOrRaw(whereRaw, value)
}

func (m Model[T]) WhereIn(pivotTableName string, pivotModelId string, pivotRelationshipId string, value string) QueryBuilder[T] {
	return m.query().WhereIn(pivotTableName, pivotModelId, pivotRelationshipId, value)
}

func (m Model[T]) OrderBy(column string) QueryBuilder[T] {
	return m.query().OrderBy(column)
}

func (m Model[T]) OrderByDesc(column string) QueryBuilder[T] {
	return m.query().OrderByDesc(column)
}

func (m Model[T]) OrderByRaw(orderByRaw string) QueryBuilder[T] {
	return m.query().OrderByRaw(orderByRaw)
}

func (m Model[T]) Limit(limit string) QueryBuilder[T] {
	return m.query().Limit(limit)
}

func (m Model[T]) Count() int64 {
	return m.query().Count()
}

func (m Model[T]) Get() []T {
	return m.query().Get()
}

func (m Model[T]) Paginate(page int, limit int) QueryBuilderPaginated[T] {
	return m.query().Paginate(page, limit)
}

func (m Model[T]) First() *T {
	return m.query().First()
}

func (m Model[T]) Find(primaryKey string) *T {
	return m.query().Where(m.PrimaryKey, primaryKey).First()
}
