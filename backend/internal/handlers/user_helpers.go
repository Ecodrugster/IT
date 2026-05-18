package handlers

import (
	"strings"
)

func normalizeRole(raw string) string {
	role := strings.ToLower(strings.TrimSpace(raw))
	switch role {
	case "admin", "teacher", "student":
		return role
	default:
		return "student"
	}
}

func normalizeUserDocument(data map[string]interface{}, uid string) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range data {
		out[k] = v
	}

	out["uid"] = uid

	displayName := getDisplayNameFromUserData(data, uid)
	out["display_name"] = displayName
	out["displayName"] = displayName

	email := asString(data["email"])
	if email != "" {
		out["email"] = email
	}

	photoURL := asString(data["photo_url"])
	if photoURL == "" {
		photoURL = asString(data["photoURL"])
	}
	out["photo_url"] = photoURL
	out["photoURL"] = photoURL

	role := normalizeRole(asString(data["role"]))
	out["role"] = role

	groupName := asString(data["group_name"])
	if groupName == "" {
		groupName = asString(data["group"])
	}
	out["group_name"] = groupName
	out["group"] = groupName

	return out
}
