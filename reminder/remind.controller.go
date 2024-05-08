package reminder

import (
	"barbershop-backend/controllers"
	"barbershop-backend/initializers"
	"barbershop-backend/models"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RemindController struct {
	DB             *gorm.DB
	UserController *controllers.UserController
}

func NewRemindController(DB *gorm.DB, UserController *controllers.UserController) RemindController {
	return RemindController{DB, UserController}
}

func (rc *RemindController) CheckAndSendBirthdayReminders() {
	// Fetch user list fom the database (example)
	uc := rc.UserController

	users, err := uc.GetAllUsers()

	if err != nil {
		log.Printf("Error fetching user list from the database: %v\n", err)
		return
	}

	// Get today's date
	today := time.Now().UTC()

	// Iterate through users and check for upcoming birthdays
	for _, user := range users {
		birthday := user.Birthday.UTC()

		// Calculate reminder date (5 days before birthday)
		reminderDate := birthday.AddDate(0, 0, -5).UTC()

		fmt.Println("birthday: year, month, day", birthday.Year(), birthday.Month(), birthday.Day())
		fmt.Println("today: year, month, day", today.Year(), today.Month(), today.Day())

		// Check if today's date matches the reminder date
		if (today.Year() == reminderDate.Year() && today.Month() == reminderDate.Month() && today.Day() == reminderDate.Day()) ||
			(birthday.Year() == today.Year() && birthday.Month() == today.Month() && birthday.Day() == today.Day()) {
			// Send birthday reminder email
			err := SendBirthdayReminder(user.Email)
			if err != nil {
				log.Printf("Error sending birthday reminder email to %s (%s): %v\n", user.Name, user.Email, err)
			}
		}
	}
}

func SendBirthdayReminder(email string) error {
	config, cError := initializers.LoadConfig(".")
	if cError != nil {
		log.Fatal("🚀 Could not load environment variables", cError)
	}
	senderEmail := config.SenderEmail
	senderPassword := config.SenderEmailKey

	// Compose the email
	subject := "Chúc mừng sinh nhật quý khách hàng"
	body := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Happy Birthday!</title>
	</head>
	<body>
		<h2>Chúc mừng sinh nhật!</h2>
		<p>
			Chúng tôi muốn bắt đầu chuỗi ngày sinh nhật của bạn với lời chúc tốt đẹp nhất và một lời mời đặc biệt.
		</p>
		<p>
			Trong vòng 5 ngày tới, chúng tôi hân hạnh mời bạn đến tham gia vào chương trình cắt tóc miễn phí trong ngày sinh nhật của mình tại cửa hàng của chúng tôi.
		</p>
		<p>
			Chúng tôi mong muốn chia sẻ niềm vui và tri ân đến quý khách hàng đã ủng hộ chúng tôi suốt thời gian qua.
		</p>
		<p>
			Hãy để chúng tôi tạo ra một trải nghiệm đẳng cấp và đầy phong cách cho bạn trong ngày đặc biệt này!
		</p>
		<p>
			Chúc quý khách có một ngày sinh nhật thật vui vẻ và ý nghĩa!
		</p>
		<p>
			Quý khách vui lòng đặt lịch tại <a href="https://roybarbershop.com/dat-lich">roybarbershop.com/dat-lich</a>.
		</p>
		<p>
			Trân trọng,<br>
			Roy Barber shop
		</p>
	</body>
	</html>
	`
	// Authentication
	auth := smtp.PlainAuth("", senderEmail, senderPassword, "smtp.gmail.com")

	// Message content
	message := fmt.Sprintf("From: %s\r\n", "Roybarbershop <roybarbershop>") +
		fmt.Sprintf("To: %s\r\n", email) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" + // Specify content type as HTML
		fmt.Sprintf("\r\n%s", body)

	// Send email using SMTP
	err := smtp.SendMail("smtp.gmail.com:587", auth, senderEmail, []string{email}, []byte(message))

	if err != nil {
		return err
	}
	return nil
}

func (rc *RemindController) AdminSendBirthdayReminder(ctx *gin.Context) {

	var user *models.UserResponse

	var payload *models.ReminderRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fmt.Println("user:", user)

	err := SendBirthdayReminder(payload.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Cannot set email birthday reminder"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Send email birthday reminder success"})
}
