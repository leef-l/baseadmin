// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// SystemUserDept is the golang structure for table system_user_dept.
type SystemUserDept struct {
	UserId uint64 `json:"userId" orm:"user_id" description:"用户ID"` // 用户ID
	DeptId uint64 `json:"deptId" orm:"dept_id" description:"部门ID"` // 部门ID
}
