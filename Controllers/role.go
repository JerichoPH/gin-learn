package Controllers

type Role struct {
	Name string `form:"name" binding:"required"`
}

type RoleController struct {
	role Role
}

func (cls *RoleController) SetRole(name string) *RoleController {
	role := Role{Name: name}
	cls.role = role
	return cls
}

func (cls *RoleController) GetRole() Role {
	return cls.role
}
