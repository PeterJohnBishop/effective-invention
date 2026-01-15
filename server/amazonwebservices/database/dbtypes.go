package database

type User struct {
	ID        string `json:"id" dynamodbav:"id"`
	Name      string `json:"name" dynamodbav:"name"`
	Email     string `json:"email" dynamodbav:"email"`
	Password  string `json:"-" dynamodbav:"password"` // "-" hides password from JSON responses
	CreatedAt int64  `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt int64  `json:"updatedAt" dynamodbav:"updatedAt"`
}

type UserFile struct {
	User      string `json:"user" dynamodbav:"user"` // partition key
	ID        string `json:"id" dynamodbav:"id"`     // sort key
	FileKey   string `json:"filekey" dynamodbav:"fileKey"`
	CreatedAt int64  `json:"createdAt" dynamodbav:"createdAt"`
}
