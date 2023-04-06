package database

import (
	"log"
	"math"
	"reflect"
	"strconv"
)

type QueryBuilderProcess[T any] func(model *T)

type QueryBuilder[T any] struct {
	Model       *Model[T]
	JoinStr     string
	Withs       []string
	WhereStr    string
	WhereValues []any
	OrderByStr  string
	OffsetInt   int
	LimitInt    int
}

type QueryBuilderColumn struct {
	Name   string
	Column string
	Type   string
}

type QueryBuilderPaginated[T any] struct {
	Data       []T `json:"data"`
	Pagination struct {
		Page  int   `json:"page"`
		Limit int   `json:"limit"`
		Total int64 `json:"total"`
	} `json:"pagination"`
}

func (qb *QueryBuilder[T]) Join(join string) *QueryBuilder[T] {
	qb.JoinStr = join
	return qb
}

func (qb *QueryBuilder[T]) With(relationships ...string) *QueryBuilder[T] {
	for _, relationship := range relationships {
		if _, ok := qb.Model.Relationships[relationship]; ok {
			qb.Withs = append(qb.Withs, relationship)
		} else {
			log.Fatalln("QueryBuilder: relationship '" + relationship + "' doesn't exists")
		}
	}
	return qb
}

func (qb *QueryBuilder[T]) FormatColumn(column string) string {
	if qb.JoinStr != "" {
		return "`" + qb.Model.TableName + "`.`" + column + "`"
	} else {
		return "`" + column + "`"
	}
}

func (qb *QueryBuilder[T]) where(column string, value any, operator string) *QueryBuilder[T] {
	for _, columnInfo := range qb.Model.Columns {
		if columnInfo.Column == column {
			if qb.WhereStr != "" {
				qb.WhereStr += " " + operator + " "
			}
			if columnInfo.Type == "uuid" {
				qb.WhereStr += qb.FormatColumn(column) + " = UUID_TO_BIN(?)"
			} else {
				qb.WhereStr += qb.FormatColumn(column) + " = ?"
			}
			qb.WhereValues = append(qb.WhereValues, value)
			return qb
		}
	}
	return qb
}
func (qb *QueryBuilder[T]) Where(column string, value any) *QueryBuilder[T] {
	return qb.where(column, value, "AND")
}
func (qb *QueryBuilder[T]) WhereOr(column string, value any) *QueryBuilder[T] {
	return qb.where(column, value, "OR")
}

func (qb *QueryBuilder[T]) whereRaw(whereRaw string, value any, operator string) *QueryBuilder[T] {
	if qb.WhereStr != "" {
		qb.WhereStr += " " + operator + " "
	}
	qb.WhereStr += whereRaw
	qb.WhereValues = append(qb.WhereValues, value)
	return qb
}
func (qb *QueryBuilder[T]) WhereRaw(whereRaw string, value any) *QueryBuilder[T] {
	return qb.whereRaw(whereRaw, value, "AND")
}
func (qb *QueryBuilder[T]) WhereOrRaw(whereRaw string, value any) *QueryBuilder[T] {
	return qb.whereRaw(whereRaw, value, "OR")
}

func (qb *QueryBuilder[T]) whereNull(column string, operator string) *QueryBuilder[T] {
	if qb.WhereStr != "" {
		qb.WhereStr += " " + operator + " "
	}
	qb.WhereStr += qb.FormatColumn(column) + " IS NULL"
	return qb
}
func (qb *QueryBuilder[T]) WhereNull(column string) *QueryBuilder[T] {
	return qb.whereNull(column, "AND")
}
func (qb *QueryBuilder[T]) WhereOrNull(column string) *QueryBuilder[T] {
	return qb.whereNull(column, "OR")
}

func (qb *QueryBuilder[T]) whereNotNull(column string, operator string) *QueryBuilder[T] {
	if qb.WhereStr != "" {
		qb.WhereStr += " " + operator + " "
	}
	qb.WhereStr += qb.FormatColumn(column) + " IS NOT NULL"
	return qb
}
func (qb *QueryBuilder[T]) WhereNotNull(column string) *QueryBuilder[T] {
	return qb.whereNotNull(column, "AND")
}
func (qb *QueryBuilder[T]) WhereOrNotNull(column string) *QueryBuilder[T] {
	return qb.whereNotNull(column, "OR")
}

func (qb *QueryBuilder[T]) WhereIn(pivotTableName string, pivotModelId string, pivotRelationshipId string, value any) *QueryBuilder[T] {
	qb.WhereStr += "`" + qb.Model.PrimaryKey + "` IN (SELECT `" + pivotModelId + "` FROM `" + pivotTableName + "` WHERE `" + pivotRelationshipId + "` = UUID_TO_BIN(?))"
	qb.WhereValues = append(qb.WhereValues, value)
	return qb
}

func (qb *QueryBuilder[T]) OrderBy(column string) *QueryBuilder[T] {
	qb.OrderByStr = qb.FormatColumn(column)
	return qb
}

func (qb *QueryBuilder[T]) OrderByDesc(column string) *QueryBuilder[T] {
	qb.OrderByStr = qb.FormatColumn(column) + " DESC"
	return qb
}

func (qb *QueryBuilder[T]) OrderByRaw(orderByRaw string) *QueryBuilder[T] {
	qb.OrderByStr = orderByRaw
	return qb
}

func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
	qb.OffsetInt = offset
	return qb
}

func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.LimitInt = limit
	return qb
}

