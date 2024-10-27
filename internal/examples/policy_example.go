package examples

import "github.com/el-mike/restrict/v2"

var ExamplePolicy = &restrict.PolicyDefinition{
	// A map of reusable Permissions. Key corresponds to a preset's name, which can
	// be later used to apply it.
	PermissionPresets: restrict.PermissionPresets{
		"updateOwn": &restrict.Permission{
			// An action that given Permission allows to perform.
			Action: "update",
			// Optional Conditions that when defined, need to be satisfied in order
			// to allow the access.
			Conditions: restrict.Conditions{
				// EqualCondition requires two values (described by ValueDescriptors)
				// to be equal in order to grant the access.
				// In this example we want to check if Conversation.CreatedBy and User.ID
				// are the same, meaning that Conversation was created by given User.
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
	},
	// A map of Roles. Key corresponds to a Role that any Subject in your system can belong to.
	Roles: restrict.Roles{
		"User": {
			// Optional, human readable description.
			Description: "This is a simple User role, with permissions for basic chat operations.",
			// Map of Permissions per Resource.
			// Grants map can be nil, meaning given Role has no Permissions (but can still inherit some).
			Grants: restrict.GrantsMap{
				"Conversation": {
					// Subject "User" can "read" any "Conversation".
					&restrict.Permission{Action: "read"},
					// Subject "User" can "create" a "Conversation".
					&restrict.Permission{Action: "create"},
					// Subject "User" can "update" ONLY a "Coversation" that was
					// created by it. Check "updateOwn" preset definition above.
					&restrict.Permission{Preset: "updateOwn"},
					// Subject "User" can "delete" ONLY inactive "Conversation".
					&restrict.Permission{
						Action: "delete",
						Conditions: restrict.Conditions{
							// EmptyCondition requires a value (described by ValueDescriptor)
							// to be empty (falsy) in order to grant the access.
							// In this example, we want Conversation.Active to be false.
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
		"Admin": {
			Description: "This is an Admin role, with permissions to manage Users.",
			// "Admin" can do everything "User" can.
			Parents: []string{"User"},
			// AND can also perform other operations that User itself
			// is not allowed to do.
			Grants: restrict.GrantsMap{
				// Please note that in order to make this work,
				// User needs to implement Resource interface.
				"User": {
					// Subject "Admin" can create a "User".
					&restrict.Permission{Action: "create"},
				},
			},
		},
	},
}
