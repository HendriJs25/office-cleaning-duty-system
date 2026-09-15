package seeders

import (
	"cleaning/internal/constants"
	"cleaning/internal/domain/model"
	"fmt"

	"gorm.io/gorm"
)

func seedRolePermission(tx *gorm.DB, roles map[string]model.Role, permissions map[string]model.Permission) error {
	rolePermissionSeeds := map[string][]string{
		constants.AdminRoleCode: {
			constants.PermissionUserRead,
			constants.PermissionUserCreate,
			constants.PermissionUserUpdate,
			constants.PermissionUserDeactivate,

			constants.PermissionRoleRead,
			constants.PermissionRoleCreate,
			constants.PermissionRoleUpdate,
			constants.PermissionRoleDeActivate,
			constants.PermissionRoleAssignPermission,

			constants.PermissionEmployeeRead,
			constants.PermissionEmployeeCreate,
			constants.PermissionEmployeeUpdate,
			constants.PermissionEmployeeDeactivate,
			constants.PermissionEmployeeInvitationCreate,

			constants.PermissionOfficeRead,
			constants.PermissionOfficeCreate,
			constants.PermissionOfficeUpdate,
			constants.PermissionOfficeDeactivate,

			constants.PermissionOfficeAreaRead,
			constants.PermissionOfficeAreaCreate,
			constants.PermissionOfficeAreaUpdate,
			constants.PermissionOfficeAreaDeactivate,

			constants.PermissionCleaningTaskRead,
			constants.PermissionCleaningTaskCreate,
			constants.PermissionCleaningTaskUpdate,
			constants.PermissionCleaningTaskDeactivate,

			constants.PermissionCleaningWeekRead,
			constants.PermissionCleaningWeekGenerate,
			constants.PermissionCleaningWeekReroll,
			constants.PermissionCleaningWeekEdit,
			constants.PermissionCleaningWeekDelete,

			constants.PermissionCleaningToolRead,
			constants.PermissionCleaningToolCreate,
			constants.PermissionCleaningToolUpdate,
			constants.PermissionCleaningToolDeactivate,

			constants.PermissionPurchaseRequestRead,
			constants.PermissionPurchaseRequestCreate,
			constants.PermissionPurchaseRequestUpdate,
			constants.PermissionPurchaseRequestDelete,
			constants.PermissionPurchaseRequestSubmit,
			constants.PermissionPurchaseRequestApprove,
			constants.PermissionPurchaseRequestReject,
			constants.PermissionPurchaseRequestMarkPurchased,
			constants.PermissionPurchaseRequestComplete,

			constants.PermissionNotificationSettingRead,
			constants.PermissionNotificationSettingUpdate,
			constants.PermissionNotificationSettingTest,

			constants.PermissionAuditLog,
		},
		constants.EditorRoleCode: {
			constants.PermissionEmployeeRead,
			constants.PermissionEmployeeCreate,
			constants.PermissionEmployeeUpdate,
			constants.PermissionEmployeeDeactivate,

			constants.PermissionOfficeRead,
			constants.PermissionOfficeCreate,
			constants.PermissionOfficeUpdate,

			constants.PermissionOfficeAreaRead,
			constants.PermissionOfficeAreaCreate,
			constants.PermissionOfficeAreaUpdate,
			constants.PermissionOfficeAreaDeactivate,

			constants.PermissionCleaningTaskRead,
			constants.PermissionCleaningTaskCreate,
			constants.PermissionCleaningTaskUpdate,
			constants.PermissionCleaningTaskDeactivate,

			constants.PermissionCleaningWeekRead,
			constants.PermissionCleaningWeekGenerate,
			constants.PermissionCleaningWeekReroll,
			constants.PermissionCleaningWeekEdit,

			constants.PermissionCleaningToolRead,
			constants.PermissionCleaningToolCreate,
			constants.PermissionCleaningToolUpdate,
			constants.PermissionCleaningToolDeactivate,

			constants.PermissionPurchaseRequestRead,
			constants.PermissionPurchaseRequestCreate,
			constants.PermissionPurchaseRequestUpdate,
			constants.PermissionPurchaseRequestDelete,
			constants.PermissionPurchaseRequestSubmit,
			constants.PermissionPurchaseRequestMarkPurchased,
			constants.PermissionPurchaseRequestComplete,

			constants.PermissionNotificationSettingRead,
			constants.PermissionNotificationSettingTest,
		},
		constants.ViewerRoleCode: {
			constants.PermissionEmployeeRead,
			constants.PermissionOfficeRead,
			constants.PermissionOfficeAreaRead,
			constants.PermissionCleaningTaskRead,
			constants.PermissionCleaningWeekRead,
			constants.PermissionCleaningToolRead,
			constants.PermissionPurchaseRequestRead,
		},
	}

	for roleCode, permissionCodes := range rolePermissionSeeds {
		rolePermissions := make([]model.Permission, 0, len(permissionCodes))

		role, exists := roles[roleCode]
		if !exists {
			return fmt.Errorf("role %q not found", roleCode)
		}

		for _, permissionCode := range permissionCodes {
			permission, exists := permissions[permissionCode]
			if !exists {
				return fmt.Errorf("permission %q not found", permissionCode)
			}
			rolePermissions = append(rolePermissions, permission)
		}

		if err := tx.Model(&role).Association("Permissions").Replace(rolePermissions); err != nil {
			return fmt.Errorf("replace permissions for role %q: %w", roleCode, err)
		}
	}
	return nil
}
