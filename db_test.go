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
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, result)
		}
	}
}

func TestSelectUserByEmail(t *testing.T) {
	testCases := []struct {
		email    string
		name     string
		expected string
	}{
		{"j0k3r@gotham.city", //not registered
			"",
			"nil"},
		{"superman@metropolis.city", // valid
			"Superman",
			"nil"},
		{"batman@gotham.city", // valid
			"BATMAN",
			"nil"},
	}
	for i, tc := range testCases {
		u, err := db.selectUserByEmail(tc.email)
		result := "nil"
		if err != nil {
			result = err.Error()
		}
		if result != tc.expected {
			t.Errorf("Case%v: failed %v\n", i, err)
			return
		} else {
			t.Log(u)
		}
	}
}

func TestSelectFieldFromTable(t *testing.T) {
	testCases := []struct {
		field    string
		table    string
		expected []string
	}{
		{"name", // check categories
			"categories",
			[]string{"General", "golang", "html", "css", "sqlite3"}},
		{"email", // check user emails
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
		} else {
			t.Log(results)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	u, _ := db.selectUserByEmail("superman@metropolis.city")
	fixedTime := time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC)
	u.name = "Ultimate Superman"
	u.pwHash = "kalEl"
	u.lastLogin = fixedTime

	err := db.updateUser(u)
	if err != nil {
		t.Errorf("Case: updateUser() failed.%s\n", err)
		return
	}
	u2, err2 := db.selectUserByEmail(u.email)
	if (err2 != nil || u2 == nil) &&
		u.name != u2.name &&
		u2.lastLogin == fixedTime {
		t.Errorf("Case: failed.\n")
	} else {
		t.Log(u2)
	}
}

func TestInsertSession(t *testing.T) {
	u, _ := db.selectUserByEmail("batman@gotham.city")
	u2, _ := db.selectUserByEmail("superman@metropolis.city")
	testCases := []struct {
		s        *session
		expected string
	}{
		{&session{ // insert valid session
			"0012",
			u.id,
			true,
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 00, 30, 0, time.UTC),
		}, "nil"},
		{&session{ // insert valid session
			"0013",
			u2.id,
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
			-100,
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
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, result)
		}
	}
}

func TestSelectActiveSessionBy(t *testing.T) {
	u, _ := db.selectUserByEmail("batman@gotham.city")
	testCases := []struct {
		field    string
		id       interface{}
		expected string
	}{
		{"id", "0012", "nil"},               // select by sessionID
		{"user_id", u.id, "nil"},            // select by userID
		{"user_id", -100, "empty"},          // invalid userID
		{"email", u.email, "invalid field"}, //invalid field
	}
	for i, tc := range testCases {
		s, err := db.selectActiveSessionBy(tc.field, tc.id)

		result := "nil"
		if err != nil {
			result = err.Error()
		}
		if err == nil && s == nil {
			result = "empty"
		}
		if result != tc.expected {
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, result)
		} else {
			t.Log(s)
		}
	}
}

func TestUpdateSession(t *testing.T) {
	s, _ := db.selectActiveSessionBy("id", "0012")
	s.isActive = false

	err := db.updateSession(s)
	s, _ = db.selectActiveSessionBy("id", "0012")
	if err != nil {
		t.Errorf("Case: expected %v, got %v\n", "nil", err)
	} else {
		t.Log(s)
	}
}

func TestInsertPost(t *testing.T) {
	s, _ := db.selectActiveSessionBy("id", "0013")
	u, _ := db.selectUserByEmail("batman@gotham.city")
	testCases := []struct {
		p        post
		expected string
	}{
		{post{ // valid entry
			UserID:     s.userID,
			Title:      "Why is Batman so parnoid?",
			Content:    "He got way too much contigencies...",
			CreatedAt:  time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			categories: []int{1, 2},
		}, "nil"},
		{post{ // invalid userID
			UserID:     -100,
			Title:      "Why is Batman so serious?",
			Content:    "He can even take a joke...",
			CreatedAt:  time.Now(),
			categories: []int{1},
		}, "FOREIGN KEY constraint failed"},
		{post{ // invalid category
			UserID:     s.userID,
			Title:      "How many mothers are named Martha?",
			Content:    "Is it a common mother name?",
			CreatedAt:  time.Now(),
			categories: []int{10},
		}, "invalid category"},
		{post{ // valid entry
			UserID:     u.id,
			Title:      "Having 72 hours at day",
			Content:    "Here is how you train, invent, investigate and more...",
			CreatedAt:  time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
			categories: []int{1, 2, 4, 5},
		}, "nil"},
		{post{ // valid entry
			UserID:     u.id,
			Title:      "Why are there more and more supervillians?",
			Content:    "Is there a deep societal problems that creates supervillians?",
			CreatedAt:  time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
			categories: []int{4, 5},
		}, "nil"},
	}
	db.deleteAllPosts()

	for i, tc := range testCases {
		err := db.insertPost(tc.p)
		result := "nil"
		if err != nil {
			result = err.Error()
		}
		if result != tc.expected {
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, result)
		}
	}
}

