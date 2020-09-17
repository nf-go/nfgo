package web

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"nfgo.ga/nfgo/nconst"
	"nfgo.ga/nfgo/nlog"

	"github.com/gin-gonic/gin"
	"nfgo.ga/nfgo/nerrors"
)

// Context -
type Context struct {
	*gin.Context
}

// APIResult -
type APIResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// Success - response and render the data by json
func (c *Context) Success(data interface{}) {
	r := &APIResult{Code: 0, Data: data}

	// logging the resp
	respLogger := nlog.Logger(c)
	if respLogger.IsLevelEnabled(nlog.DebugLevel) {
		if respJSON, err := json.Marshal(r); err == nil {
			respLogger.WithField("resp", string(respJSON)).Debug()
		}
	}

	c.JSON(http.StatusOK, r)
}

// Fail - response and render the error by json
func (c *Context) Fail(err error) {
	c.Error(err)

	// handle biz error
	if bizErr, ok := err.(nerrors.BizError); ok {
		var statusCode int
		switch bizErr {
		case nerrors.ErrForbidden:
			statusCode = http.StatusForbidden
		case nerrors.ErrUnauthorized:
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusOK
		}

		c.JSON(statusCode, &APIResult{
			Code: bizErr.Code(),
			Msg:  bizErr.Msg(),
		})
		nlog.Logger(c).WithError(err).Info()
		return
	}

	c.JSON(http.StatusInternalServerError, &APIResult{
		Code: nerrors.ErrInternal.Code(),
		Msg:  nerrors.ErrInternal.Msg(),
	})
	nlog.Logger(c).WithError(err).Error()
}

// FormFileBytes - returns the first file bytes for the provided form key.
func (c *Context) FormFileBytes(name string) ([]byte, string, error) {
	file, err := c.FormFile(name)
	if err != nil {
		return nil, "", err
	}
	filename := filepath.Base(file.Filename)

	src, err := file.Open()
	if err != nil {
		return nil, filename, err
	}
	defer src.Close()

	bytes, err := ioutil.ReadAll(src)
	return bytes, filename, err
}

// IsMultipartReq -
func (c *Context) IsMultipartReq() bool {
	if c.Request.Method != http.MethodPost {
		return false
	}
	contentType := c.ContentType()
	return strings.HasPrefix(contentType, "multipart/")
}

// ClientIP -
func (c *Context) ClientIP() string {
	remoteAddr := c.GetHeader(nconst.HeaderRealIP)
	if existsRemoteAddr(remoteAddr) {
		return remoteAddr
	}

	remoteAddr = c.GetHeader(nconst.HeaderForwardedFor)
	if remoteAddr != "" {
		ips := strings.Split(remoteAddr, ",")
		for _, ip := range ips {
			if existsRemoteAddr(ip) {
				return ip
			}
		}
	}

	remoteAddr = c.GetHeader("Proxy-Client-IP")
	if existsRemoteAddr(remoteAddr) {
		return remoteAddr
	}

	remoteAddr = c.GetHeader("WL-Proxy-Client-IP")
	if existsRemoteAddr(remoteAddr) {
		return remoteAddr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}

	return c.Request.RemoteAddr
}

func existsRemoteAddr(remoteAddr string) bool {
	if remoteAddr != "" {
		lowerRemoteAddr := strings.ToLower(remoteAddr)
		return lowerRemoteAddr != "unknown" && lowerRemoteAddr != "null"
	}
	return false
}

/************************************/
/***** context.Context *****/
/************************************/

// Deadline -
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Request.Context().Deadline()
}

// Done -
func (c *Context) Done() <-chan struct{} {
	return c.Request.Context().Done()
}

// Err -
func (c *Context) Err() error {
	return c.Request.Context().Err()
}

// Value -
func (c *Context) Value(key interface{}) interface{} {
	return c.Request.Context().Value(key)
}
