package services

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type BaseService struct {
	DB *gorm.DB
}

// Create a new instance with a transaction DB
func (s *BaseService) WithTx(tx *gorm.DB) *BaseService {
	return &BaseService{DB: tx}
}

// Generic CRUD operations
func (s *BaseService) Create(model any) error {
	return s.DB.Create(model).Error
}

func (s *BaseService) GetByID(model any, id uint64) error {
	return s.DB.First(model, id).Error
}

func (s *BaseService) GetByIDWithPreload(model any, id uint64, preloads ...string) error {
	query := s.DB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.First(model, id).Error
}

func (s *BaseService) GetAll(models any, limit, offset int) error {
	query := s.DB
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	return query.Find(models).Error
}

func (s *BaseService) Update(model any, id uint64, updates any) error {
	return s.DB.Model(model).Where("id = ?", id).Updates(updates).Error
}

func (s *BaseService) Delete(model any, id uint64) error {
	return s.DB.Delete(model, id).Error
}

func (s *BaseService) GetWithPreload(model any, id uint64, preloads ...string) error {
	query := s.DB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.First(model, id).Error
}

func (s *BaseService) Count(model any) (int64, error) {
	var count int64
	err := s.DB.Model(model).Count(&count).Error
	return count, err
}

func (s *BaseService) Exists(model any, id uint64) (bool, error) {
	var count int64
	err := s.DB.Model(model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func WhereHas(db *gorm.DB, relationTable string, fk string, parentTable string, subFilter func(*gorm.DB) *gorm.DB) *gorm.DB {
	sub := db.Table(relationTable).
		Select("1").
		Where(fmt.Sprintf("%s.%s = %s.id", relationTable, fk, parentTable))

	sub = subFilter(sub)

	return db.Where("EXISTS (?)", sub)
}

func singularize(table string) string {
	if strings.HasSuffix(table, "ies") {
		// companies → company
		return strings.TrimSuffix(table, "ies") + "y"
	}

	if strings.HasSuffix(table, "es") {
		// matches → match, boxes → box
		return strings.TrimSuffix(table, "es")
	}

	if strings.HasSuffix(table, "s") && !strings.HasSuffix(table, "ss") {
		// users → user
		return strings.TrimSuffix(table, "s")
	}

	return table
}

func Paginate[T any](
	db *gorm.DB,
	model T,
	joins []string,
	filters map[string]any,
	filterFuncs map[string]func(*gorm.DB, any) *gorm.DB,
	includes map[string]func(*gorm.DB) *gorm.DB,
	searchFields []string,
	searchTerm string,
	page int,
	pageSize int,
	orderBy string,
	order string,
) ([]T, int64, error) {
	var results []T
	var total int64

	// Default pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := db.Model(model)
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, 0, err
	}
	parentTable := stmt.Table

	// Apply filters
	for field, value := range filters {
		if filterFuncs, ok := filterFuncs[field]; ok {
			query = filterFuncs(query, value)
			continue
		}

		// Check for nested field
		if strings.Contains(field, ".") {
			parts := strings.Split(field, ".")
			relationTable := parts[0]
			fieldName := parts[1]
			fk := fmt.Sprintf("%s_id", singularize(parentTable))

			query = query.Where("EXISTS (?)",
				db.Table(relationTable).
					Select("1").
					Where(fmt.Sprintf("%s.%s = %s.id", relationTable, fk, parentTable)).
					Where(fmt.Sprintf("%s.%s IN (?)", relationTable, fieldName), value),
			)
		} else {
			query = query.Where(fmt.Sprintf("%s IN (?)", field), value)
		}
	}

	// Apply search
	if searchTerm != "" && len(searchFields) > 0 {
		query = query.Scopes(func(db *gorm.DB) *gorm.DB {
			subQuery := db
			for _, col := range searchFields {
				if strings.Contains(col, ".") {
					parts := strings.Split(col, ".")
					relationTable := parts[0]
					fieldName := parts[1]
					fk := fmt.Sprintf("%s_id", parentTable)
					subQuery = subQuery.Or("EXISTS (?)",
						db.Table(relationTable).
							Select("1").
							Where(fmt.Sprintf("%s.%s = %s.id", relationTable, fk, parentTable)).
							Where(fmt.Sprintf("%s.%s LIKE ?", relationTable, fieldName), "%"+searchTerm+"%"),
					)
				} else {
					subQuery = subQuery.Or(fmt.Sprintf("%s LIKE ?", col), "%"+searchTerm+"%")
				}
			}
			return subQuery
		})
	}

	for key, fn := range includes {
		if fn != nil {
			query = query.Preload(key, fn)
		} else {
			query = query.Preload(key)
		}
	}

	for _, join := range joins {
		query = query.Joins(join)
	}

	query = query.Order(orderBy + " " + order)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if err := query.Limit(pageSize).Offset(offset).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
