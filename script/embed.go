package script

import _ "embed"

//go:embed top10.pgsql
var Top10Script string

//go:embed archive_events.pgsql
var ArchiveOldEventsScript string

//go:embed update_status.pgsql
var UpdateStatusesScript string