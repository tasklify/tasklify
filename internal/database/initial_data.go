package database

// var systemRoles = []SystemRole{
// 	{Key: "admin", Title: "Administrator sistema"},
// 	{Key: "user", Title: "Uporabnik sistema"},
// }

// var projectRoles = []ProjectRole{
// 	{Key: "manager", Title: "Produktni vodja"},
// 	{Key: "master", Title: "Skrbnik metodologije"},
// 	{Key: "developer", Title: "Član razvojne skupine"},
// }

var users = []User{
	{Username: "admin", Password: "admin", FirstName: "Admin", LastName: "Uporabnik", Email: "admin@mail.com", SystemRole: SystemRoleAdmin},
	{Username: "navadni", Password: "navadni", FirstName: "Navadni", LastName: "Uporabnik", Email: "navadni@mail.com", SystemRole: SystemRoleUser},
}
