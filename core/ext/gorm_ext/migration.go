package gorm_ext

import (
	"cardamom/core/models"
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetColumnField(dbModel any, columnName string) *schema.Field {
	stmt, err := getTableStatement(models.DB, dbModel)
	if err != nil {
		panic(err)
	}

	return stmt.Schema.FieldsByDBName[columnName]
}

func GetColumnType(dbModel any, columnName string) gorm.ColumnType {
	columns, err := models.DB.Migrator().ColumnTypes(dbModel)
	if err != nil {
		panic(err)
	}

	return funk.Filter(columns, func(col gorm.ColumnType) bool {
		return col.Name() == columnName
	}).([]gorm.ColumnType)[0]
}

func GenerateMigration(dbModels ...any) ([]string, error) {

	var changes []string
	oldDBMigrator := models.DB.Session(&gorm.Session{}).Migrator()
	for _, model := range dbModels {

		if !oldDBMigrator.HasTable(model) {
			changes = append(changes, fmt.Sprintf(`mg.CreateTable(&md.%s{})`, reflect.TypeOf(model).Name()))
			continue
		}

		stmt, err := getTableStatement(models.DB, model)
		if err != nil {
			return nil, err
		}
		newSchema := stmt.Schema

		if newChanges, err := getColumnChanges(model, newSchema, oldDBMigrator); err == nil {
			changes = append(changes, newChanges...)
		} else {
			return nil, err
		}
		changes = append(changes, getConstraintChanges(model, newSchema, oldDBMigrator)...)
		changes = append(changes, getIndexChanges(model, newSchema, oldDBMigrator)...)

	}

	return changes, nil
}

func getTableStatement(db *gorm.DB, value any) (*gorm.Statement, error) {
	stmt := &gorm.Statement{DB: db}
	if db.Statement != nil {
		stmt.Table = db.Statement.Table
		stmt.TableExpr = db.Statement.TableExpr
	}

	if table, ok := value.(string); ok {
		stmt.Table = table
	} else if err := stmt.ParseWithSpecialTableName(value, stmt.Table); err != nil {
		return nil, err
	}

	return stmt, nil
}

func checkFieldUpdates(field *schema.Field, col gorm.ColumnType) bool {

	if length, ok := col.Length(); ok && length != int64(field.Size) {
		if length > 0 && field.Size > 0 {
			return true
		}
	}

	if pk, ok := col.PrimaryKey(); ok && pk != field.PrimaryKey {
		return true
	}

	if precision, _, ok := col.DecimalSize(); ok && precision != int64(field.Precision) {
		return true
	}

	if unique, ok := col.Unique(); ok && unique != field.Unique {
		return true
	}

	if nullable, ok := col.Nullable(); ok && !field.PrimaryKey && nullable == field.NotNull && nullable {
		return true
	}

	// check default value
	if dv, ok := col.DefaultValue(); ok && !field.PrimaryKey && dv != field.DefaultValue {
		return true
	}

	if comment, ok := col.Comment(); ok && comment != field.Comment {
		return true
	}

	return false
}

func getColumnChanges(model any, newSchema *schema.Schema, oldDBMigrator gorm.Migrator) ([]string, error) {

	db_cols, err := oldDBMigrator.ColumnTypes(model)
	if err != nil {
		return nil, err
	}

	var changes []string
	for _, col := range newSchema.DBNames {
		results := funk.Filter(db_cols,
			func(c gorm.ColumnType) bool { return c.Name() == col }).([]gorm.ColumnType)

		if len(results) == 0 {
			changes = append(changes, fmt.Sprintf(
				`mg.AddColumn(&md.%s{}, "%s")`, reflect.TypeOf(model).Name(), col))
		} else if checkFieldUpdates(newSchema.FieldsByDBName[col], results[0]) {
			modelName := reflect.TypeOf(model).Name()
			columnName := results[0].Name()
			changes = append(changes, fmt.Sprintf(
				`mg.MigrateColumn(&md.%s{}, ge.GetColumnField(md.%s{}, "%s"), ge.GetColumnType(md.%s{}, "%s"))`,
				modelName, modelName, columnName, modelName, columnName,
			))
		}
	}

	unusedCols := funk.Filter(db_cols, func(col gorm.ColumnType) bool {
		return !funk.Contains(newSchema.DBNames, col.Name())
	}).([]gorm.ColumnType)
	changes = append(changes, funk.Map(unusedCols, func(col gorm.ColumnType) string {
		return fmt.Sprintf(`mg.DropColumn(&md.%s{}, "%s")`, reflect.TypeOf(model).Name(), col.Name())
	}).([]string)...)

	return changes, nil
}

func getConstraintChanges(model any, newSchema *schema.Schema, oldDBMigrator gorm.Migrator) []string {

	var changes []string
	for _, rel := range newSchema.Relationships.Relations {
		if !models.DB.Config.DisableForeignKeyConstraintWhenMigrating {
			if constraint := rel.ParseConstraint(); constraint != nil &&
				constraint.Schema == newSchema && !oldDBMigrator.HasConstraint(model, constraint.Name) {
				changes = append(changes, fmt.Sprintf(`mg.CreateConstraint(&md.%s{}, "%s")`, reflect.TypeOf(model).Name(), constraint.Name))
			}
		}
	}

	for _, chk := range newSchema.ParseCheckConstraints() {
		if !oldDBMigrator.HasConstraint(model, chk.Name) {
			changes = append(changes, fmt.Sprintf(`mg.CreateConstraint(&md.%s{}, "%s")`, reflect.TypeOf(model).Name(), chk.Name))
		}
	}

	return changes
}

func getIndexChanges(model any, newSchema *schema.Schema, oldDBMigrator gorm.Migrator) []string {

	var changes []string
	// new indexes
	for _, idx := range newSchema.ParseIndexes() {
		if !oldDBMigrator.HasIndex(model, idx.Name) {
			changes = append(changes, fmt.Sprintf(`mg.CreateIndex(&md.%s{}, "%s")`, reflect.TypeOf(model).Name(), idx.Name))
		}
	}

	return changes
}
