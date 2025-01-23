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
			TypeID:    3,
			Name:      "BATMAN",
			Email:     "batman@gotham.city",
			PwHash:    "bruc3W47N3",
			RegDate:   fixedTime,
			LastLogin: fixedTime},
			"nil"},
		{&user{ // duplicate user Email
			TypeID:    3,
			Name:      "BATMAN",
			Email:     "batman@gotham.city",
			PwHash:    "dickGrayson",
			RegDate:   fixedTime,
			LastLogin: fixedTime},
			"UNIQUE constraint failed: users.email"},
		{&user{ // invalid TypeID
			TypeID:    666,
			Name:      "JOKER",
			Email:     "j0k3r@gotham.city",
			PwHash:    "alfredPennyless",
			RegDate:   fixedTime,
			LastLogin: fixedTime},
			"FOREIGN KEY constraint failed"},
		{&user{
			TypeID:    1,
			Name:      "Superman",
			Email:     "superman@metropolis.city",
			PwHash:    "clarkKent",
			RegDate:   fixedTime,
			LastLogin: fixedTime},
			"nil"},
	}
	db.DeleteAllUsers()
	for i, tc := range testCases {
		err := db.InsertUser(tc.u)
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
		Email    string
		Name     string
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
		u, err := db.SelectUserByEmail(tc.Email)
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
		{"Name", // check categories
			"categories",
			[]string{"General", "golang", "html", "css", "sqlite3"}},
		{"Email", // check user Emails
			"users",
			[]string{"batman@gotham.city", "superman@metropolis.city"}},
		{"title", // empty result
			"posts",
			[]string{}},
	}
	for i, tc := range testCases {
		results, err := db.SelectFieldFromTable(tc.field, tc.table)
		if !isEqualStringSlice(results, tc.expected) && err == nil {
			t.Errorf("case%v: failed. %v\n", i, results)
		} else {
			t.Log(results)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	u, _ := db.SelectUserByEmail("superman@metropolis.city")
	fixedTime := time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC)
	u.Name = "Ultimate Superman"
	u.PwHash = "kalEl"
	u.LastLogin = fixedTime

	err := db.UpdateUser(u)
	if err != nil {
		t.Errorf("Case: updateUser() failed.%s\n", err)
		return
	}
	u2, err2 := db.SelectUserByEmail(u.Email)
	if (err2 != nil || u2 == nil) &&
		u.Name != u2.Name &&
		u2.LastLogin == fixedTime {
		t.Errorf("Case: failed.\n")
	} else {
		t.Log(u2)
	}
}

