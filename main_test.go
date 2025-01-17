package main

import (
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func isEqualStringSlice(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestInsertUser(t *testing.T) {
	fixedTime := time.Date(2025, 1, 17, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		u        *user
		expected string
	}{
		{&user{
			typeID:    3,
			name:      "BATMAN",
			email:     "batman@gotham.city",
			pwHash:    "bruc3W47N3",
			regDate:   fixedTime,
			lastLogin: fixedTime},
			"nil"},
		{&user{ // duplicate user email
			typeID:    3,
			name:      "BATMAN",
			email:     "batman@gotham.city",
			pwHash:    "dickGrayson",
			regDate:   fixedTime,
			lastLogin: fixedTime},
			"UNIQUE constraint failed: users.email"},
		{&user{ // invalid typeID
			typeID:    666,
			name:      "JOKER",
			email:     "j0k3r@gotham.city",
			pwHash:    "alfredPennyless",
			regDate:   fixedTime,
			lastLogin: fixedTime},
			"FOREIGN KEY constraint failed"},
		{&user{
			typeID:    1,
			name:      "Superman",
			email:     "superman@metropolis.city",
			pwHash:    "clarkKent",
			regDate:   fixedTime,
			lastLogin: fixedTime},
			"nil"},
	}
	db.deleteAllUsers()
	for i, tc := range testCases {
		err := db.insertUser(tc.u)
		result := "nil"
		if err != nil {
			result = err.Error()
		}
		if result != tc.expected {
			t.Errorf("Case %d: expected error %v, got %v\n", i, tc.expected, result)
		}
	}
}

func TestSelectFieldFromTable(t *testing.T) {
	testCases := []struct {
		field    string
		table    string
		expected []string
	}{
		{"name",
			"categories",
			[]string{"General", "golang", "html", "css", "sqlite3"}},
		{"email",
			"users",
			[]string{"batman@gotham.city", "superman@metropolis.city"}},
		{"title", // empty result
			"posts",
			[]string{}},
	}
	for i, tc := range testCases {
		results, err := db.selectFieldFromTable(tc.field, tc.table)
		if !isEqualStringSlice(results, tc.expected) && err == nil {
			t.Errorf("case%v: failed. %v\n", i, results)
		}
	}
}

func TestUpdateSelectUserByEmail(t *testing.T) {
	fixedTime := time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC)
	u := &user{
		name:      "Superman",
		email:     "superman@metropolis.city",
		pwHash:    "kalEl",
		lastLogin: fixedTime}

	err := db.updateUser(u)
	u2, err2 := db.selectUserByEmail(u.email)
	if (err != nil || err2 != nil) &&
		u.name != u2.name &&
		u2.lastLogin == fixedTime {
		t.Errorf("Case: failed.\n")
	}
}

func TestInsertSession(t *testing.T) {
	u, _ := db.selectUserByEmail("batman@gotham.city")
	testCases := []struct {
		s        *session
		expected string
	}{
		{&session{
			"0012",
			u.id,
			true,
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 00, 30, 0, time.UTC),
		}, "nil"},
		{&session{ // duplicated session id
			"0012",
			1,
			true,
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 00, 30, 0, time.UTC),
		}, "UNIQUE constraint failed: sessions.id"},
		{&session{ // invlaid user id
			"0015",
			10,
			true,
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 00, 30, 0, time.UTC),
		}, "FOREIGN KEY constraint failed"},
	}
	db.deleteAllSessions()
	for i, tc := range testCases {
		err := db.insertSession(tc.s)
		result := "nil"
		if err != nil {
			result = err.Error()
		}
		if result != tc.expected {
			t.Errorf("Case %d: expected error %v, got %v\n", i, tc.expected, result)
		}
	}
}