func TestInsertFeedback(t *testing.T) {
	s, _ := db.selectActiveSessionBy("id", "0013")
	u, _ := db.selectUserByEmail("batman@gotham.city")
	posts, err := db.selectPosts("", "", 0)
	if err != nil {
		t.Error(err)
		return
	}
	testCases := []struct {
		tgt      string
		fb       feedback
		expected string
	}{
		{ // first like
			"post",
			feedback{
				userID:    s.userID,
				parentID:  (*posts)[1].ID,
				rating:    1,
				createdAt: time.Now(),
			}, "nil"},
		{ // second like by same user
			"post",
			feedback{
				userID:    s.userID,
				parentID:  (*posts)[1].ID,
				rating:    1,
				createdAt: time.Now(),
			}, "UNIQUE constraint failed: post_feedback.user_id, post_feedback.parent_id"},
		{ // second like by different user
			"post",
			feedback{
				userID:    u.id,
				parentID:  (*posts)[1].ID,
				rating:    1,
				createdAt: time.Now(),
			}, "nil"},
		{ // like a different post
			"post",
			feedback{
				userID:    u.id,
				parentID:  (*posts)[0].ID,
				rating:    1,
				createdAt: time.Now(),
			}, "nil"},
		{ // like a different post
			"post",
			feedback{
				userID:    s.userID,
				parentID:  (*posts)[2].ID,
				rating:    1,
				createdAt: time.Now(),
			}, "nil"},
		{ // invalid postID
			"post",
			feedback{
				userID:    s.userID,
				parentID:  -100,
				rating:    1,
				createdAt: time.Now(),
			}, "FOREIGN KEY constraint failed"},
	}
	for i, tc := range testCases {
		err := db.insertFeedback(tc.tgt, tc.fb)
		result := "nil"
		if err != nil {
			result = err.Error()
		}
		if result != tc.expected {
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, result)
		}
	}
}

func TestSelectUpdateFeedback(t *testing.T) {
	u, _ := db.selectUserByEmail("batman@gotham.city")
	// selectFeedback made in posts by user
	feedbacks, err := db.selectFeedback("post", u.id)
	if err != nil {
		t.Error(err)
		return
	}
	// change and update feedback on a post by user
	// this change should reflect in the results of TestSelectPosts
	(*feedbacks)[0].rating = 0
	err = db.updateFeedback("post", (*feedbacks)[0])
	if err != nil {
		t.Error(err)
	}
}

func TestSelectPosts(t *testing.T) {
	s, _ := db.selectActiveSessionBy("id", "0013")
	u, _ := db.selectUserByEmail("batman@gotham.city")
	testCase := []struct {
		filterBy string
		orderBy  string
		catID    int
		expected []time.Time
	}{
		{"", "", 0, []time.Time{ // default sorting
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
		}},
		{"", "oldest", 0, []time.Time{ // oldest post first
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
		}},
		{"createdBy", "", s.userID, []time.Time{ // filterBy superman
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
		}},
		{"createdBy", "", u.id, []time.Time{ // filterBy batman
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
		}},
		{"catergory", "", 1, []time.Time{ // filterBy batman
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
		}},
		{"catergory", "oldest", 4, []time.Time{ // filterBy batman
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
		}},
		{"catergory", "", 3, []time.Time{}}, // empty result
		{"likedBy", "likeCount", s.userID, []time.Time{ //likedBy superman
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
		}},
		{"likedBy", "likeCount", u.id, []time.Time{ //likedBy batman
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
		}},
	}

	for i, tc := range testCase {
		posts, err := db.selectPosts(tc.filterBy, tc.orderBy, tc.catID)
		if err != nil {
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, err)
			continue
		}
		result := true
		if len(tc.expected) != len(*posts) {
			result = false
		} else {
			for i, p := range *posts {
				if p.CreatedAt != tc.expected[i] {
					result = false
					break
				}
			}
		}
		if !result {
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, result)
		} else {
			t.Logf("Case %d: passed\n", i)
		}
		if posts != nil {
			for _, p := range *posts {
				t.Log(p)
			}
		}
	}
}
