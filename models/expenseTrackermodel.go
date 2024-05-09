package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)
type ExpenseTrackModel struct{

	Id				primitive.ObjectID		`bson:"_id"`
	UserExpense		primitive.ObjectID		`bson:"user_expense"`
	ExpenseType		*string					`bson:"expense_type" json:"expense_type"  validate:"required"`
	ExpensePrice	*int					`bson:"expense_price" json:"expense_price" validate:"required"`
	CreateAt		time.Time               `bson:"created_at"`
	UpdatedAt		time.Time				`bson:"updated_at"`
}	