package storage

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/httpclient"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// WebDAVObjectStorage represents WebDAV object storage
type WebDAVObjectStorage struct {
	httpClient   *http.Client
	webDavConfig *settings.WebDAVConfig
	rootPath     string
}

// NewWebDAVObjectStorage returns a WebDAV object storage
func NewWebDAVObjectStorage(config *settings.Config, pathPrefix string) (*WebDAVObjectStorage, error) {
	webDavConfig := config.WebDAVConfig

	storage := &WebDAVObjectStorage{
		httpClient:   httpclient.NewHttpClient(webDavConfig.RequestTimeout, webDavConfig.Proxy, webDavConfig.SkipTLSVerify, core.GetOutgoingUserAgent(), false),
		webDavConfig: webDavConfig,
		rootPath:     webDavConfig.RootPath,
	}

	storage.rootPath = storage.getFinalPath(pathPrefix)
	storage.rootPath = strings.ReplaceAll(storage.rootPath, "\\", "/")

	ctx := core.NewNullContext()
	exists, err := storage.directoryExists(ctx, storage.rootPath)

	if err != nil {
		return nil, err
	}

	if !exists {
		err := storage.createAllDirectories(ctx, "", storage.rootPath)

		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}

// Exists returns whether the file exists
func (s *WebDAVObjectStorage) Exists(ctx core.Context, path string) (bool, error) {
	req, err := http.NewRequest("HEAD", s.getFinalFileUrl(path), nil)

	if err != nil {
		return false, err
	}

	req.SetBasicAuth(s.webDavConfig.Username, s.webDavConfig.Password)
	resp, err := s.httpClient.Do(req)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.Exists] cannot check file exists, because %s", err.Error())
		return false, err
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	log.Errorf(ctx, "[webdav_storage.Exists] cannot check file exists, http status code is %d", resp.StatusCode)
	return false, errs.ErrSystemError
}

// Read returns the object instance according to specified the file path
func (s *WebDAVObjectStorage) Read(ctx core.Context, path string) (ObjectInStorage, error) {
	req, err := http.NewRequest("GET", s.getFinalFileUrl(path), nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.webDavConfig.Username, s.webDavConfig.Password)
	resp, err := s.httpClient.Do(req)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.Read] cannot get file, because %s", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.Read] cannot read response (http status code %d) body, because %s", resp.StatusCode, err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf(ctx, "[webdav_storage.Read] cannot get file, http status code is %d, response is %s", resp.StatusCode, string(body))
		return nil, errs.ErrSystemError
	}

	return newByteSliceObject(body), nil
}

// Save returns whether save the object instance successfully
func (s *WebDAVObjectStorage) Save(ctx core.Context, path string, object ObjectInStorage) error {
	finalPath := s.getFinalPath(path)
	dir := strings.ReplaceAll(filepath.Dir(finalPath), "\\", "/")

	exists, err := s.directoryExists(ctx, dir)

	if err != nil {
		return err
	}

	if !exists {
		rootExists, err := s.directoryExists(ctx, s.rootPath)

		if err != nil {
			return err
		}

		if !rootExists {
			err := s.createAllDirectories(ctx, "", s.rootPath)

			if err != nil {
				return err
			}
		}

		err = s.createAllDirectories(ctx, s.rootPath, strings.ReplaceAll(filepath.Dir(path), "\\", "/"))

		if err != nil {
			return err
		}
	}

	data, err := io.ReadAll(object)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", s.getFinalFileUrl(path), bytes.NewReader(data))

	if err != nil {
		return err
	}

	req.SetBasicAuth(s.webDavConfig.Username, s.webDavConfig.Password)
	resp, err := s.httpClient.Do(req)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.Save] cannot save file, because %s", err.Error())
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.Save] cannot read response (http status code %d) body, because %s", resp.StatusCode, err.Error())
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.Errorf(ctx, "[webdav_storage.Save] cannot save file, http status code is %d, response is %s", resp.StatusCode, string(body))
		return errs.ErrSystemError
	}

	return nil
}

