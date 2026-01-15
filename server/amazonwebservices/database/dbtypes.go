package database

type User struct {
	ID        string `dynamodbav:"id"`
	Name      string `dynamodbav:"name"`
	Email     string `dynamodbav:"email"`
	Password  string `dynamodbav:"password"`
	CreatedAt int64  `dynamodbav:"createdAt"`
	UpdatedAt int64  `dynamodbav:"updatedAt"`
}

type UserFile struct {
	UserID    string `dynamodbav:"userId"` // partition key
	ID        string `dynamodbav:"fileId"` // sort key
	FileKey   string `dynamodbav:"fileKey"`
	CreatedAt int64  `dynamodbav:"createdAt"`
}
