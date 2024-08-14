package models

import "github.com/golang-jwt/jwt/v4"

// Meta Permissions model
type Meta struct {
	Version      int    `json:"version"`
	ManagedBy    string `json:"managed_by"`
	LastModified string `json:"last_modified"`
	Operation    string `json:"operation"`
}

type Role struct {
	Resource map[string]PermissionAttr `json:"resource"`
	Scope    []string                  `json:"scope"`
}

type Permissions struct {
	Meta  *Meta           `json:"meta"`
	Roles map[string]Role `json:"roles"`
}

type Permission map[string]PermissionAttr

type PermissionAttr struct {
	Grants []string            `json:"grants"`
	Data   map[string][]string `json:"data,omitempty"`
}

type CustomClaims struct {
	*jwt.RegisteredClaims
	ApiKeyTags
	Scope       string     `json:"scope"`
	Permissions Permission `json:"permissions"`
}

type ApiKeyTags struct {
	TenantId *string `json:"tenant_id"`
	UserRole *string `json:"user_role"`
	UserId   *string `json:"user_id"`
}