func (qb *QueryBuilder[T]) Count() int64 {
	countQuery := "SELECT COUNT(" + qb.FormatColumn(qb.Model.PrimaryKey) + ") FROM `" + qb.Model.TableName + "`"
	if qb.JoinStr != "" {
		countQuery += " " + qb.JoinStr
	}
	if qb.WhereStr != "" {
		countQuery += " WHERE " + qb.WhereStr
	}
	if qb.LimitInt != 0 {
		if qb.OffsetInt != 0 {
			countQuery += " LIMIT " + strconv.Itoa(qb.OffsetInt) + ", " + strconv.Itoa(qb.LimitInt)
		} else {
			countQuery += " LIMIT " + strconv.Itoa(qb.LimitInt)
		}
	}

	query := Query(countQuery, qb.WhereValues...)
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
	for _, column := range qb.Model.Columns {
		if column.Type == "uuid" {
			selectQuery += "BIN_TO_UUID(" + qb.FormatColumn(column.Column) + ")"
		} else {
			selectQuery += qb.FormatColumn(column.Column)
		}
		if index != len(qb.Model.Columns)-1 {
			selectQuery += ", "
		}
		index++
	}
	selectQuery += " FROM `" + qb.Model.TableName + "`"
	if qb.JoinStr != "" {
		selectQuery += " " + qb.JoinStr
	}
	if qb.WhereStr != "" {
		selectQuery += " WHERE " + qb.WhereStr
	}
	if qb.OrderByStr != "" {
		selectQuery += " ORDER BY " + qb.OrderByStr
	}
	if qb.LimitInt != 0 {
		if qb.OffsetInt != 0 {
			selectQuery += " LIMIT " + strconv.Itoa(qb.OffsetInt) + ", " + strconv.Itoa(qb.LimitInt)
		} else {
			selectQuery += " LIMIT " + strconv.Itoa(qb.LimitInt)
		}
	}

	// Execute query and read models
	query := Query(selectQuery, qb.WhereValues...)
	models := []T{}
	for query.Next() {
		var model T
		modelValue := reflect.Indirect(reflect.ValueOf(&model))
		ptrs := []any{}
		for _, column := range qb.Model.Columns {
			ptrs = append(ptrs, modelValue.FieldByName(column.Name).Addr().Interface())
		}
		_ = query.Scan(ptrs...)
		models = append(models, model)
	}
	query.Close()

	// Process models and run relationships
	for i := 0; i < len(models); i++ {
		if qb.Model.Process != nil {
			qb.Model.Process(&models[i])
		}
		for _, with := range qb.Withs {
			qb.Model.Relationships[with](&models[i])
		}
	}
	return models
}

func (qb *QueryBuilder[T]) Update(values Map) {
	updateQuery := "UPDATE `" + qb.Model.TableName + "` SET "
	index := 0
	queryValues := []any{}
	for column, value := range values {
		for _, columnInfo := range qb.Model.Columns {
			if columnInfo.Column == column {
				if columnInfo.Type == "uuid" {
					updateQuery += qb.FormatColumn(column) + " = UUID_TO_BIN(?)"
				} else {
					updateQuery += qb.FormatColumn(column) + " = ?"
				}
				queryValues = append(queryValues, value)
				break
			}
		}
		if index != len(values)-1 {
			updateQuery += ", "
		}
		index++
	}
	if qb.WhereStr != "" {
		updateQuery += " WHERE " + qb.WhereStr
	}
	queryValues = append(queryValues, qb.WhereValues...)
	if qb.LimitInt != 0 {
		if qb.OffsetInt != 0 {
			updateQuery += " LIMIT " + strconv.Itoa(qb.OffsetInt) + ", " + strconv.Itoa(qb.LimitInt)
		} else {
			updateQuery += " LIMIT " + strconv.Itoa(qb.LimitInt)
		}
	}
	Exec(updateQuery, queryValues...)
}

func (qb *QueryBuilder[T]) Delete() {
	deleteQuery := "DELETE FROM `" + qb.Model.TableName + "` "
	if qb.JoinStr != "" {
		deleteQuery += " " + qb.JoinStr
	}
	if qb.WhereStr != "" {
		deleteQuery += " WHERE " + qb.WhereStr
	}
	if qb.LimitInt != 0 {
		if qb.OffsetInt != 0 {
			deleteQuery += " LIMIT " + strconv.Itoa(qb.OffsetInt) + ", " + strconv.Itoa(qb.LimitInt)
		} else {
			deleteQuery += " LIMIT " + strconv.Itoa(qb.LimitInt)
		}
	}
	Exec(deleteQuery, qb.WhereValues...)
}

func (qb *QueryBuilder[T]) Paginate(page int, limit int) QueryBuilderPaginated[T] {
	paginated := QueryBuilderPaginated[T]{}
	paginated.Pagination.Page = page
	paginated.Pagination.Limit = limit
	paginated.Pagination.Total = qb.Count()
	paginated.Data = qb.Offset((page - 1) * limit).Limit(limit).Get()
	return paginated
}

func (qb *QueryBuilder[T]) Chunk(limit int, callback func(items []T)) {
	total := qb.Count()
	for page := 1; page <= int(math.Ceil(float64(total)/float64(limit))); page++ {
		callback(qb.Offset((page - 1) * limit).Limit(limit).Get())
	}
}

func (qb *QueryBuilder[T]) First() *T {
	models := qb.Limit(1).Get()
	if len(models) == 0 {
		return nil
	}
	return &models[0]
}

func (qb *QueryBuilder[T]) Find(primaryKey any) *T {
	return qb.Where(qb.Model.PrimaryKey, primaryKey).First()
}
