package database

import (
	"babblegraph/util/testutils"
	"fmt"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	builder, err := NewBulkInsertQueryBuilder("things", "column1", "column2")
	if err != nil {
		t.Errorf("Did not expect error creating query builder, but got one: %s", err.Error())
		return
	}
	if err := builder.AddValues(5, "value"); err != nil {
		t.Errorf("Did not expect error adding values, but got one: %s", err.Error())
		return
	}
	if err := builder.AddValues(100, "another value"); err != nil {
		t.Errorf("Did not expect error adding values a second time, but got one: %s", err.Error())
		return
	}
	query, err := builder.buildQuery()
	if err != nil {
		t.Errorf("Did not expect error building query, but got one: %s", err.Error())
		return
	}
	expectedQuery := `INSERT INTO "things" ("column1", "column2") VALUES ($1, $2), ($3, $4)`
	if *query != expectedQuery {
		t.Errorf("Unexpected query. Expected %s, but got %s", expectedQuery, *query)
	}
}

func TestErrors(t *testing.T) {
	_, err := NewBulkInsertQueryBuilder("things")
	if err := testutils.CompareErrors(err, fmt.Errorf("Must supply at least one column")); err != nil {
		t.Errorf("Error on comparing creating bulk insert builder without columns: %s", err.Error())
	}
	_, err = NewBulkInsertQueryBuilder("", "column1")
	if err := testutils.CompareErrors(err, fmt.Errorf("Table name is invalid: must have at least one character")); err != nil {
		t.Errorf("Error on comparing creating bulk insert builder without valid table name: %s", err.Error())
	}
	validBuilder, err := NewBulkInsertQueryBuilder("things", "column1", "column2")
	if err := testutils.CompareErrors(err, nil); err != nil {
		t.Errorf("Error on comparing creating valid bulk insert builder: %s", err.Error())
		return
	}
	_, err = validBuilder.buildQuery()
	if err := testutils.CompareErrors(err, fmt.Errorf("Must add at least one value to insert query")); err != nil {
		t.Errorf("Error on comparing build query with no values: %s", err.Error())
	}
	err = validBuilder.AddValues(1)
	if err := testutils.CompareErrors(err, fmt.Errorf("Number of values must equal number of columns")); err != nil {
		t.Errorf("Error on comparing add values with too few values: %s", err.Error())
	}
	err = validBuilder.AddValues(1, 2, 3, 4)
	if err := testutils.CompareErrors(err, fmt.Errorf("Number of values must equal number of columns")); err != nil {
		t.Errorf("Error on comparing add values with too many values: %s", err.Error())
	}
	err = validBuilder.Execute(nil)
	if err := testutils.CompareErrors(err, fmt.Errorf("Must add at least one value to insert query")); err != nil {
		t.Errorf("Error on comparing execute with no values: %s", err.Error())
	}
}
