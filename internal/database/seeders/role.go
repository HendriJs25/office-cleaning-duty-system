package seeders

import (
	"cleaning/internal/constants"
	"cleaning/internal/domain/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func seedRoles(tx *gorm.DB) (map[string]model.Role, error) {
	roleSeeds := []model.Role{
		{
			Code: constants.AdminRoleCode,
			Name: "管理者",
		},
		{
			Code: constants.EditorRoleCode,
			Name: "編集者",
		},
		{
			Code: constants.ViewerRoleCode,
			Name: "閲覧者",
		},
	}

	roles := make(map[string]model.Role, len(roleSeeds))

	for _, roleSeed := range roleSeeds {
		role, err := findOrCreateRole(tx, roleSeed)
		if err != nil {
			return nil, err
		}
		roles[role.Code] = role
	}
	return roles, nil
}

func findOrCreateRole(tx *gorm.DB, roleSeed model.Role) (model.Role, error) {
	var role model.Role

	err := tx.Where("code = ?", roleSeed.Code).First(&role).Error

	switch {
	case err == nil:
		if !role.IsActive {
			if err := tx.Model(&role).Update("is_active", true).Error; err != nil {
				return model.Role{}, fmt.Errorf("restore role %q: %w", roleSeed.Code, err)
			}
		}
		return role, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		role = roleSeed
		if err := tx.Create(&role).Error; err != nil {
			return model.Role{}, fmt.Errorf("create role %q: %w", roleSeed.Code, err)
		}
		return role, nil
	default:
		return model.Role{}, fmt.Errorf("find role %q: %w", roleSeed.Code, err)
	}
}
