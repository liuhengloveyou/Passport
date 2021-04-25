package accessctl

import (
	casbin "github.com/casbin/casbin/v2"
	_ "github.com/go-sql-driver/mysql"
)

var (
	enforcer *casbin.SyncedEnforcer
)

func enforce(sub, domain, obj, act string) (bool, error) {
	return enforcer.Enforce(sub, domain, obj, act)
}

func addPolicy(sub, domain, obj, act string) (err error) {
	if _, err = enforcer.AddPolicy(sub, domain, obj, act); err != nil {
		return
	}

	if err = enforcer.SavePolicy(); err != nil {
		return
	}

	return
}

func removePolicy(sub, domain, obj, act string) (err error) {
	if _, err = enforcer.RemovePolicy(sub, domain, obj, act); err != nil {
		return
	}

	if err = enforcer.SavePolicy(); err != nil {
		return
	}

	return
}

func getFilteredPolicy(domain string) [][]string {
	return enforcer.GetFilteredPolicy(1, domain)
}

func addRoleForUserInDomain(user, role, domain string) (err error) {
	if _, err = enforcer.AddRoleForUserInDomain(user, role, domain); err != nil {
		return
	}

	if err = enforcer.SavePolicy(); err != nil {
		return
	}

	return
}

func deleteRoleForUserInDomain(user, role, domain string) (err error) {
	if _, err = enforcer.DeleteRoleForUserInDomain(user, role, domain); err != nil {
		return
	}

	if err = enforcer.SavePolicy(); err != nil {
		return
	}

	return
}

func getRoleForUserInDomain(user, domain string) []string {
	return enforcer.GetRolesForUserInDomain(user, domain)
}

func getUsersForRoleInDomain(role, domain string) []string {
	return enforcer.GetUsersForRoleInDomain(role, domain)
}
