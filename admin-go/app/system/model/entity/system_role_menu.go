// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// SystemRoleMenu is the golang structure for table system_role_menu.
type SystemRoleMenu struct {
	RoleId uint64 `json:"roleId" orm:"role_id" description:"角色ID"` // 角色ID
	MenuId uint64 `json:"menuId" orm:"menu_id" description:"菜单ID"` // 菜单ID
}
