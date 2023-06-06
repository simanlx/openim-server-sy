package base_info

import crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"

type CreateDepartmentReq struct {
	*crazy_server_sdk.Department
	OperationID string `json:"operationID" binding:"required"`
}
type CreateDepartmentResp struct {
	CommResp
	Department *crazy_server_sdk.Department `json:"-"`
	Data       map[string]interface{}       `json:"data" swaggerignore:"true"`
}

type UpdateDepartmentReq struct {
	*crazy_server_sdk.Department
	DepartmentID string `json:"departmentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}
type UpdateDepartmentResp struct {
	CommResp
}

type GetSubDepartmentReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	DepartmentID string `json:"departmentID" binding:"required"`
}
type GetSubDepartmentResp struct {
	CommResp
	DepartmentList []*crazy_server_sdk.Department `json:"-"`
	Data           []map[string]interface{}       `json:"data" swaggerignore:"true"`
}

type DeleteDepartmentReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	DepartmentID string `json:"departmentID" binding:"required"`
}
type DeleteDepartmentResp struct {
	CommResp
}

type CreateOrganizationUserReq struct {
	OperationID string `json:"operationID" binding:"required"`
	*crazy_server_sdk.OrganizationUser
}
type CreateOrganizationUserResp struct {
	CommResp
}

type UpdateOrganizationUserReq struct {
	OperationID string `json:"operationID" binding:"required"`
	*crazy_server_sdk.OrganizationUser
}
type UpdateOrganizationUserResp struct {
	CommResp
}

type CreateDepartmentMemberReq struct {
	OperationID string `json:"operationID" binding:"required"`
	*crazy_server_sdk.DepartmentMember
}

type CreateDepartmentMemberResp struct {
	CommResp
}

type GetUserInDepartmentReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type GetUserInDepartmentResp struct {
	CommResp
	UserInDepartment *crazy_server_sdk.UserInDepartment `json:"-"`
	Data             map[string]interface{}             `json:"data" swaggerignore:"true"`
}

type UpdateUserInDepartmentReq struct {
	OperationID string `json:"operationID" binding:"required"`
	*crazy_server_sdk.DepartmentMember
}
type UpdateUserInDepartmentResp struct {
	CommResp
}

type DeleteOrganizationUserReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type DeleteOrganizationUserResp struct {
	CommResp
}

type GetDepartmentMemberReq struct {
	DepartmentID string `json:"departmentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}
type GetDepartmentMemberResp struct {
	CommResp
	UserInDepartmentList []*crazy_server_sdk.UserDepartmentMember `json:"-"`
	Data                 []map[string]interface{}                 `json:"data" swaggerignore:"true"`
}

type DeleteUserInDepartmentReq struct {
	DepartmentID string `json:"departmentID" binding:"required"`
	UserID       string `json:"userID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}
type DeleteUserInDepartmentResp struct {
	CommResp
}

type GetUserInOrganizationReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}

type GetUserInOrganizationResp struct {
	CommResp
	OrganizationUserList []*crazy_server_sdk.OrganizationUser `json:"-"`
	Data                 []map[string]interface{}             `json:"data" swaggerignore:"true"`
}
