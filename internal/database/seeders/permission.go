package seeders

import (
	"cleaning/internal/constants"
	"cleaning/internal/domain/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func seedPermissions(tx *gorm.DB) (map[string]model.Permission, error) {
	permissionSeeds := []model.Permission{
		{Code: constants.PermissionUserRead, Name: "ユーザー閲覧"},
		{Code: constants.PermissionUserCreate, Name: "ユーザー作成"},
		{Code: constants.PermissionUserUpdate, Name: "ユーザー更新"},
		{Code: constants.PermissionUserDeactivate, Name: "ユーザー無効化"},

		{Code: constants.PermissionRoleRead, Name: "ロール閲覧"},
		{Code: constants.PermissionRoleCreate, Name: "ロール作成"},
		{Code: constants.PermissionRoleUpdate, Name: "ロール更新"},
		{Code: constants.PermissionRoleDeActivate, Name: "ロール無効化"},
		{Code: constants.PermissionRoleAssignPermission, Name: "ロール権限設定"},

		{Code: constants.PermissionEmployeeRead, Name: "従業員閲覧"},
		{Code: constants.PermissionEmployeeCreate, Name: "従業員作成"},
		{Code: constants.PermissionEmployeeUpdate, Name: "従業員更新"},
		{Code: constants.PermissionEmployeeDeactivate, Name: "従業員無効化"},
		{Code: constants.PermissionEmployeeInvitationCreate, Name: "従業員招待作成"},

		{Code: constants.PermissionOfficeRead, Name: "オフィス閲覧"},
		{Code: constants.PermissionOfficeCreate, Name: "オフィス作成"},
		{Code: constants.PermissionOfficeUpdate, Name: "オフィス更新"},
		{Code: constants.PermissionOfficeDeactivate, Name: "オフィス無効化"},

		{Code: constants.PermissionOfficeAreaRead, Name: "オフィスエリア閲覧"},
		{Code: constants.PermissionOfficeAreaCreate, Name: "オフィスエリア作成"},
		{Code: constants.PermissionOfficeAreaUpdate, Name: "オフィスエリア更新"},
		{Code: constants.PermissionOfficeAreaDeactivate, Name: "オフィスエリア無効化"},

		{Code: constants.PermissionCleaningTaskRead, Name: "清掃タスク閲覧"},
		{Code: constants.PermissionCleaningTaskCreate, Name: "清掃タスク作成"},
		{Code: constants.PermissionCleaningTaskUpdate, Name: "清掃タスク更新"},
		{Code: constants.PermissionCleaningTaskDeactivate, Name: "清掃タスク無効化"},

		{Code: constants.PermissionCleaningWeekRead, Name: "清掃当番週閲覧"},
		{Code: constants.PermissionCleaningWeekGenerate, Name: "清掃当番週生成"},
		{Code: constants.PermissionCleaningWeekReroll, Name: "清掃当番再抽選"},
		{Code: constants.PermissionCleaningWeekEdit, Name: "清掃当番週編集"},
		{Code: constants.PermissionCleaningWeekDelete, Name: "清掃当番週削除"},

		{Code: constants.PermissionCleaningToolRead, Name: "清掃用具閲覧"},
		{Code: constants.PermissionCleaningToolCreate, Name: "清掃用具作成"},
		{Code: constants.PermissionCleaningToolUpdate, Name: "清掃用具更新"},
		{Code: constants.PermissionCleaningToolDeactivate, Name: "清掃用具無効化"},

		{Code: constants.PermissionPurchaseRequestRead, Name: "購入申請閲覧"},
		{Code: constants.PermissionPurchaseRequestCreate, Name: "購入申請作成"},
		{Code: constants.PermissionPurchaseRequestUpdate, Name: "購入申請更新"},
		{Code: constants.PermissionPurchaseRequestDelete, Name: "購入申請削除"},
		{Code: constants.PermissionPurchaseRequestSubmit, Name: "購入申請提出"},
		{Code: constants.PermissionPurchaseRequestApprove, Name: "購入申請承認"},
		{Code: constants.PermissionPurchaseRequestReject, Name: "購入申請却下"},
		{Code: constants.PermissionPurchaseRequestMarkPurchased, Name: "購入済み登録"},
		{Code: constants.PermissionPurchaseRequestComplete, Name: "購入申請完了"},

		{Code: constants.PermissionNotificationSettingRead, Name: "通知設定閲覧"},
		{Code: constants.PermissionNotificationSettingUpdate, Name: "通知設定更新"},
		{Code: constants.PermissionNotificationSettingTest, Name: "通知テスト実行"},

		{Code: constants.PermissionAuditLog, Name: "監査ログ閲覧"},
	}

	permissions := make(map[string]model.Permission, len(permissionSeeds))

	for _, permissionSeed := range permissionSeeds {
		permission, err := findOrCreatePermission(tx, permissionSeed)
		if err != nil {
			return nil, err
		}

		permissions[permission.Code] = permission
	}

	return permissions, nil
}

func findOrCreatePermission(tx *gorm.DB, permissionSeed model.Permission) (model.Permission, error) {
	var permission model.Permission

	err := tx.Where("code = ?", permissionSeed.Code).First(&permission).Error

	switch {
	case err == nil:
		updates := map[string]any{}

		if !permission.IsActive {
			updates["is_active"] = true
		}

		if permission.Name != permissionSeed.Name {
			updates["name"] = permissionSeed.Name
		}

		if len(updates) > 0 {
			if err := tx.Model(&permission).Updates(updates).Error; err != nil {
				return model.Permission{}, fmt.Errorf("restore/update permission: %w", err)
			}
		}
		permission.Name = permissionSeed.Name
		return permission, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		permission = permissionSeed

		if err := tx.Create(&permission).Error; err != nil {
			return model.Permission{}, fmt.Errorf("create permission: %w", err)
		}

		return permission, nil
	default:
		return model.Permission{}, fmt.Errorf("find permission: %w", err)
	}
}
