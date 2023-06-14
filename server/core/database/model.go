package database

import (
	"reflect"
	"strings"

	"github.com/bplaat/bassiemusic/core/uuid"
)

type Map map[string]any

type ModelColumn struct {
	FieldName  string
	ColumnName string
	Type       string
}

type ModelProcessFunc[T any] func(model *T)
type ModelRelationshipFunc[T any] func(model *T, args []any)

type Model[T any] struct {
	TableName     string
	PrimaryKey    string
	Process       ModelProcessFunc[T]
	Relationships map[string]ModelRelationshipFunc[T]
	Columns       []*ModelColumn
	ColumnsLookup map[string]*ModelColumn
}

func (m *Model[T]) Init() *Model[T] {
	if m.PrimaryKey == "" {
		m.PrimaryKey = "id"
	}

	// Get columns info from struct tags
	m.ColumnsLookup = map[string]*ModelColumn{}
	var model T
	modelType := reflect.TypeOf(model)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("column")
		if tag != "" {
			parts := strings.Split(tag, ",")
			column := ModelColumn{
				FieldName:  field.Name,
				ColumnName: parts[0],
				Type:       parts[1],
			}
			m.Columns = append(m.Columns, &column)
			m.ColumnsLookup[column.ColumnName] = &column
		}
	}
	return m
}

func (m *Model[T]) Create(values Map) *T {
	// Create uuid id when it is not given
	if _, ok := values[m.PrimaryKey]; !ok {
		column := m.ColumnsLookup[m.PrimaryKey]
		if column.Type == "uuid" {
			values[m.PrimaryKey] = uuid.New().String()
		}
	}

	// Create insert SQL query
	insertQuery := "INSERT INTO `" + m.TableName + "` ("
	valuesQueryPart := ""
	queryValues := []any{}
	index := 0
	for columnName, value := range values {
		insertQuery += "`" + columnName + "`"
		if index != len(values)-1 {
			insertQuery += ", "
		}

		column := m.ColumnsLookup[columnName]
		if column.Type == "uuid" {
			valuesQueryPart += "UUID_TO_BIN(?)"
		} else {
			valuesQueryPart += "?"
		}
		queryValues = append(queryValues, value)
		if index != len(values)-1 {
			valuesQueryPart += ", "
		}
		index++
	}
	insertQuery += ") VALUES (" + valuesQueryPart + ")"

	// Run insert SQL query
	Exec(insertQuery, queryValues...)

	// Fetch newly created model
	return m.Find(values[m.PrimaryKey])
}

func (m *Model[T]) query() *QueryBuilder[T] {
	return &QueryBuilder[T]{model: m, withs: map[string][]any{}}
}

func (m *Model[T]) Select(columnNames ...string) *QueryBuilder[T] {
	return m.query().Select(columnNames...)
}

func (m *Model[T]) Join(join string) *QueryBuilder[T] {
	return m.query().Join(join)
}

func (m *Model[T]) With(relationships ...string) *QueryBuilder[T] {
	return m.query().With(relationships...)
}

func (m *Model[T]) WithArgs(relationship string, args ...any) *QueryBuilder[T] {
	return m.query().WithArgs(relationship, args...)
}

func (m *Model[T]) Where(columnName string, value any) *QueryBuilder[T] {
	return m.query().Where(columnName, value)
}
func (m *Model[T]) WhereOr(columnName string, value any) *QueryBuilder[T] {
	return m.query().WhereOr(columnName, value)
}

func (m *Model[T]) WhereRaw(whereRaw string, value any) *QueryBuilder[T] {
	return m.query().WhereRaw(whereRaw, value)
}
func (m *Model[T]) WhereOrRaw(whereRaw string, value any) *QueryBuilder[T] {
	return m.query().WhereOrRaw(whereRaw, value)
}

func (m *Model[T]) WhereNull(columnName string) *QueryBuilder[T] {
	return m.query().WhereNull(columnName)
}
func (m *Model[T]) WhereOrNull(columnName string) *QueryBuilder[T] {
	return m.query().WhereOrNull(columnName)
}

func (m *Model[T]) WhereNotNull(columnName string) *QueryBuilder[T] {
	return m.query().WhereNotNull(columnName)
}
func (m *Model[T]) WhereOrNotNull(columnName string) *QueryBuilder[T] {
	return m.query().WhereOrNotNull(columnName)
}

func (m *Model[T]) WhereIn(columnName string, list []any) *QueryBuilder[T] {
	return m.query().WhereIn(columnName, list)
}

func (m *Model[T]) WhereInQuery(columnName string, queryBuilder QueryBuilderSelectQuery) *QueryBuilder[T] {
	return m.query().WhereInQuery(columnName, queryBuilder)
}

func (m *Model[T]) OrderBy(columnName string) *QueryBuilder[T] {
	return m.query().OrderBy(columnName)
}

func (m *Model[T]) OrderByDesc(columnName string) *QueryBuilder[T] {
	return m.query().OrderByDesc(columnName)
}

func (m *Model[T]) OrderByRaw(orderByRaw string) *QueryBuilder[T] {
	return m.query().OrderByRaw(orderByRaw)
}

func (m *Model[T]) Offset(offset int64) *QueryBuilder[T] {
	return m.query().Offset(offset)
}

func (m *Model[T]) Limit(limit int64) *QueryBuilder[T] {
	return m.query().Limit(limit)
}

func (m *Model[T]) Count() int64 {
	return m.query().Count()
}

func (m *Model[T]) Get() []T {
	return m.query().Get()
}

func (m *Model[T]) Update(values Map) {
	m.query().Update(values)
}

func (m *Model[T]) Delete() {
	m.query().Delete()
}

func (m *Model[T]) Paginate(page int64, limit int64) QueryBuilderPaginated[T] {
	return m.query().Paginate(page, limit)
}

func (m *Model[T]) Chunk(limit int64, callback func(items []T)) {
	m.query().Chunk(limit, callback)
}

func (m *Model[T]) First() *T {
	return m.query().First()
}

func (m *Model[T]) Find(primaryKey any) *T {
	return m.query().Where(m.PrimaryKey, primaryKey).First()
}
