package handler

type UserManagementHandler struct {
	UserHandler       UserHandler
	PermissionHandler PermissionHandler
	ModuleHandler     ModuleHandler
	RoleHandler       RoleHandler
	ProfileHandler    ProfileHandler
}

type Handlers struct {
	UserManagementHandler *UserManagementHandler
	AuthHandler           AuthHandler
	RegistrationHandler   RegistrationHandler
}
