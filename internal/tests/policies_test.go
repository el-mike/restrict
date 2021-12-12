package tests

import "github.com/el-Mike/restrict"

var PolicyOne = &restrict.PolicyDefinition{
	PermissionPresets: restrict.PermissionPresets{
		"readOwn": &restrict.Permission{
			Action: "read",
			Conditions: restrict.Conditions{
				&restrict.EqualCondition{
					ID: "isOwner",
					Left: &restrict.ValueDescriptor{
						Source: restrict.ResourceField,
						Field:  "CreatedBy",
					},
					Right: &restrict.ValueDescriptor{
						Source: restrict.SubjectField,
						Field:  "ID",
					},
				},
			},
		},
		"accessSelf": &restrict.Permission{
			Conditions: restrict.Conditions{
				&restrict.EqualCondition{
					ID: "self",
					Left: &restrict.ValueDescriptor{
						Source: restrict.ResourceField,
						Field:  "ID",
					},
					Right: &restrict.ValueDescriptor{
						Source: restrict.SubjectField,
						Field:  "ID",
					},
				},
			},
		},
	},
	Roles: restrict.Roles{
		BasicUserRole: {
			Grants: restrict.GrantsMap{
				UserResource: {
					&restrict.Permission{
						Action: "read",
						Preset: "accessSelf",
					},
				},
			},
		},
		UserRole: {
			Parents: []string{BasicUserRole},
			Grants: restrict.GrantsMap{
				ConversationResource: {
					&restrict.Permission{Preset: "readOwn"},
					&restrict.Permission{Action: "create"},
					&restrict.Permission{
						Action: "delete",
						Conditions: restrict.Conditions{
							&restrict.EmptyCondition{
								ID: "deleteActive",
								Value: &restrict.ValueDescriptor{
									Source: restrict.ResourceField,
									Field:  "Active",
								},
							},
						},
					},
				},
			},
		},
		AdminRole: {
			Parents: []string{UserRole},
			Grants: restrict.GrantsMap{
				ConversationResource: {
					&restrict.Permission{Action: "read"},
				},
				UserResource: {
					&restrict.Permission{Action: "create"},
					&restrict.Permission{Action: "read"},
					&restrict.Permission{Action: "update"},
					&restrict.Permission{Action: "delete"},
				},
			},
		},
	},
}
