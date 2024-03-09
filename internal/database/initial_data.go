package database

var users = []User{
	{Username: "admin", Password: "admin", FirstName: "Admin", LastName: "Uporabnik", Email: "admin@mail.com", SystemRole: SystemRoleAdmin},
	{Username: "navadni", Password: "navadni", FirstName: "Navadni", LastName: "Uporabnik", Email: "navadni@mail.com", SystemRole: SystemRoleUser},
}
