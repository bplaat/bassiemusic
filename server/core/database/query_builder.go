package database

import (
	"log"
	"reflect"
	"strconv"
)

type QueryBuilder[T any] struct {
	model             *Model[T]
	selectColumnNames []string
	joinQueryPart     string
	withs             map[string][]any
	whereQueryPart    string
	whereValues       []any
	orderBy           string
	offset            int64
	limit             int64
}

type QueryBuilderSelectQuery interface {
	SelectQuery(whereInQuery bool) (string, []any)
}

type QueryBuilderPaginated[T any] struct {
	Data       []T `json:"data"`
	Pagination struct {
		Page  int64 `json:"page"`
		Limit int64 `json:"limit"`
		Total int64 `json:"total"`
	} `json:"pagination"`
}

func (qb *QueryBuilder[T]) Select(columnNames ...string) *QueryBuilder[T] {
	// Check if columns exist
	for _, columnName := range columnNames {
		if _, ok := qb.model.ColumnsLookup[columnName]; !ok {
			log.Fatalf("QueryBuilder: column '%s' doesn't exist on model '%s'\n", columnName, qb.model.TableName)
		}
	}

	qb.selectColumnNames = append(qb.selectColumnNames, columnNames...)
	return qb
}

func (qb *QueryBuilder[T]) Join(join string) *QueryBuilder[T] {
	qb.joinQueryPart = join
	return qb
}

func (qb *QueryBuilder[T]) With(relationships ...string) *QueryBuilder[T] {
	for _, relationship := range relationships {
		if _, ok := qb.model.Relationships[relationship]; ok {
			qb.withs[relationship] = []any{}
		} else {
			log.Fatalf("QueryBuilder: relationship '%s' doesn't exists on model '%s'\n", relationship, qb.model.TableName)
		}
	}
	return qb
}

func (qb *QueryBuilder[T]) WithArgs(relationship string, args ...any) *QueryBuilder[T] {
	if _, ok := qb.model.Relationships[relationship]; ok {
		qb.withs[relationship] = append(qb.withs[relationship], args...)
	} else {
		log.Fatalf("QueryBuilder: relationship '%s' doesn't exists on model '%s'\n", relationship, qb.model.TableName)
	}
	return qb
}

func (qb *QueryBuilder[T]) FormatColumn(columnName string) string {
	// Check if columns exist
	if _, ok := qb.model.ColumnsLookup[columnName]; !ok {
		log.Fatalf("QueryBuilder: column '%s' doesn't exist on model '%s'\n", columnName, qb.model.TableName)
	}

	if qb.joinQueryPart != "" {
		return "`" + qb.model.TableName + "`.`" + columnName + "`"
	} else {
		return "`" + columnName + "`"
	}
}

func (qb *QueryBuilder[T]) where(columnName string, value any, operator string) *QueryBuilder[T] {
	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " " + operator + " "
	}
	qb.whereQueryPart += qb.FormatColumn(columnName) + " = ?"
	qb.whereValues = append(qb.whereValues, value)
	return qb
}
func (qb *QueryBuilder[T]) Where(columnName string, value any) *QueryBuilder[T] {
	return qb.where(columnName, value, "AND")
}
func (qb *QueryBuilder[T]) WhereOr(columnName string, value any) *QueryBuilder[T] {
	return qb.where(columnName, value, "OR")
}

func (qb *QueryBuilder[T]) whereRaw(whereRaw string, value any, operator string) *QueryBuilder[T] {
	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " " + operator + " "
	}
	qb.whereQueryPart += whereRaw
	qb.whereValues = append(qb.whereValues, value)
	return qb
}
func (qb *QueryBuilder[T]) WhereRaw(whereRaw string, value any) *QueryBuilder[T] {
	return qb.whereRaw(whereRaw, value, "AND")
}
func (qb *QueryBuilder[T]) WhereOrRaw(whereRaw string, value any) *QueryBuilder[T] {
	return qb.whereRaw(whereRaw, value, "OR")
}

func (qb *QueryBuilder[T]) whereNull(columnName string, operator string) *QueryBuilder[T] {
	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " " + operator + " "
	}
	qb.whereQueryPart += qb.FormatColumn(columnName) + " IS NULL"
	return qb
}
func (qb *QueryBuilder[T]) WhereNull(columnName string) *QueryBuilder[T] {
	return qb.whereNull(columnName, "AND")
}
func (qb *QueryBuilder[T]) WhereOrNull(columnName string) *QueryBuilder[T] {
	return qb.whereNull(columnName, "OR")
}

