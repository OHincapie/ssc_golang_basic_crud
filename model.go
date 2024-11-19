package main

type User struct {
	ID          uint64 `bson:"identification" json:"identification"`
	Name        string `bson:"name" json:"name"`
	LastName    string `bson:"lastName" json:"last_name"`
	PublicForce string `bson:"publicForce" json:"public_force"`
	Range       string `bson:"range" json:"range"`
	ForceID     int    `bson:"forceId" json:"force_id"`
	Email       string `bson:"email" json:"email"`
	Password    string `bson:"password" json:"password"`
}
