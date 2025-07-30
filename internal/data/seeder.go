package data

import (
	"fmt"
	"go-42/internal/data/entity"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// // seed all data
// func SeedAll(db *gorm.DB) error {
// 	return db.Transaction(func(tx *gorm.DB) error {
// 		seeds := dataSeeds()
// 		for i := range seeds {
// 			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(seeds[i]).Error
// 			if nil != err {
// 				name := reflect.TypeOf(seeds[i]).String()
// 				errorMessage := err.Error()
// 				return fmt.Errorf("%s seeder fail with %s", name, errorMessage)
// 			}
// 		}
// 		return nil
// 	})
// }

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		seeds := dataSeeds()

		for _, item := range seeds {
			// ambil nama struct
			tableName := structToTableName(item)

			// cek apakah table ada data
			var count int64
			err := tx.Table(tableName).Count(&count).Error
			if err != nil {
				return fmt.Errorf("failed to count table %s: %w", tableName, err)
			}

			if count > 0 {
				// hapus semua data
				fmt.Printf("Table %s has %d rows, deleting...\n", tableName, count)

				if err := tx.Exec(fmt.Sprintf("DELETE FROM %s", tableName)).Error; err != nil {
					return fmt.Errorf("failed to delete table %s: %w", tableName, err)
				}

				// reset sequence (optional, kalau pakai auto-increment ID)
				if err := resetSequence(tx, tableName); err != nil {
					fmt.Printf("Warning: failed to reset sequence for table %s: %s\n", tableName, err.Error())
				}
			}

			// insert fresh data
			if err := tx.Create(item).Error; err != nil {
				return fmt.Errorf("failed to insert seed into %s: %w", tableName, err)
			}

			fmt.Printf("Seeding %s finished.\n", tableName)
		}

		return nil
	})
}

// DataSeeds data
func dataSeeds() []interface{} {
	return []interface{}{
		entity.SeedUsers(),
		// entity.SeedWallets(),
	}
}

func structToTableName(data interface{}) string {
	typ := reflect.TypeOf(data)

	// kalau slice, ambil type elemennya
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	name := typ.Name()

	// convert CamelCase → snake_case
	snake := camelToSnake(name)

	// plural → tambahkan "s"
	return snake + "s"
}

func camelToSnake(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func resetSequence(tx *gorm.DB, tableName string) error {
	seqName := tableName + "_id_seq"
	sql := fmt.Sprintf(`ALTER SEQUENCE %s RESTART WITH 1`, seqName)
	return tx.Exec(sql).Error
}
