package database

import (
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type bulkInsertQueryBuilder struct {
	tableName   string
	columnNames []string
	onConflict  *string
	values      []interface{}
}

func NewBulkInsertQueryBuilder(tableName string, columnNames ...string) (*bulkInsertQueryBuilder, error) {
	switch {
	case len(tableName) == 0:
		return nil, fmt.Errorf("Table name is invalid: must have at least one character")
	case len(columnNames) == 0:
		return nil, fmt.Errorf("Must supply at least one column")
	default:
		return &bulkInsertQueryBuilder{
			tableName:   tableName,
			columnNames: columnNames,
		}, nil
	}
}

// This is a hack and definitely needs to be fixed
func (b *bulkInsertQueryBuilder) AddConflictResolution(conflictString string) {
	b.onConflict = ptr.String(conflictString)
}

func (b *bulkInsertQueryBuilder) AddValues(values ...interface{}) error {
	if len(values) != len(b.columnNames) {
		return fmt.Errorf("Number of values must equal number of columns")
	}
	b.values = append(b.values, values...)
	return nil
}

func (b *bulkInsertQueryBuilder) buildQuery() (*string, error) {
	if len(b.values) == 0 {
		return nil, fmt.Errorf("Must add at least one value to insert query")
	}
	var quotedColumnNames []string
	for _, columnName := range b.columnNames {
		quotedColumnNames = append(quotedColumnNames, fmt.Sprintf(`"%s"`, columnName))
	}
	queryStart := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES`, b.tableName, strings.Join(quotedColumnNames, ", "))
	numberOfValueGroups := len(b.values) / len(b.columnNames)
	var valueGroups []string
	for idx := 0; idx < numberOfValueGroups; idx++ {
		var valueStringParts []string
		for columnIdx, _ := range b.columnNames {
			position := (columnIdx + 1) + (idx * len(b.columnNames))
			valueStringParts = append(valueStringParts, fmt.Sprintf("$%d", position))
		}
		valueGroups = append(valueGroups, fmt.Sprintf("(%s)", strings.Join(valueStringParts, ", ")))
	}
	var conflictString *string
	if b.onConflict != nil {
		conflictString = ptr.String(fmt.Sprintf(" ON CONFLICT %s", *b.onConflict))
	}
	query := fmt.Sprintf("%s %s%s", queryStart, strings.Join(valueGroups, ", "), deref.String(conflictString, ""))
	return ptr.String(query), nil
}

func (b *bulkInsertQueryBuilder) Execute(tx *sqlx.Tx) error {
	query, err := b.buildQuery()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(*query, b.values...); err != nil {
		return err
	}
	return nil
}
