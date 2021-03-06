package resource

import (
	"fmt"

	"github.com/aelsabbahy/goss/system"
	"github.com/aelsabbahy/goss/util"
)

type User struct {
	Title    string  `json:"title,omitempty" yaml:"title,omitempty"`
	Meta     meta    `json:"meta,omitempty" yaml:"meta,omitempty"`
	Username string  `json:"-" yaml:"-"`
	Exists   bool    `json:"exists" yaml:"exists"`
	UID      matcher `json:"uid,omitempty" yaml:"uid,omitempty"`
	GID      matcher `json:"gid,omitempty" yaml:"gid,omitempty"`
	Groups   matcher `json:"groups,omitempty" yaml:"groups,omitempty"`
	Home     matcher `json:"home,omitempty" yaml:"home,omitempty"`
	Shell    matcher `json:"shell,omitempty" yaml:"shell,omitempty"`
}

func (u *User) ID() string      { return u.Username }
func (u *User) SetID(id string) { u.Username = id }

func (u *User) GetTitle() string { return u.Title }
func (u *User) GetMeta() meta    { return u.Meta }

func (u *User) Validate(sys *system.System) []TestResult {
	sysuser := sys.NewUser(u.Username, sys, util.Config{})

	var results []TestResult

	results = append(results, ValidateValue(u, "exists", u.Exists, sysuser.Exists))

	if u.UID != nil {
		uUID := deprecateAtoI(u.UID, fmt.Sprintf("%s: user.uid", u.Username))
		results = append(results, ValidateValue(u, "uid", uUID, sysuser.UID))
	}
	if u.GID != nil {
		uGID := deprecateAtoI(u.GID, fmt.Sprintf("%s: user.gid", u.Username))
		results = append(results, ValidateValue(u, "gid", uGID, sysuser.GID))
	}
	if u.Home != nil {
		results = append(results, ValidateValue(u, "home", u.Home, sysuser.Home))
	}
	if u.Groups != nil {
		results = append(results, ValidateValue(u, "groups", u.Groups, sysuser.Groups))
	}
	if u.Shell != nil {
		results = append(results, ValidateValue(u, "shell", u.Shell, sysuser.Shell))
	}

	return results
}

func NewUser(sysUser system.User, config util.Config) (*User, error) {
	username := sysUser.Username()
	exists, _ := sysUser.Exists()
	u := &User{
		Username: username,
		Exists:   exists,
	}
	if !contains(config.IgnoreList, "uid") {
		if uid, err := sysUser.UID(); err == nil {
			u.UID = uid
		}
	}
	if !contains(config.IgnoreList, "gid") {
		if gid, err := sysUser.GID(); err == nil {
			u.GID = gid
		}
	}
	if !contains(config.IgnoreList, "groups") {
		if groups, err := sysUser.Groups(); err == nil {
			u.Groups = groups
		}
	}
	if !contains(config.IgnoreList, "home") {
		if home, err := sysUser.Home(); err == nil {
			u.Home = home
		}
	}
	if !contains(config.IgnoreList, "shell") {
		if shell, err := sysUser.Shell(); err == nil {
			u.Shell = shell
		}
	}
	return u, nil
}