// Delete returns whether delete the object according to specified the file path successfully
func (s *WebDAVObjectStorage) Delete(ctx core.Context, path string) error {
	req, err := http.NewRequest("DELETE", s.getFinalFileUrl(path), nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(s.webDavConfig.Username, s.webDavConfig.Password)
	resp, err := s.httpClient.Do(req)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.Delete] cannot delete file, because %s", err.Error())
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.Delete] cannot read response (http status code %d) body, because %s", resp.StatusCode, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		log.Errorf(ctx, "[webdav_storage.Delete] cannot delete file, http status code is %d, response is %s", resp.StatusCode, string(body))
		return errs.ErrSystemError
	}

	return nil
}

func (s *WebDAVObjectStorage) directoryExists(ctx core.Context, path string) (bool, error) {
	req, err := http.NewRequest("PROPFIND", s.getFinalDirectoryUrl(path), nil)

	if err != nil {
		return false, err
	}

	req.SetBasicAuth(s.webDavConfig.Username, s.webDavConfig.Password)
	resp, err := s.httpClient.Do(req)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.directoryExists] cannot check directory exists, because %s", err.Error())
		return false, err
	}

	if resp.StatusCode == http.StatusMultiStatus || resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	log.Errorf(ctx, "[webdav_storage.directoryExists] cannot check directory exists, http status code is %d", resp.StatusCode)
	return false, errs.ErrSystemError
}

func (s *WebDAVObjectStorage) createDirectory(ctx core.Context, path string) error {
	req, err := http.NewRequest("MKCOL", s.getFinalDirectoryUrl(path), nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(s.webDavConfig.Username, s.webDavConfig.Password)
	resp, err := s.httpClient.Do(req)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.createDirectory] cannot create directory, because %s", err.Error())
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(ctx, "[webdav_storage.createDirectory] cannot read response (http status code %d) body, because %s", resp.StatusCode, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMethodNotAllowed {
		log.Errorf(ctx, "[webdav_storage.createDirectory] cannot create directory, http status code is %d, response is %s", resp.StatusCode, string(body))
		return errs.ErrSystemError
	}

	return nil
}

func (s *WebDAVObjectStorage) createAllDirectories(ctx core.Context, currentPath string, path string) error {
	directories := strings.Split(path, "/")

	for _, dir := range directories {
		if len(dir) == 0 {
			continue
		}

		currentPath = currentPath + "/" + dir
		exists, err := s.directoryExists(ctx, currentPath)

		if err != nil {
			return err
		}

		if !exists {
			err = s.createDirectory(ctx, currentPath)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *WebDAVObjectStorage) getFinalFileUrl(filePath string) string {
	finalUrl := s.webDavConfig.Url

	if len(finalUrl) < 1 || finalUrl[len(finalUrl)-1] != '/' {
		finalUrl = finalUrl + "/"
	}

	finalPath := s.getFinalPath(filePath)

	if len(finalPath) > 0 && finalPath[0] == '/' {
		finalPath = finalPath[1:]
	}

	return finalUrl + finalPath
}

func (s *WebDAVObjectStorage) getFinalDirectoryUrl(dirPath string) string {
	finalUrl := s.webDavConfig.Url

	if len(finalUrl) < 1 || finalUrl[len(finalUrl)-1] != '/' {
		finalUrl = finalUrl + "/"
	}

	if len(dirPath) > 0 && dirPath[0] == '/' {
		dirPath = dirPath[1:]
	}

	if len(dirPath) > 0 && dirPath[len(dirPath)-1] != '/' {
		dirPath = dirPath + "/"
	}

	return finalUrl + dirPath
}

func (s *WebDAVObjectStorage) getFinalPath(path string) string {
	rootPath := s.rootPath

	if len(rootPath) < 1 || rootPath[len(rootPath)-1] != '/' {
		rootPath = rootPath + "/"
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	path = strings.ReplaceAll(path, "\\", "/")

	return rootPath + path
}
