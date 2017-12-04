package main

import (
	"fmt"
	"time"
)

type User struct {
	ID           int      `json:"id" storm:"id,increment"`
	Name         string   `json:"name" storm:"index"`
	CanRead      bool     `json:"can_read"`
	CanPost      bool     `json:"can_post"`
	IsAdmin      bool     `json:"is_admin"`
	RegionFilter []string `json:"region_filter"`
}

func (u User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s, CanRead: %t, CanPost: %t, IsAdmin: %t}",
		u.ID, u.Name, u.CanRead, u.CanPost, u.IsAdmin)
}

func dbGetUsers() ([]User, error) {
	var users []User
	err := db.All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func dbDeleteUser() error {
	return nil
}

func dbUpdateUser() () {}

func dbCreateUser(name string, canRead bool, canPost bool,
	isAdmin bool) error {
	user := User{
		Name:    name,
		CanRead: canRead,
		CanPost: canPost,
		IsAdmin: isAdmin,
	}
	err := db.Save(&user)
	if err != nil {
		return err
	}
	return nil
}

type Timer struct {
	ID            int       `json:"id" storm:"id,increment"`
	Region        string    `json:"region" storm:"index"`
	System        string    `json:"system" storm:"index"`
	StructureType string    `json:"structure_type"`
	ReinforceType string    `json:"reinforce_type"`
	Comment       string    `json:"comment"`
	DaysLeft      int       `json:"days_left"`
	HoursLeft     int       `json:"hours_left"`
	MinutesLeft   int       `json:"minutes_left"`
	CreatedAt     time.Time `json:"created_at"`
	StartAt       time.Time `json:"start_at"`
}

func (t Timer) String() string {
	return fmt.Sprintf("Timer{ID: %d, Region: %s, System: %s, "+
		"StructureType: %s, ReinforceType: %s, Comment: %s, StartAt: %v}",
		t.ID, t.Region, t.System, t.StructureType, t.ReinforceType,
		t.Comment, t.StartAt)
}

func dbGetTimers() ([]Timer, error) {
	var timers []Timer
	err := db.All(&timers)
	if err != nil {
		return nil, err
	}
	return timers, nil
}

func dbDeleteTimer(id int) error {
	var timer Timer

	err := db.One("ID", id, &timer)
	if err != nil {
		return err
	}

	err = db.DeleteStruct(&timer)
	if err != nil {
		return err
	}

	return nil
}

func dbUpdateTimer(id int, region string, system string, structureType string,
	reinforceType string, comment string, daysLeft int, hoursLeft int,
	minutesLeft int) error {
	var timer Timer

	err := db.One("ID", id, &timer)
	if err != nil {
		return err
	}

	startAt := time.Now().UTC().Add(time.Hour*time.Duration(daysLeft)*24 +
		time.Hour*time.Duration(hoursLeft) +
		time.Minute*time.Duration(minutesLeft))
	err = db.Update(&Timer{
		ID:   id,
		Region:        region,
		System:        system,
		StructureType: structureType,
		ReinforceType: reinforceType,
		Comment:       comment,
		DaysLeft:      daysLeft,
		HoursLeft:     hoursLeft,
		MinutesLeft:   minutesLeft,
		StartAt:       startAt,
	})

	return nil
}

func dbCreateTimer(region string, system string, structureType string,
	reinforceType string, comment string, daysLeft int, hoursLeft int,
	minutesLeft int) error {
	createdAt := time.Now().UTC()
	startAt := time.Now().UTC().Add(time.Hour*time.Duration(daysLeft)*24 +
		time.Hour*time.Duration(hoursLeft) +
		time.Minute*time.Duration(minutesLeft))
	timer := Timer{
		Region:        region,
		System:        system,
		StructureType: structureType,
		ReinforceType: reinforceType,
		Comment:       comment,
		DaysLeft:      daysLeft,
		HoursLeft:     hoursLeft,
		MinutesLeft:   minutesLeft,
		CreatedAt:     createdAt,
		StartAt:       startAt,
	}
	err := db.Save(&timer)
	if err != nil {
		return err
	}
	return nil
}
