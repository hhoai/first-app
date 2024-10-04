package controller

import (
	"first-app/database"
	"first-app/models"
	"first-app/share"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LoginController(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func LoginPostController(c *fiber.Ctx) error {
	var p models.UserLogin
	// var errorText string

	// errorsMessage := fiber.Map{
	// 	"NameError": errorText,
	// }

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	if strings.TrimSpace(p.Username) == "" {
		errorText := "Tên dăng nhập không được để trống!"
		errorsMessage := fiber.Map{
			"NameError": errorText,
		}
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}

	var existingAccount models.User
	// Kiểm tra sự tồn tại của tài khoản
	result := database.DB.First(&existingAccount, "username = ?", p.Username)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		// Xử lý lỗi nếu có lỗi ngoài lỗi không tìm thấy
		// return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())

		// Không tìm thấy bản ghi, trả về lỗi đăng nhập thất bại
		errorText := "Sai tên đăng nhập hoặc mật khẩu!"
		errorsMessage := fiber.Map{
			"NameError": errorText,
		}
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}

	if result.Error != nil {
		if models.CheckPassword(&existingAccount, p.Password) {
			// Không tìm thấy bản ghi, trả về lỗi đăng nhập thất bại
			log.Print("loi pwd")
			errorText := "Sai tên đăng nhập hoặc mật khẩu!"
			errorsMessage := fiber.Map{
				"NameError": errorText,
			}
			return c.Status(fiber.StatusUnauthorized).JSON(errorsMessage)
			// return c.Status(fiber.StatusConflict).JSON(errorsMessage)
		} else {
			// Lỗi khác, trả về lỗi server
			errorsMessage := fiber.Map{
				"NameError": "Lỗi server",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(errorsMessage)
		}
	}

	database.DB.Where("username = ?", existingAccount.Username).First(&existingAccount)

	// Lưu thông tin người dùng vào session
	sess, err := share.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	existingAccount.SessionID = p.Username + "_session"
	sess.Set("username", existingAccount.Username)
	sess.Set("login_success", existingAccount.Username)
	sess.Set("sessionID", existingAccount.SessionID)

	database.DB.Model(&models.User{}).Where("username = ?", existingAccount.Username).Update("SessionID", existingAccount.SessionID)

	// log.Println(existingAccount.SessionID)

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	log.Println("username: ")

	return c.Redirect("/")
}

// logout
func LogoutController(c *fiber.Ctx) error {
	sess, err := share.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	sess.Destroy()
	return c.Redirect("/login")
}
