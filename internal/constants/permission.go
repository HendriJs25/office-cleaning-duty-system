package constants

const (
	PermissionUserRead       = "user.read"
	PermissionUserCreate     = "user.create"
	PermissionUserUpdate     = "user.update"
	PermissionUserDeactivate = "user.deactivate"

	PermissionRoleRead             = "role.read"
	PermissionRoleCreate           = "role.create"
	PermissionRoleUpdate           = "role.update"
	PermissionRoleDeActivate       = "role.deactivate"
	PermissionRoleAssignPermission = "role.assign_permission"

	PermissionEmployeeRead             = "employee.read"
	PermissionEmployeeCreate           = "employee.create"
	PermissionEmployeeUpdate           = "employee.update"
	PermissionEmployeeDeactivate       = "employee.deactivate"
	PermissionEmployeeInvitationCreate = "employee_invitation.create"

	PermissionOfficeRead       = "office.read"
	PermissionOfficeCreate     = "office.create"
	PermissionOfficeUpdate     = "office.update"
	PermissionOfficeDeactivate = "office.deactivate"

	PermissionOfficeAreaRead       = "office_area.read"
	PermissionOfficeAreaCreate     = "office_area.create"
	PermissionOfficeAreaUpdate     = "office_area.update"
	PermissionOfficeAreaDeactivate = "office_area.deactivate"

	PermissionCleaningTaskRead       = "cleaning_task.read"
	PermissionCleaningTaskCreate     = "cleaning_task.create"
	PermissionCleaningTaskUpdate     = "cleaning_task.update"
	PermissionCleaningTaskDeactivate = "cleaning_task.deactivate"

	PermissionCleaningWeekRead     = "cleaning_week.read"
	PermissionCleaningWeekGenerate = "cleaning_week.generate"
	PermissionCleaningWeekReroll   = "cleaning_week.reroll"
	PermissionCleaningWeekEdit     = "cleaning_week.edit"
	PermissionCleaningWeekDelete   = "cleaning_week.delete"

	PermissionCleaningToolRead       = "cleaning_tool.read"
	PermissionCleaningToolCreate     = "cleaning_tool.create"
	PermissionCleaningToolUpdate     = "cleaning_tool.update"
	PermissionCleaningToolDeactivate = "cleaning_tool.deactivate"

	PermissionPurchaseRequestRead          = "purchase_request.read"
	PermissionPurchaseRequestCreate        = "purchase_request.create"
	PermissionPurchaseRequestUpdate        = "purchase_request.update"
	PermissionPurchaseRequestDelete        = "purchase_request.delete"
	PermissionPurchaseRequestSubmit        = "purchase_request.submit"
	PermissionPurchaseRequestApprove       = "purchase_request.approve"
	PermissionPurchaseRequestReject        = "purchase_request.reject"
	PermissionPurchaseRequestMarkPurchased = "purchase_request.mark_purchased"
	PermissionPurchaseRequestComplete      = "purchase_request.complete"

	PermissionNotificationSettingRead   = "notification_setting.read"
	PermissionNotificationSettingUpdate = "notification_setting.update"
	PermissionNotificationSettingTest   = "notification_setting.test"

	PermissionAuditLog = "audit_log.read"
)
