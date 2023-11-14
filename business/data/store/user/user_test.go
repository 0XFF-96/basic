package user_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/yourusername/basic-a/business/core/user"
	userStore "github.com/yourusername/basic-a/business/data/store/user"
	"github.com/yourusername/basic-a/business/data/store/user/usercache"
	"github.com/yourusername/basic-a/business/data/tests"
	"github.com/yourusername/basic-a/foundation/docker"
	"runtime/debug"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = tests.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tests.StopDB(c)

	m.Run()
}

func Test_User(t *testing.T) {
	log, db, teardown := tests.NewUnit(t, c, "testuser")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		teardown()
	}()

	core := user.NewCore(usercache.NewStore(log, userStore.NewStore(log, db)))

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			ctx := context.Background()

			nu := user.NewUser{
				Name:            "Bill Kennedy",
				Email:           "bill@ardanlabs.com",
				Roles:           []string{user.RoleAdmin},
				Password:        "gophers",
				PasswordConfirm: "gophers",
			}

			usr, err := core.Create(ctx, nu)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			saved, err := core.QueryByID(ctx, usr.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", tests.Success, testID)

			if saved.DateCreated.Sub(usr.DateCreated) > time.Millisecond {
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.DateCreated)
				t.Logf("\t\tTest %d:\tExp: %v", testID, usr.DateCreated)
				t.Logf("\t\tTest %d:\tDiff: %v", testID, saved.DateCreated.Sub(usr.DateCreated))
				t.Fatalf("\t%s\tTest %d:\tShould get back the same date created.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same date created.", tests.Success, testID)

			if saved.DateUpdated.Sub(usr.DateUpdated) > time.Millisecond {
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.DateUpdated)
				t.Logf("\t\tTest %d:\tExp: %v", testID, usr.DateUpdated)
				t.Logf("\t\tTest %d:\tDiff: %v", testID, saved.DateUpdated.Sub(usr.DateUpdated))
				t.Fatalf("\t%s\tTest %d:\tShould get back the same date updated.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same date updated.", tests.Success, testID)

			usr.DateCreated = time.Time{}
			usr.DateUpdated = time.Time{}
			saved.DateCreated = time.Time{}
			saved.DateUpdated = time.Time{}

			if diff := cmp.Diff(usr, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", tests.Success, testID)

			upd := user.UpdateUser{
				Name:  tests.StringPointer("Jacob Walker"),
				Email: tests.StringPointer("jacob@ardanlabs.com"),
			}

			if err := core.Update(ctx, usr.ID, upd); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", tests.Success, testID)

			saved, err = core.QueryByEmail(ctx, *upd.Email)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by Email : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by Email.", tests.Success, testID)

			diff := usr.DateUpdated.Sub(saved.DateUpdated)
			if diff > 0 {
				t.Fatalf("Should have a larger DateUpdated : sav %v, usr %v, dif %v", saved.DateUpdated, usr.DateUpdated, diff)
			}

			if saved.Name != *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", tests.Success, testID)
			}

			if saved.Email != *upd.Email {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Email.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Email)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Email)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Email.", tests.Success, testID)
			}

			if err := core.Delete(ctx, usr.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", tests.Success, testID)

			_, err = core.QueryByID(ctx, usr.ID)
			if !errors.Is(err, user.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve user.", tests.Success, testID)
		}
	}
}

func Test_PagingUser(t *testing.T) {
	log, db, teardown := tests.NewUnit(t, c, "testpaging")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		teardown()
	}()

	core := user.NewCore(userStore.NewStore(log, db))

	t.Log("Given the need to page through User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 users.", testID)
		{
			ctx := context.Background()

			name := "User Gopher"
			users1, err := core.Query(ctx, user.QueryFilter{Name: &name}, user.DefaultOrderBy, 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user %q : %s.", tests.Failed, testID, name, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user %q.", tests.Success, testID, name)

			if len(users1) != 1 && users1[0].Name == name {
				t.Fatalf("\t%s\tTest %d:\tShould have a single user for %q : %s.", tests.Failed, testID, name, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single user.", tests.Success, testID)

			name = "Admin Gopher"
			users2, err := core.Query(ctx, user.QueryFilter{Name: &name}, user.DefaultOrderBy, 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user %q : %s.", tests.Failed, testID, name, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve users %q.", tests.Success, testID, name)

			if len(users2) != 1 && users2[0].Name == name {
				t.Fatalf("\t%s\tTest %d:\tShould have a single user for %q : %s.", tests.Failed, testID, name, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single user.", tests.Success, testID)

			users3, err := core.Query(ctx, user.QueryFilter{}, user.DefaultOrderBy, 1, 2)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve 2 users for page 1 : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve 2 users for page 1.", tests.Success, testID)

			if len(users3) != 2 {
				t.Logf("\t\tTest %d:\tgot: %v", testID, len(users3))
				t.Logf("\t\tTest %d:\texp: %v", testID, 2)
				t.Fatalf("\t%s\tTest %d:\tShould have 2 users for page 1 : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have 2 users for page 1.", tests.Success, testID)

			if users3[0].ID == users3[1].ID {
				t.Logf("\t\tTest %d:\tUser1: %v", testID, users3[0].ID)
				t.Logf("\t\tTest %d:\tUser2: %v", testID, users3[1].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different users : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different users.", tests.Success, testID)
		}
	}
}
