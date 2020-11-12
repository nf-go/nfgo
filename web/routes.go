// Copyright 2020 The nfgo Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import "net/http"

// Routes defines all router handle interface.
type Routes interface {
	Use(handlers ...HandlerFunc)

	Handle(httpMethod, relativePath string, handlers ...HandlerFunc)

	Any(relativePath string, handlers ...HandlerFunc)

	GET(relativePath string, handlers ...HandlerFunc)

	POST(relativePath string, handlers ...HandlerFunc)

	DELETE(relativePath string, handlers ...HandlerFunc)

	PATCH(relativePath string, handlers ...HandlerFunc)

	PUT(relativePath string, handlers ...HandlerFunc)

	OPTIONS(relativePath string, handlers ...HandlerFunc)

	HEAD(relativePath string, handlers ...HandlerFunc)

	StaticFile(relativePath string, filepath string)

	Static(relativePath string, root string)

	StaticFS(relativePath string, fs http.FileSystem)
}
