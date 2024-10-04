package routers

import domains "kesbekes/Domains"

func Migrate() {
	// Migrate the schema
	DB.AutoMigrate(&domains.User{}, &domains.ChatInfo{})
	// DB.Migrator().AddColumn(&domains.User{}, "Preferenses")
	// DB.Migrator().AddColumn(&domains.User{}, "UserID")
	// DB.Migrator().AddColumn(&domains.User{}, "Chats")
	// DB.Migrator().AddColumn(&domains.ChatInfo{}, "Users")

	// // Create the many-to-many relationship
	// DB.SetupJoinTable(&domains.User{}, "Chats", &domains.ChatInfo{})
	// DB.SetupJoinTable(&domains.ChatInfo{}, "Users", &domains.User{})

}
