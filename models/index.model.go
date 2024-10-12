package models

import (
	"time"
)

// Bảng users
type User struct {
	UserID      int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:255"`
	Username    string `gorm:"size:255"`
	Password    string `gorm:"size:255"`
	Email       string `gorm:"size:255"`
	PhoneNumber string `gorm:"size:255"`
	Avatar      string `gorm:"size:255"`
	RoleID      int    // Liên kết đến bảng roles
	Type        string `gorm:"size:255"`
	CreatedAt   time.Time
	SessionID   string
}

// Bảng roles
type Role struct {
	RoleID         int    `gorm:"primaryKey;autoIncrement"`
	RoleName       string `gorm:"size:255" form:"role" json:"role"`
	UserID         int    // Liên kết đến bảng users
	CreatedAt      time.Time
	LastModifiedAt time.Time
}

// Bảng permissions
type Permission struct {
	PermissionID   int    `gorm:"primaryKey;autoIncrement"`
	PermissionName string `gorm:"unique;not null" form:"permission" json:"permission"`
	UserID         int    // Liên kết đến bảng users
	CreatedAt      time.Time
	LastModifiedAt time.Time
}

type RolePermission struct {
	RoleID       int `gorm:"primaryKey"`
	PermissionID int `gorm:"primaryKey"`
}

// Bảng courses
type Course struct {
	CourseID          int    `gorm:"primaryKey;autoIncrement"`
	CourseTitle       string `gorm:"size:255" json:"courseName" form:"courseName"`
	CourseDescription string `gorm:"size:255" json:"description" form:"description"`
	UserID            int    // Liên kết đến bảng users
}

// Bảng courses_users
type CourseUser struct {
	CourseUserID int `gorm:"primaryKey;autoIncrement"`
	UserID       int `json:"user_id" form:"user_id"`
	CourseID     int `json:"course_id" form:"course_id"`
}

// Bảng lessons
type Lesson struct {
	LessonID    int    `gorm:"primaryKey"`
	LessonTitle string `gorm:"size:255"`
	UserID      int    // Liên kết đến bảng users
	CourseID    int    // Liên kết đến bảng courses
}

// Bảng posts
type Post struct {
	PostID         int    `gorm:"primaryKey"`
	PostTitle      string `gorm:"size:255"`
	PostBody       string
	LessonID       int // Liên kết đến bảng lessons
	UserID         int // Liên kết đến bảng users
	CreatedAt      time.Time
	LastModifiedAt time.Time
}

// Bảng assignments
type Assignment struct {
	AssignmentID     int    `gorm:"primaryKey"`
	AssignmentTitle  string `gorm:"size:255"`
	AssignmentBody   string
	UserID           int    // Liên kết đến bảng users
	LessonID         int    // Liên kết đến bảng lessons
	AssignmentStatus string `gorm:"size:255"`
	CreatedAt        time.Time
	DueDate          time.Time
}

// Bảng students_assignments
type StudentAssignment struct {
	StudentAssignmentID     int `gorm:"primaryKey"`
	UserID                  int `gorm:"primaryKey"`
	AssignmentID            int `gorm:"primaryKey"`
	Score                   int
	StudentAssignmentStatus string `gorm:"size:255"`
}

// Bảng file_assignments
type FileAssignment struct {
	FileAssignmentID int    `gorm:"primaryKey"`
	URL              string `gorm:"size:255"`
	AssignmentID     int    // Liên kết đến bảng assignments
	UserID           int    // Liên kết đến bảng users
}

// Bảng file_posts
type FilePost struct {
	FilePostID int    `gorm:"primaryKey"`
	URL        string `gorm:"size:255"`
	PostID     int    // Liên kết đến bảng posts
}
