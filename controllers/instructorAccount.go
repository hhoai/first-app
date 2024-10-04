package controller

import (
	"first-app/database"
	"first-app/models"
	"first-app/share"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// instructor
func GetDataInstructor(c *fiber.Ctx) error {
	draw, _ := strconv.Atoi(c.Query("draw"))
	start, _ := strconv.Atoi(c.Query("start"))
	length, _ := strconv.Atoi(c.Query("length"))
	searchValue := c.Query("search[value]")

	var totalRecords int64
	var filteredRecords int64
	var users []models.User

	// Get total number of records
	database.DB.Table("users").Count(&totalRecords)

	// Apply search filter if provided
	query := database.DB.Table("users")

	// search theo cot
	if searchValue != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+searchValue+"%", "%"+searchValue+"%")
	}

	// Get total number of filtered records
	query.Count(&filteredRecords)

	// Apply pagination
	// database.DB.Table("user_entities").Offset(start).Limit(length).Find(&users)
	database.DB.Offset(start).Limit(length).Find(&users)

	log.Println(users)
	var account []models.UserWithRowNumber
	result := query.
		Where("type = 'instructor'").
		Joins("INNER JOIN roles ON users.role_id = roles.role_id").
		Select("ROW_NUMBER() OVER (ORDER BY ID) AS RowNumber, users.user_id AS ID, name, username, email, phone_number, roles.role_name AS role_name").
		Offset(start).
		Limit(length).
		Find(&account)
	if result.Error != nil {
		log.Println(result.Error)
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	log.Println("acc: ", account)

	// Prepare response
	response := map[string]interface{}{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            account,
	}
	// Return JSON response
	return c.JSON(response)
}

func InstructorAccountController(c *fiber.Ctx) error {
	return c.Render("accountInstructor", fiber.Map{
		// "SaleData": sales,
		"Ctx": c,
	}, "layouts/main")
}

// create account
func CreateInstructorAccountController(c *fiber.Ctx) error {
	var roles []models.Role
	result := database.DB.Find(&roles)
	if result.Error != nil {
		log.Println(result.Error)
	}

	data := fiber.Map{
		"Role": roles,
		"Ctx":  c,
	}

	return c.Render("createInstructorAccount", data, "layouts/main")
}