func TestInsertSession(t *testing.T) {
	u, _ := db.SelectUserByEmail("batman@gotham.city")
	u2, _ := db.SelectUserByEmail("superman@metropolis.city")
	testCases := []struct {
		s        *session
		expected string
	}{
		{&session{ // insert valid session
			ID:         "0012",
			UserID:     u.ID,
			IsActive:   true,
			StartTime:  time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			ExpireTime: time.Date(2025, 1, 17, 13, 11, 59, 0, time.UTC),
			LastAccess: time.Date(2025, 1, 17, 13, 00, 30, 0, time.UTC),
		}, "nil"},
		{&session{ // insert valid session
			ID:         "0013",
			UserID:     u2.ID,
			IsActive:   true,
			StartTime:  time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			ExpireTime: time.Date(2025, 1, 17, 13, 11, 59, 0, time.UTC),
			LastAccess: time.Date(2025, 1, 17, 13, 00, 30, 0, time.UTC),
		}, "nil"},
		{&session{ // duplicated session id
			ID:         "0012",
			UserID:     1,
			IsActive:   true,
			StartTime:  time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			ExpireTime: time.Date(2025, 1, 17, 13, 11, 59, 0, time.UTC),
			LastAccess: time.Date(2025, 1, 17, 13, 00, 30, 0, time.UTC),
		}, "UNIQUE constraint failed: sessions.id"},
		{&session{ // invlaid user id
			ID:         "0015",
			UserID:     -100,
			IsActive:   true,
			StartTime:  time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			ExpireTime: time.Date(2025, 1, 17, 13, 11, 59, 0, time.UTC),
			LastAccess: time.Date(2025, 1, 17, 13, 00, 30, 0, time.UTC),
		}, "FOREIGN KEY constraint failed"},
	}
	db.DeleteAllSessions()
	for i, tc := range testCases {
		err := db.InsertSession(tc.s)
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
	u, _ := db.SelectUserByEmail("batman@gotham.city")
	testCases := []struct {
		field    string
		id       interface{}
		expected string
	}{
		{"id", "0012", "nil"},               // select by sessionID
		{"user_id", u.ID, "nil"},            // select by userID
		{"user_id", -100, "empty"},          // invalid userID
		{"Email", u.Email, "invalid field"}, // invalid field
	}
	for i, tc := range testCases {
		s, err := db.SelectActiveSessionBy(tc.field, tc.id)

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
	s, _ := db.SelectActiveSessionBy("id", "0012")
	s.IsActive = false

	err := db.UpdateSession(s)
	s, _ = db.SelectActiveSessionBy("id", "0012")
	if err != nil {
		t.Errorf("Case: expected %v, got %v\n", "nil", err)
	} else {
		t.Log(s)
	}
}

func TestInsertPost(t *testing.T) {
	s, _ := db.SelectActiveSessionBy("id", "0013")
	u, _ := db.SelectUserByEmail("batman@gotham.city")
	testCases := []struct {
		p        post
		expected string
	}{
		{post{ // valid entry
			UserID:     s.UserID,
			Title:      "Why is Batman so parnoid?",
			Content:    "He got way too much contigencies...",
			CreatedAt:  time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			Categories: []int{0, 1},
		}, "nil"},
		{post{ // invalid userID
			UserID:     -100,
			Title:      "Why is Batman so serious?",
			Content:    "He can even take a joke...",
			CreatedAt:  time.Now(),
			Categories: []int{0},
		}, "FOREIGN KEY constraint failed"},
		{post{ // invalid category
			UserID:     s.UserID,
			Title:      "How many mothers are Named Martha?",
			Content:    "Is it a common mother Name?",
			CreatedAt:  time.Now(),
			Categories: []int{10},
		}, "invalid category"},
		{post{ // valid entry
			UserID:     u.ID,
			Title:      "Having 72 hours at day",
			Content:    "Here is how you train, invent, investigate and more...",
			CreatedAt:  time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
			Categories: []int{0, 1, 3, 4},
		}, "nil"},
		{post{ // valid entry
			UserID:     u.ID,
			Title:      "Why are there more and more supervillians?",
			Content:    "Is there a deep societal problems that creates supervillians?",
			CreatedAt:  time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
			Categories: []int{3, 4},
		}, "nil"},
	}
	db.DeleteAllPosts()

	for i, tc := range testCases {
		err := db.InsertPost(tc.p)
		result := "nil"
		if err != nil {
			result = err.Error()
		}
		if result != tc.expected {
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, result)
		}
	}
}

func TestInsertComment(t *testing.T) {
	s, _ := db.SelectActiveSessionBy("id", "0013")
	u, _ := db.SelectUserByEmail("batman@gotham.city")
	posts, _ := db.SelectPosts("", "oldest", 0)

	testCases := []struct {
		c        comment
		expected string
	}{
		{ // valid comment in first post
			comment{
				UserID:    s.UserID,
				PostID:    posts[0].ID,
				Content:   "Why can't be more trustful of us?",
				CreatedAt: time.Now(),
			}, "nil",
		},
		{ // valid comment in first post
			comment{
				UserID:    u.ID,
				PostID:    posts[0].ID,
				Content:   "Not to be xenophobic... but you are not from around here",
				CreatedAt: time.Now(),
			}, "nil",
		},
		{ // invalid comment to first post
			comment{
				UserID:    -100,
				PostID:    posts[0].ID,
				Content:   "he just so serious all the time",
				CreatedAt: time.Now(),
			}, "FOREIGN KEY constraint failed",
		},
		{ // valid comment in third post
			comment{
				UserID:    s.UserID,
				PostID:    posts[2].ID,
				Content:   "You can always move faster. Oh! You can't travel at the speed of light.",
				CreatedAt: time.Now(),
			}, "nil",
		},
		{ // valid comment in third post
			comment{
				UserID:    u.ID,
				PostID:    posts[2].ID,
				Content:   "At least I don't get defeated by some rocks.",
				CreatedAt: time.Now(),
			}, "nil",
		},
		{ // valid comment in third post
			comment{
				UserID:    s.UserID,
				PostID:    posts[2].ID,
				Content:   "Humans get affected by radioactive materials too.",
				CreatedAt: time.Now(),
			}, "nil",
		},
		{ // valid comment in third post
			comment{
				UserID:    s.UserID,
				PostID:    posts[2].ID,
				Content:   "So we both bleed...",
				CreatedAt: time.Now(),
			}, "nil",
		},
	}
	db.DeleteAllComments()

	for i, tc := range testCases {
		err := db.InsertComment(tc.c)
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
	s, _ := db.SelectActiveSessionBy("id", "0013")
	u, _ := db.SelectUserByEmail("batman@gotham.city")
	posts, err := db.SelectPosts("", "", 0)
	if err != nil {
		t.Error(err)
		return
	}
	comments, err := db.SelectComments(posts[0].ID, "")
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
				UserID:    s.UserID,
				ParentID:  posts[1].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "nil"},
		{ // second like by same user
			"post",
			feedback{
				UserID:    s.UserID,
				ParentID:  posts[1].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "UNIQUE constraint failed: post_feedback.user_id, post_feedback.parent_id"},
		{ // second like by different user
			"post",
			feedback{
				UserID:    u.ID,
				ParentID:  posts[1].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "nil"},
		{ // like a different post
			"post",
			feedback{
				UserID:    u.ID,
				ParentID:  posts[0].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "nil"},
		{ // like a different post
			"post",
			feedback{
				UserID:    s.UserID,
				ParentID:  posts[2].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "nil"},
		{ // invalid postID
			"post",
			feedback{
				UserID:    s.UserID,
				ParentID:  -100,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "FOREIGN KEY constraint failed"},
		{ // like invalid comment
			"comment",
			feedback{
				UserID:    s.UserID,
				ParentID:  -100,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "FOREIGN KEY constraint failed"},
		{ // like comment
			"comment",
			feedback{
				UserID:    s.UserID,
				ParentID:  comments[0].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "nil"},
		{ // like comment again
			"comment",
			feedback{
				UserID:    s.UserID,
				ParentID:  comments[0].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "UNIQUE constraint failed: comment_feedback.user_id, comment_feedback.parent_id"},
		{ // like comment by different user
			"comment",
			feedback{
				UserID:    u.ID,
				ParentID:  comments[0].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "nil"},
		{ // like a different comment
			"comment",
			feedback{
				UserID:    u.ID,
				ParentID:  comments[1].ID,
				Rating:    1,
				CreatedAt: time.Now(),
			}, "nil"},
	}
	for i, tc := range testCases {
		err := db.InsertFeedback(tc.tgt, tc.fb)
		result := "nil"
		if err != nil {
			result = err.Error()
		}
		if result != tc.expected {
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, result)
		}
	}
}

func TestSelectUpdateFeedbacks(t *testing.T) {
	u, _ := db.SelectUserByEmail("batman@gotham.city")
	// selectFeedback made in posts by user
	feedbacks, err := db.SelectFeedbacks("post", u.ID)
	if err != nil {
		t.Error(err)
		return
	}
	// change and update feedback on a post by user
	// this change should reflect in the results of TestSelectPosts
	feedbacks[0].Rating = 0
	err = db.UpdateFeedback("post", *feedbacks[0])
	if err != nil {
		t.Error(err)
	}
}

func TestSelectPosts(t *testing.T) {
	s, _ := db.SelectActiveSessionBy("id", "0013")
	u, _ := db.SelectUserByEmail("batman@gotham.city")
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
		{"createdBy", "", s.UserID, []time.Time{ // filterBy superman
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
		}},
		{"createdBy", "", u.ID, []time.Time{ // filterBy batman
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
		}},
		{"catergory", "", 0, []time.Time{ // filterBy batman
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
		}},
		{"catergory", "oldest", 4, []time.Time{ // filterBy batman
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
		}},
		{"catergory", "", 2, []time.Time{}}, // empty result
		{"likedBy", "likeCount", s.UserID, []time.Time{ //likedBy superman
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
		}},
		{"likedBy", "likeCount", u.ID, []time.Time{ //likedBy batman
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
		}},
		{"", "commentCount", -1, []time.Time{ //sortBy comment_count
			time.Date(2025, 1, 17, 13, 12, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 11, 59, 0, time.UTC),
			time.Date(2025, 1, 17, 12, 12, 0, 0, time.UTC),
		}},
	}

	for i, tc := range testCase {
		posts, err := db.SelectPosts(tc.filterBy, tc.orderBy, tc.catID)
		if err != nil {
			t.Errorf("Case %d: expected %v, got %v\n", i, tc.expected, err)
			continue
		}
		result := true
		if len(tc.expected) != len(posts) {
			result = false
		} else {
			for i, p := range posts {
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
		for _, p := range posts {
			t.Log(p)
		}
	}
}