func (qb *QueryBuilder[T]) whereNotNull(columnName string, operator string) *QueryBuilder[T] {
	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " " + operator + " "
	}
	qb.whereQueryPart += qb.FormatColumn(columnName) + " IS NOT NULL"
	return qb
}
func (qb *QueryBuilder[T]) WhereNotNull(columnName string) *QueryBuilder[T] {
	return qb.whereNotNull(columnName, "AND")
}
func (qb *QueryBuilder[T]) WhereOrNotNull(columnName string) *QueryBuilder[T] {
	return qb.whereNotNull(columnName, "OR")
}

func (qb *QueryBuilder[T]) WhereIn(columnName string, list []any) *QueryBuilder[T] {
	if len(list) == 1 {
		return qb.Where(columnName, list[0])
	}

	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " AND "
	}
	inQueryPart := ""
	for i := 0; i < len(list); i++ {
		inQueryPart += "?"
		if i != len(list)-1 {
			inQueryPart += ", "
		}
	}
	qb.whereQueryPart += qb.FormatColumn(columnName) + " IN (" + inQueryPart + ")"
	qb.whereValues = append(qb.whereValues, list...)
	return qb
}

func (qb *QueryBuilder[T]) WhereInQuery(columnName string, queryBuilder QueryBuilderSelectQuery) *QueryBuilder[T] {
	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " AND "
	}
	query, whereValues := queryBuilder.SelectQuery(true)
	qb.whereQueryPart += qb.FormatColumn(columnName) + " IN (" + query + ")"
	qb.whereValues = append(qb.whereValues, whereValues...)
	return qb
}

func (qb *QueryBuilder[T]) OrderBy(columnName string) *QueryBuilder[T] {
	qb.orderBy = qb.FormatColumn(columnName)
	return qb
}

func (qb *QueryBuilder[T]) OrderByDesc(columnName string) *QueryBuilder[T] {
	qb.orderBy = qb.FormatColumn(columnName) + " DESC"
	return qb
}

func (qb *QueryBuilder[T]) OrderByRaw(orderByRaw string) *QueryBuilder[T] {
	qb.orderBy = orderByRaw
	return qb
}

func (qb *QueryBuilder[T]) Offset(offset int64) *QueryBuilder[T] {
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder[T]) Limit(limit int64) *QueryBuilder[T] {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder[T]) Count() int64 {
	countQuery := "SELECT COUNT(" + qb.FormatColumn(qb.model.PrimaryKey) + ") FROM `" + qb.model.TableName + "`"
	if qb.joinQueryPart != "" {
		countQuery += " " + qb.joinQueryPart
	}
	if qb.whereQueryPart != "" {
		countQuery += " WHERE " + qb.whereQueryPart
	}
	if qb.limit != 0 {
		if qb.offset != 0 {
			countQuery += " LIMIT " + strconv.FormatInt(qb.offset, 10) + ", " + strconv.FormatInt(qb.limit, 10)
		} else {
			countQuery += " LIMIT " + strconv.FormatInt(qb.limit, 10)
		}
	}

	query := Query(countQuery, qb.whereValues...)
	defer query.Close()
	query.Next()
	var count int64
	_ = query.Scan(&count)
	return count
}

func (qb *QueryBuilder[T]) selectColumns() []*ModelColumn {
	var selectColumns []*ModelColumn
	if len(qb.selectColumnNames) > 0 {
		for _, columnName := range qb.selectColumnNames {
			selectColumns = append(selectColumns, qb.model.ColumnsLookup[columnName])
		}
	} else {
		selectColumns = qb.model.Columns
	}
	return selectColumns
}

func (qb *QueryBuilder[T]) SelectQuery(whereInQuery bool) (string, []any) {
	selectQuery := "SELECT "

	// Add selected columns to the query
	selectColumns := qb.selectColumns()
	index := 0
	for _, column := range selectColumns {
		selectQuery += qb.FormatColumn(column.ColumnName)
		if index != len(selectColumns)-1 {
			selectQuery += ", "
		}
		index++
	}

	// Add rest of the stuff to the query
	selectQuery += " FROM `" + qb.model.TableName + "`"
	if qb.joinQueryPart != "" {
		selectQuery += " " + qb.joinQueryPart
	}
	if qb.whereQueryPart != "" {
		selectQuery += " WHERE " + qb.whereQueryPart
	}
	if qb.orderBy != "" {
		selectQuery += " ORDER BY " + qb.orderBy
	}
	if qb.limit != 0 {
		if qb.offset != 0 {
			selectQuery += " LIMIT " + strconv.FormatInt(qb.offset, 10) + ", " + strconv.FormatInt(qb.limit, 10)
		} else {
			selectQuery += " LIMIT " + strconv.FormatInt(qb.limit, 10)
		}
	}
	return selectQuery, qb.whereValues
}

