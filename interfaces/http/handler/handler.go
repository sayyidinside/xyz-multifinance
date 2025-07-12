package handler

type UserManagementHandler struct {
	UserHandler       UserHandler
	PermissionHandler PermissionHandler
	ModuleHandler     ModuleHandler
	RoleHandler       RoleHandler
	ProfileHandler    ProfileHandler
	DocumentHandler   DocumentHandler
}

type TransactionManagementHandler struct {
	LimitHandler       LimitHandler
	TransactionHandler TransactionHandler
	InstallmentHandler InstallmentHandler
	PaymentHandler     PaymentHandler
}

type Handlers struct {
	UserManagementHandler        *UserManagementHandler
	AuthHandler                  AuthHandler
	RegistrationHandler          RegistrationHandler
	TransactionManagementHandler *TransactionManagementHandler
}
