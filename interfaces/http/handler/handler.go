package handler

type UserManagementHandler struct {
	UserHandler       UserHandler
	PermissionHandler PermissionHandler
	ModuleHandler     ModuleHandler
	RoleHandler       RoleHandler
	ProfileHandler    ProfileHandler
}

type TransactionManagementHandler struct {
	LimitHandler LimitHandler
}

type Handlers struct {
	UserManagementHandler        *UserManagementHandler
	AuthHandler                  AuthHandler
	RegistrationHandler          RegistrationHandler
	TransactionManagementHandler *TransactionManagementHandler
}
