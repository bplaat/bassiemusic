package database

import (
	"log"
	"reflect"
	"strconv"
)

type QueryBuilder[T any] struct {
	model          *Model[T]
	joinQueryPart  string
	withs          map[string][]any
	whereQueryPart string
	whereValues    []any
	orderBy        string
	offset         int64
	limit          int64
}

type QueryBuilderPaginated[T any] struct {
	Data       []T `json:"data"`
	Pagination struct {
		Page  int64 `json:"page"`
		Limit int64 `json:"limit"`
		Total int64 `json:"total"`
	} `json:"pagination"`
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
			log.Fatalln("QueryBuilder: relationship '" + relationship + "' doesn't exists")
		}
	}
	return qb
}

func (qb *QueryBuilder[T]) WithArgs(relationship string, args ...any) *QueryBuilder[T] {
	if _, ok := qb.model.Relationships[relationship]; ok {
		qb.withs[relationship] = append(qb.withs[relationship], args...)
	} else {
		log.Fatalln("QueryBuilder: relationship '" + relationship + "' doesn't exists")
	}
	return qb
}

func (qb *QueryBuilder[T]) FormatColumn(column string) string {
	if qb.joinQueryPart != "" {
		return "`" + qb.model.TableName + "`.`" + column + "`"
	} else {
		return "`" + column + "`"
	}
}

func (qb *QueryBuilder[T]) where(column string, value any, operator string) *QueryBuilder[T] {
	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " " + operator + " "
	}
	columnInfo := qb.model.ColumnsLookup[column]
	if columnInfo.Type == "uuid" {
		qb.whereQueryPart += qb.FormatColumn(column) + " = UUID_TO_BIN(?)"
	} else {
		qb.whereQueryPart += qb.FormatColumn(column) + " = ?"
	}
	qb.whereValues = append(qb.whereValues, value)
	return qb
}
func (qb *QueryBuilder[T]) Where(column string, value any) *QueryBuilder[T] {
	return qb.where(column, value, "AND")
}
func (qb *QueryBuilder[T]) WhereOr(column string, value any) *QueryBuilder[T] {
	return qb.where(column, value, "OR")
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

func (qb *QueryBuilder[T]) whereNull(column string, operator string) *QueryBuilder[T] {
	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " " + operator + " "
	}
	qb.whereQueryPart += qb.FormatColumn(column) + " IS NULL"
	return qb
}
func (qb *QueryBuilder[T]) WhereNull(column string) *QueryBuilder[T] {
	return qb.whereNull(column, "AND")
}
func (qb *QueryBuilder[T]) WhereOrNull(column string) *QueryBuilder[T] {
	return qb.whereNull(column, "OR")
}

func (qb *QueryBuilder[T]) whereNotNull(column string, operator string) *QueryBuilder[T] {
	if qb.whereQueryPart != "" {
		qb.whereQueryPart += " " + operator + " "
	}
	qb.whereQueryPart += qb.FormatColumn(column) + " IS NOT NULL"
	return qb
}
func (qb *QueryBuilder[T]) WhereNotNull(column string) *QueryBuilder[T] {
	return qb.whereNotNull(column, "AND")
}
func (qb *QueryBuilder[T]) WhereOrNotNull(column string) *QueryBuilder[T] {
	return qb.whereNotNull(column, "OR")
}

func (qb *QueryBuilder[T]) WhereIn(pivotTableName string, pivotModelId string, pivotRelationshipId string, value any) *QueryBuilder[T] {
	qb.whereQueryPart += "`" + qb.model.PrimaryKey + "` IN (SELECT `" + pivotModelId + "` FROM `" + pivotTableName + "` WHERE `" + pivotRelationshipId + "` = UUID_TO_BIN(?))"
	qb.whereValues = append(qb.whereValues, value)
	return qb
}

func (qb *QueryBuilder[T]) OrderBy(column string) *QueryBuilder[T] {
	qb.orderBy = qb.FormatColumn(column)
	return qb
}

func (qb *QueryBuilder[T]) OrderByDesc(column string) *QueryBuilder[T] {
	qb.orderBy = qb.FormatColumn(column) + " DESC"
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

func (qb *QueryBuilder[T]) Get() []T {
	// Build select query string
	selectQuery := "SELECT "
	index := 0
	for _, column := range qb.model.Columns {
		if column.Type == "uuid" {
			selectQuery += "BIN_TO_UUID(" + qb.FormatColumn(column.Column) + ")"
		} else {
			selectQuery += qb.FormatColumn(column.Column)
		}
		if index != len(qb.model.Columns)-1 {
			selectQuery += ", "
		}
		index++
	}
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

	// Execute query and read models
	query := Query(selectQuery, qb.whereValues...)
	models := []T{}
	for query.Next() {
		var model T
		modelValue := reflect.ValueOf(&model).Elem()
		ptrs := []any{}
		for _, column := range qb.model.Columns {
			ptrs = append(ptrs, modelValue.FieldByName(column.Name).Addr().Interface())
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
	if len(values) == 0 {
		return
	}
	updateQuery := "UPDATE `" + qb.model.TableName + "` SET "
	index := 0
	queryValues := []any{}
	for column, value := range values {
		columnInfo := qb.model.ColumnsLookup[column]
		if columnInfo.Type == "uuid" {
			updateQuery += qb.FormatColumn(column) + " = UUID_TO_BIN(?)"
		} else {
			updateQuery += qb.FormatColumn(column) + " = ?"
		}
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
	Exec(updateQuery, queryValues...)
}

func (qb *QueryBuilder[T]) Delete() {
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