// post
func CreateInstructorAccountPostController(c *fiber.Ctx) error {
	var p models.Account

	if err := c.BodyParser(&p); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(p.RoleID)
	nameError := validateUsername(p.Username)

	errorsMessage := fiber.Map{
		"NameError":     nameError,
		"PasswordError": validatePassword(p.Password),
		"EmailError":    validateEmail(p.Email),
		"AddressError":  validateAddress(p.Name),
		"PhoneError":    validatePhoneNumber(p.PhoneNumber),
	}

	var existingAccount models.User
	// Kiểm tra sự tồn tại của tài khoản
	resultFindByUserName := database.DB.Where("name = ? ", p.Username).First(&existingAccount)
	resultFindByEmail := database.DB.Where("email = ?", p.Email).First(&existingAccount)

	if resultFindByUserName.Error != nil && resultFindByUserName.Error != gorm.ErrRecordNotFound {
		// Xử lý lỗi nếu có lỗi ngoài lỗi không tìm thấy
		return c.Status(fiber.StatusInternalServerError).JSON(resultFindByUserName.Error.Error())
	}

	if resultFindByEmail.Error != nil && resultFindByEmail.Error != gorm.ErrRecordNotFound {
		// Xử lý lỗi nếu có lỗi ngoài lỗi không tìm thấy
		return c.Status(fiber.StatusInternalServerError).JSON(resultFindByEmail.Error.Error())
	}

	if validateUsername(p.Username) != "" || validatePassword(p.Password) != "" || validateEmail(p.Email) != "" || validatePhoneNumber(p.PhoneNumber) != "" || validateAddress(p.Name) != "" {
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}

	if resultFindByUserName.Error == nil {
		errorsMessage := fiber.Map{
			"NameError": "UserName already exists",
		}
		// Nếu không có lỗi và tìm thấy bản ghi, trả về lỗi trùng lặp
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}

	if resultFindByEmail.Error == nil {
		errorsMessage := fiber.Map{
			"EmailError": "Email already exists",
		}
		// Nếu không có lỗi và tìm thấy bản ghi, trả về lỗi trùng lặp
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}
	// Tạo tài khoản mới
	newAccount := models.User{
		Username:    p.Username,
		Password:    p.Password,
		Email:       p.Email,
		Name:        p.Name,
		PhoneNumber: p.PhoneNumber,
		RoleID:      p.RoleID,
		Type:        "instructor",
	}

	if err := database.DB.Create(&newAccount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	sess, err := share.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	username := sess.Get("username")

	log.Println(username, "create", p.RoleID)

	return c.Redirect("/admin/account-instructor")
	// ...
}

// update account
func UpdateInstructorAccountController(c *fiber.Ctx) error {
	var p models.User

	id := c.Params("id")

	result := database.DB.First(&p, "user_id = ?", id)

	if result.Error != nil {
		log.Println(result.Error)
	}

	var roles []models.Role

	rs := database.DB.Find(&roles)
	if rs.Error != nil {
		log.Println(rs.Error)
	}

	// Tạo dữ liệu để truyền vào template
	data := fiber.Map{
		"ID":          id,
		"Username":    p.Username,
		"Email":       p.Email,
		"PhoneNumber": p.PhoneNumber,
		"Name":        p.Name,
		"RoleID":      p.RoleID,
		"Roles":       roles,
		"Ctx":         c,
	}
	return c.Render("updateInstructorAccount", data, "layouts/main")
}

// put
func UpdateInstructorAccountPutController(c *fiber.Ctx) error {
	id := c.Params("id")
	var p models.Account

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	nameError := validateUsername(p.Username)

	errorsMessage := fiber.Map{
		"NameError":    nameError,
		"EmailError":   validateEmail(p.Email),
		"AddressError": validateAddress(p.Name),
		"PhoneError":   validatePhoneNumber(p.PhoneNumber),
	}

	var existingAccount models.User
	// Kiểm tra sự tồn tại của tài khoản
	resultFindByUserName := database.DB.Where("name = ? ", p.Username).First(&existingAccount)
	resultFindByEmail := database.DB.Where("email = ?", p.Email).First(&existingAccount)

	if resultFindByUserName.Error != nil && resultFindByUserName.Error != gorm.ErrRecordNotFound {
		// Xử lý lỗi nếu có lỗi ngoài lỗi không tìm thấy
		return c.Status(fiber.StatusInternalServerError).JSON(resultFindByUserName.Error.Error())
	}

	if resultFindByEmail.Error != nil && resultFindByEmail.Error != gorm.ErrRecordNotFound {
		// Xử lý lỗi nếu có lỗi ngoài lỗi không tìm thấy
		return c.Status(fiber.StatusInternalServerError).JSON(resultFindByEmail.Error.Error())
	}

	if validateUsername(p.Username) != "" || validateEmail(p.Email) != "" || validatePhoneNumber(p.PhoneNumber) != "" || validateAddress(p.Name) != "" {
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}

	var account models.User
	rs := database.DB.First(&account, "user_id = ?", id)

	if rs.Error != nil {
		log.Println(rs.Error)
	}

	account.Username = p.Username
	account.Email = p.Email
	account.Name = p.Name
	account.PhoneNumber = p.PhoneNumber
	account.RoleID = p.RoleID

	if err := database.DB.Save(&account).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	sess, err := share.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	username := sess.Get("username")

	log.Println(username, "update")
	var users []models.User

	err1 := database.DB.Find(&users).Error
	if err1 != nil {
		log.Println(err1)
	}

	user := fiber.Map{
		"Users": users,
		"Ctx":   c,
	}
	return c.Render("accountInstructor", user, "layouts/main")
	// return c.Redirect("/admin")
	// ...
}

// delete

func DeleteInstructorAccountController(c *fiber.Ctx) error {
	var p models.User

	id := c.Params("id")

	err := database.DB.Delete(&p, "user_id = ?", id).Error

	if err != nil {
		log.Println(err)
	}

	sess, err := share.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	username := sess.Get("username")

	log.Println(username, "delete")

	var users []models.User

	err1 := database.DB.Find(&users).Error
	if err1 != nil {
		log.Println(err1)
	}

	user := fiber.Map{
		"Users": users,
		"Ctx":   c,
	}
	return c.Render("accountInstructor", user, "layouts/main")
}

func DeleteMultipleInsructorAccounts(c *fiber.Ctx) error {
	// Định nghĩa một struct để nhận mảng ID từ client
	type RequestBody struct {
		AccountIDs []int `json:"account_id"`
	}

	var reqBody RequestBody

	// Parse JSON từ request body
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu không hợp lệ",
		})
	}
	// Kiểm tra nếu không có account ID nào được gửi lên
	if len(reqBody.AccountIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Vui lòng chọn bản ghi cần xóa",
		})
	}

	// Xóa nhiều bản ghi trong bảng user_entities dựa vào mảng ID
	if err := database.DB.Where("id IN ?", reqBody.AccountIDs).Delete(&models.User{}).Error; err != nil {
		log.Println("Đã xảy ra lỗi:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Đã xảy ra lỗi",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Xóa thành công",
	})
}
