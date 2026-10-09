package extservices

import (
	"strconv"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// Person is somebody who works in a business, for showing who did what
type Person struct {
	Uid    int64
	Name   string
	Role   extmodels.Role
	Active bool // false for people who were removed or left
}

// DisplayName returns the nickname of a user, or the user name when there is none
func DisplayName(c core.Context, uid int64) string {
	user, err := services.Users.GetUserById(c, uid)

	if err != nil || user == nil {
		return ""
	}

	if name := strings.TrimSpace(user.Nickname); name != "" {
		return name
	}

	return user.Username
}

// People lists the owner and everyone who has worked in the business, owner first
func People(c core.Context, ownerUid int64) ([]*Person, error) {
	people := []*Person{{Uid: ownerUid, Name: DisplayName(c, ownerUid), Role: extmodels.RoleOwner, Active: true}}

	memberships, err := Memberships.ListAllByOwner(c, ownerUid)

	if err != nil {
		return nil, err
	}

	for _, m := range memberships {
		people = append(people, &Person{Uid: m.StaffUid, Name: DisplayName(c, m.StaffUid), Role: m.Role, Active: m.Status == extmodels.MembershipStatusActive})
	}

	return people, nil
}

// StampText marks a transaction comment with the person who recorded it, so it shows everywhere the comment does,
// including the app's own transaction list and exports. No name means no mark. The result never exceeds maxLength
// characters: the original text is shortened, never the mark.
func StampText(base string, name string, maxLength int) string {
	name = strings.TrimSpace(name)

	if name == "" {
		return base
	}

	mark := "by " + name
	base = strings.TrimSpace(base)

	if base != "" {
		mark = " · " + mark
	}

	room := maxLength - len([]rune(mark))

	if room < 0 {
		return string([]rune(mark)[:maxLength])
	}

	runes := []rune(base)

	if len(runes) > room {
		runes = runes[:room]
	}

	return string(runes) + mark
}

// stamped returns the transaction comment for something recorded by actorUid in ownerUid's books.
// The owner's own entries are left unmarked: no mark means the owner did it.
func stamped(c core.Context, ownerUid int64, actorUid int64, comment string) string {
	if actorUid <= 0 || actorUid == ownerUid {
		return truncate(comment, 255)
	}

	return StampText(comment, DisplayName(c, actorUid), 255)
}

func fmtInt(n int64) string {
	return strconv.FormatInt(n, 10)
}
