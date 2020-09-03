// See LICENSE.txt for license information.

package einterfaces

import "github.com/mattermost/mattermost-server/v5/model"

type GenericSocketExporter interface {
	Export(*model.Config, *model.WebSocketEvent) bool
	InitExporter(*model.Config)
}
