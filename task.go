package main

import "time"

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Time        time.Time `json:"time"`
	Assigned_To string    `json:"assigned_to"`
}
