package models

type AssignmentDetail struct {
	AssignmentID     int    `gorm:"primaryKey;autoIncrement"`
	AssignmentTitle  string `gorm:"size:255" form:"assignment_title"`
	AssignmentBody   string `form:"assignment_body"`
	UserID           int    // Liên kết đến bảng users
	LessonID         int    // Liên kết đến bảng lessons
	AssignmentStatus string `gorm:"size:255"`
	CreatedAt        string
	DueDate          string
}
