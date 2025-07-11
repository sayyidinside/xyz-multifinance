package handler

type UserManagementHandler struct {
	UserHandler       UserHandler
	PermissionHandler PermissionHandler
	ModuleHandler     ModuleHandler
	RoleHandler       RoleHandler
	ProfileHandler    ProfileHandler
}

type TransactionManagementHandler struct {
	LimitHandler       LimitHandler
	TransactionHandler TransactionHandler
	InstallmentHandler InstallmentHandler
}

type Handlers struct {
	UserManagementHandler        *UserManagementHandler
	AuthHandler                  AuthHandler
	RegistrationHandler          RegistrationHandler
	TransactionManagementHandler *TransactionManagementHandler
}
