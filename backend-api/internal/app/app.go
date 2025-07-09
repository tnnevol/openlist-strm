package app

import "os"

func Run() {
	InitLogger()
	db, err := InitDatabase()
	if err != nil {
		os.Exit(1)
	}
	r := RegisterRouter(db)
	r.Run(":8890")
} 
