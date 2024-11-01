package controller

import (
	"first-app/database"
	"first-app/models"
	"first-app/share"

	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func IsPermissionSelected(selectedPermissions []int, permissionID int) bool {
	for _, p := range selectedPermissions {
		if p == permissionID {
			return true
		}
	}
	return false
}

func RoleController(c *fiber.Ctx) error {
	sess, err := share.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	log.Println(sess.Get("username"), "admin")

	var role []models.RolePermissionWithRowNumber

	result := database.DB.Table("roles").
		Select("ROW_NUMBER() OVER (ORDER BY roles.role_id) AS RowNumber, roles.role_id AS ID, roles.role_name AS role_name").
		Find(&role)

	if result.Error != nil {
		log.Println(result.Error)
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	log.Println("role", role)

	return c.Render("role/index", fiber.Map{
		"Roles": role,
		"Ctx":   c,
	}, "layouts/main")

}

// create role
func CreateRoleController(c *fiber.Ctx) error {
	var permission []models.Permission
	if err := database.DB.Find(&permission).Error; err != nil {
		log.Println(err)
	}

	data := fiber.Map{
		"Title":      "Thêm vai trò",
		"Permission": permission,
		"Ctx":        c,
	}

	return c.Render("role/create", data, "layouts/main")
}

// handle create role
func CreateRolePostController(c *fiber.Ctx) error {
	var p models.Role

	if err := c.BodyParser(&p); err != nil {
		log.Println(err.Error())
		return err
	}

	var existingRole models.Role
	// Kiểm tra sự tồn tại của tài khoản
	resultFindByRoleName := database.DB.Where("role_name = ?", p.RoleName).First(&existingRole)

	if resultFindByRoleName.Error != nil && resultFindByRoleName.Error != gorm.ErrRecordNotFound {
		// Xử lý lỗi nếu có lỗi ngoài lỗi không tìm thấy
		return c.Status(fiber.StatusInternalServerError).JSON(resultFindByRoleName.Error.Error())
	}

	if resultFindByRoleName.Error == nil {
		errorsMessage := fiber.Map{
			"RoleNameError": "Role already exists",
		}
		// Nếu không có lỗi và tìm thấy bản ghi, trả về lỗi trùng lặp
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}

	if err := database.DB.Create(&p).Error; err != nil {
		log.Println(err)
	}

	var data models.PermissionsRequest
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Lấy giá trị từ form
	permissions := data.Permissions

	log.Println("permissions", permissions)

	var newRole models.Role
	if err := database.DB.Where("role_name = ?", p.RoleName).First(&newRole).Error; err != nil {
		log.Println(err)
	}

	for _, permission := range permissions {
		log.Println("permission: ", permission)
		rolPer := models.RolePermission{
			RoleID:       newRole.RoleID,
			PermissionID: permission,
		}
		if err := database.DB.Create(rolPer).Error; err != nil {
			log.Println(err)
		}
		log.Println(rolPer)

	}

	var role []models.RolePermissionWithRowNumber

	result := database.DB.Table("roles").
		Select("ROW_NUMBER() OVER (ORDER BY roles.role_id) AS RowNumber, roles.role_name AS role_name, role_id AS ID").
		Find(&role)

	if result.Error != nil {
		log.Println(result.Error)
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	log.Println(role)
	// roles :=

	return c.Render("role/index", fiber.Map{
		"Roles": role,
		"Ctx":   c,
	}, "layouts/main")
	// ...
}

// update role
func UpdateRoleController(c *fiber.Ctx) error {
	var p models.Role

	id := c.Params("id")

	result := database.DB.First(&p, "role_id = ?", id)

	if result.Error != nil {
		log.Println(result.Error)
	}

	var roles []models.Role
	rs := database.DB.Find(&roles)
	if rs.Error != nil {
		log.Println(rs.Error)
	}

	var permission []models.Permission
	if err := database.DB.Find(&permission).Error; err != nil {
		log.Println(err)
	}

	var rolPer []models.RolePermission
	if err := database.DB.Where("role_id = ?", p.RoleID).Find(&rolPer).Error; err != nil {
		log.Println(err)
	}

	var accPer []models.Permission
	var permissionID []int
	for _, p := range rolPer {
		if err := database.DB.Where("permission_id = ?", p.PermissionID).Find(&accPer).Error; err != nil {
			log.Println(err)
		} else {
			permissionID = append(permissionID, p.PermissionID)
		}
	}

	// Tạo dữ liệu để truyền vào template
	data := fiber.Map{
		"id":           id,
		"Role":         p.RoleName,
		"Permission":   permission,
		"PermissionID": permissionID,
		"Ctx":          c,
		"Title":        "Cập nhật vai trò",
	}

	// tmpl := template.Must(template.New("").Funcs(template.FuncMap{
	// 	"isPermissionSelected": IsPermissionSelected,
	// }).ParseGlob("views/*.html"))

	// return tmpl.ExecuteTemplate(c.Response().BodyWriter(), "updateRole.html", data)

	return c.Render("role/edit", data, "layouts/main")
}

// handle update role

func UpdateRolePutController(c *fiber.Ctx) error {
	id := c.Params("id")

	var data models.PermissionsRequest
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Lấy giá trị từ form
	permissions := data.Permissions

	log.Println("permissions", permissions)

	var updateRole models.Role

	if err := database.DB.Where("role_id = ?", id).First(&updateRole).Error; err != nil {
		log.Println(err)
	}

	updateRole.RoleName = data.RoleName
	if err := database.DB.Save(&updateRole).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	var rolePermission models.RolePermission
	if err := database.DB.Where("role_id = ?", updateRole.RoleID).Delete(&rolePermission).Error; err != nil {
		log.Println(err)
	}

	for _, permission := range permissions {
		log.Println("permission: ", permission)
		rolPer := models.RolePermission{
			RoleID:       updateRole.RoleID,
			PermissionID: permission,
		}
		if err := database.DB.Create(rolPer).Error; err != nil {
			log.Println(err)
		}
		log.Println(rolPer)
	}

	var role []models.RolePermissionWithRowNumber

	result := database.DB.Table("roles").
		Select("ROW_NUMBER() OVER (ORDER BY roles.role_id) AS RowNumber, roles.role_name AS role_name, role_id AS ID").
		Find(&role)

	if result.Error != nil {
		log.Println(result.Error)
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	log.Println(role)
	// roles :=

	return c.Render("role/index", fiber.Map{
		"Roles": role,
		"Ctx":   c,
	}, "layouts/main")
}

// delete role

func DeleteRoleController(c *fiber.Ctx) error {
	var p models.Role
	var rolPer models.RolePermission

	id := c.Params("id")

	if err := database.DB.Delete(&p, "role_id = ?", id).Error; err != nil {
		log.Println(err)
	}

	if err := database.DB.Delete(&rolPer, "role_id = ?", id).Error; err != nil {
		log.Println(err)
	}

	return c.Redirect("/admin/account/role")
}
