// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// SystemRoleDept is the golang structure for table system_role_dept.
type SystemRoleDept struct {
	RoleId uint64 `json:"roleId" orm:"role_id" description:"角色ID"` // 角色ID
	DeptId uint64 `json:"deptId" orm:"dept_id" description:"部门ID"` // 部门ID
}
