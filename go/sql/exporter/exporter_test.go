package exporter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTables struct {
	TableOne []tableOneRow
	TableTwo []tableTwoRow
}

type tableOneRow struct {
	ColumnOne string `sql:"column_one STRING PRIMARY KEY"`
	ColumnTwo int    `sql:"column_two INT8 NOT NULL"`
}

type tableTwoRow struct {
	CompositeKeyOne []byte `sql:"comp_one BYTES"`
	CompositeKeyTwo []byte `sql:"comp_two BYTES"`
	// This is a test comment
	primaryKey struct{} `sql:"PRIMARY KEY (comp_one, comp_two)"`
	// We have this index for good reasons. Spaces are intentional to make sure they are trimmed
	indexForSearch struct{} `sql:"   INDEX comp_two_desc_idx (comp_two DESC)   "`
	// We have this index for even better reasons.
	otherIndex struct{} `sql:"INDEX comp_two_asc_idx (comp_two ASC)"`
}

func TestGenerateSQL_WellFormedInput_CorrectOutput(t *testing.T) {

	gen := GenerateSQL(testTables{}, "test_package_one", SchemaOnly, CockroachDB, nil)
	// We cannot have backticks in the multistring literal, so we substitute $$ for them.
	expectedOutput := strings.ReplaceAll(`package test_package_one

// Generated by //go/sql/exporter/
// DO NOT EDIT

const Schema = $$CREATE TABLE IF NOT EXISTS TableOne (
  column_one STRING PRIMARY KEY,
  column_two INT8 NOT NULL
);
CREATE TABLE IF NOT EXISTS TableTwo (
  comp_one BYTES,
  comp_two BYTES,
  PRIMARY KEY (comp_one, comp_two),
  INDEX comp_two_desc_idx (comp_two DESC),
  INDEX comp_two_asc_idx (comp_two ASC)
);
$$
`, "$$", "`")

	assert.Equal(t, expectedOutput, gen)
}

func TestGenerateSQL_WellFormedInput_CorrectOutputIncludingColumnNames(t *testing.T) {

	gen := GenerateSQL(testTables{}, "test_package_one", SchemaAndColumnNames, CockroachDB, nil)
	// We cannot have backticks in the multistring literal, so we substitute $$ for them.
	expectedOutput := strings.ReplaceAll(`package test_package_one

// Generated by //go/sql/exporter/
// DO NOT EDIT

const Schema = $$CREATE TABLE IF NOT EXISTS TableOne (
  column_one STRING PRIMARY KEY,
  column_two INT8 NOT NULL
);
CREATE TABLE IF NOT EXISTS TableTwo (
  comp_one BYTES,
  comp_two BYTES,
  PRIMARY KEY (comp_one, comp_two),
  INDEX comp_two_desc_idx (comp_two DESC),
  INDEX comp_two_asc_idx (comp_two ASC)
);
$$

var TableOne = []string{
	"column_one",
	"column_two",
}

var TableTwo = []string{
	"comp_one",
	"comp_two",
}
`, "$$", "`")

	assert.Equal(t, expectedOutput, gen)
}

type malformedTable struct {
	TableOne []tableOneRow
	TableTwo tableTwoRow
}

func TestGenerateSQL_NonSliceField_Panics(t *testing.T) {

	assert.Panics(t, func() {
		GenerateSQL(malformedTable{}, "test_package_one", SchemaOnly, CockroachDB, nil)
	})
}

type missingSQLStructs struct {
	TableOne []tableOneRow
	TableTwo []malformedTable
}

func TestGenerateSQL_MissingSQLStructTag_Panics(t *testing.T) {

	assert.Panics(t, func() {
		GenerateSQL(missingSQLStructs{}, "test_package_one", SchemaOnly, CockroachDB, nil)
	})
}