func (qb *QueryBuilder[T]) Get() []T {
	// Build select query string
	selectQuery, whereValues := qb.SelectQuery(false)

	// Execute query and read models
	selectColumns := qb.selectColumns()
	query := Query(selectQuery, whereValues...)
	models := []T{}
	for query.Next() {
		var model T
		modelValue := reflect.ValueOf(&model).Elem()
		ptrs := []any{}
		for _, column := range selectColumns {
			ptrs = append(ptrs, modelValue.FieldByName(column.FieldName).Addr().Interface())
		}
		_ = query.Scan(ptrs...)
		models = append(models, model)
	}
	query.Close()

	// Process models and run relationships
	for i := 0; i < len(models); i++ {
		if qb.model.Process != nil {
			qb.model.Process(&models[i])
		}
		for relationship, args := range qb.withs {
			qb.model.Relationships[relationship](&models[i], args)
		}
	}
	return models
}

func (qb *QueryBuilder[T]) Update(values Map) {
	// Check if we want to update something
	if len(values) == 0 {
		return
	}

	// Check if columns exists
	for columnName := range values {
		if _, ok := qb.model.ColumnsLookup[columnName]; !ok {
			log.Fatalf("QueryBuilder: column '%s' doesn't exist on model '%s'\n", columnName, qb.model.TableName)
		}
	}

	// Create update SQL query
	updateQuery := "UPDATE `" + qb.model.TableName + "` SET "
	index := 0
	queryValues := []any{}
	for columnName, value := range values {
		updateQuery += qb.FormatColumn(columnName) + " = ?"
		queryValues = append(queryValues, value)
		if index != len(values)-1 {
			updateQuery += ", "
		}
		index++
	}
	if qb.whereQueryPart != "" {
		updateQuery += " WHERE " + qb.whereQueryPart
	}
	queryValues = append(queryValues, qb.whereValues...)
	if qb.limit != 0 {
		if qb.offset != 0 {
			updateQuery += " LIMIT " + strconv.FormatInt(qb.offset, 10) + ", " + strconv.FormatInt(qb.limit, 10)
		} else {
			updateQuery += " LIMIT " + strconv.FormatInt(qb.limit, 10)
		}
	}

	// Execute update query
	Exec(updateQuery, queryValues...)
}

func (qb *QueryBuilder[T]) Delete() {
	// Create delete SQL query
	deleteQuery := "DELETE FROM `" + qb.model.TableName + "` "
	if qb.joinQueryPart != "" {
		deleteQuery += " " + qb.joinQueryPart
	}
	if qb.whereQueryPart != "" {
		deleteQuery += " WHERE " + qb.whereQueryPart
	}
	if qb.limit != 0 {
		if qb.offset != 0 {
			deleteQuery += " LIMIT " + strconv.FormatInt(qb.offset, 10) + ", " + strconv.FormatInt(qb.limit, 10)
		} else {
			deleteQuery += " LIMIT " + strconv.FormatInt(qb.limit, 10)
		}
	}

	// Execute delete query
	Exec(deleteQuery, qb.whereValues...)
}

func (qb *QueryBuilder[T]) Paginate(page int64, limit int64) QueryBuilderPaginated[T] {
	paginated := QueryBuilderPaginated[T]{}
	paginated.Pagination.Page = page
	paginated.Pagination.Limit = limit
	paginated.Pagination.Total = qb.Count()
	paginated.Data = qb.Offset((page - 1) * limit).Limit(limit).Get()
	return paginated
}

func (qb *QueryBuilder[T]) Chunk(limit int64, callback func(items []T)) {
	total := qb.Count()
	for offset := int64(0); offset < total; offset += limit {
		callback(qb.Offset(offset).Limit(limit).Get())
	}
}

func (qb *QueryBuilder[T]) First() *T {
	models := qb.Limit(1).Get()
	if len(models) > 0 {
		return &models[0]
	}
	return nil
}

func (qb *QueryBuilder[T]) Find(primaryKey any) *T {
	return qb.Where(qb.model.PrimaryKey, primaryKey).First()
}
