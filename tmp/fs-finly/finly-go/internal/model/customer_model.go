package model

type Customer struct {
	ID      string `bson:"_id,omitempty" json:"id"`
	Name    string `bson:"name" json:"name"`
	Email   string `bson:"email" json:"email"`
	Phone   string `bson:"phone" json:"phone"`
	Address string `bson:"address" json:"address"`
	Owner   string `bson:"owner" json:"owner"`
}
