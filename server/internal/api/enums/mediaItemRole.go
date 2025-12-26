package enums

type MediaItemRole string

const (
	MediaItemRoleOriginal  MediaItemRole = "original"
	MediaItemRoleThumbnail MediaItemRole = "thumbnail"
	MediaItemRoleViewer    MediaItemRole = "viewer"
	MediaItemRolePreview   MediaItemRole = "preview"
	MediaItemRoleLive      MediaItemRole = "live"
	MediaItemRoleStream    MediaItemRole = "stream"
)
