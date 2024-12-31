package folder_mod

import (
	"goa.design/goa/v3/http"
	"mai.today/api/gen/folder"
	"mai.today/api/gen/http/folder/server"
)

func NewFolderServer() *server.Server {
	endpoints := folder.NewEndpoints(Instance())

	dec := http.RequestDecoder
	enc := http.ResponseEncoder
	svr := server.New(endpoints, http.NewMuxer(), dec, enc, nil, nil)

	return svr
}
