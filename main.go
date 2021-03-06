/*
   Copyright (C) 2017 The BlameWarrior Authors.
   This file is a part of BlameWarrior service.
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bmizerany/pat"

	"github.com/blamewarrior/users/blamewarrior"
)

func main() {

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("missing test database name (expected to be passed via ENV['DB_NAME'])")
	}

	opts := &blamewarrior.DatabaseOptions{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	db, err := blamewarrior.ConnectDatabase(dbName, opts)
	if err != nil {
		log.Fatalf("failed to establish connection with db %s using connection string %s: %s", dbName, opts.ConnectionString(), err)
	}

	handlers := NewUserHandlers(db)

	mux := pat.New()

	mux.Post("/users", http.HandlerFunc(handlers.SaveUser))
	mux.Get("/users/:nickname", http.HandlerFunc(handlers.GetUserByNickname))

	http.Handle("/", mux)

	log.Printf("blamewarrior users is running on 8080 port")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}

}
