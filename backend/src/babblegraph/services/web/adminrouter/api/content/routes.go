package content

import (
	"babblegraph/model/admin"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
)

var Routes = router.RouteGroup{
	Prefix: "content",
	Routes: []router.Route{
		{
			Path: "get_all_topics_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				getAllContentTopics,
			),
		}, {
			Path: "add_topic_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				addContentTopic,
			),
		}, {
			Path: "get_topic_by_id_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				getTopicByID,
			),
		}, {
			Path: "update_is_topic_active_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				updateIsContentTopicActive,
			),
		}, {
			Path: "get_all_topic_display_names_for_topic_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				getAllTopicDisplayNamesForTopic,
			),
		}, {
			Path: "add_topic_display_name_for_topic_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				addTopicDisplayNameForTopic,
			),
		}, {
			Path: "update_topic_display_name_label_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				updateTopicDisplayNameLabel,
			),
		}, {
			Path: "toggle_topic_display_name_is_active_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				toggleTopicDisplayNameIsActive,
			),
		}, {
			Path: "get_all_sources_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				getAllSources,
			),
		}, {
			Path: "get_source_by_id_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				getSourceByID,
			),
		}, {
			Path: "add_source_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				addSource,
			),
		}, {
			Path: "update_source_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				updateSource,
			),
		}, {
			Path: "get_all_source_seeds_for_source_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				getAllSourceSeedsForSource,
			),
		}, {
			Path: "add_source_seed_for_source_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				addSourceSeed,
			),
		}, {
			Path: "update_source_seed_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				updateSourceSeed,
			),
		}, {
			Path: "upsert_source_seed_mappings_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				upsertSourceSeedMappings,
			),
		}, {
			Path: "get_source_seed_mappings_for_source_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				getSourceSourceSeedMappingsForSource,
			),
		},
	},
}
