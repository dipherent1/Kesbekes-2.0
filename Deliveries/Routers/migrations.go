package routers

import domains "kesbekes/Domains"

func Migrate() {
	// Migrate the schema
	DB.AutoMigrate(&domains.User{}, &domains.ChatInfo{})
}
