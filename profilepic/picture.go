package profilepic

type Picture struct {
	FileName   string `json:"file_name" binding:"required" gorm:"not null"`
	ID         uint   `json:"id" binding:"required" gorm:"unique;not null"`
	USERID     uint   `json:"user_id" binding:"required" gorm:"not null"`
	URL        string `json:"url" binding:"required" gorm:"not null"`
	UploadDate string `json:"upload_date" binding:"required" gorm:"not null"`
}
