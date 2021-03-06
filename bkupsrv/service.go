package bkupsrv

/*
 Six910 is a shopping cart and E-commerce system.
 Copyright (C) 2020 Ulbora Labs LLC. (www.ulboralabs.com)
 All rights reserved.
 Copyright (C) 2020 Ken Williamson
 All rights reserved.
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	lg "github.com/Ulbora/Level_Logger"
	ds "github.com/Ulbora/json-datastore"
)

//BackupService BackupService
type BackupService interface {
	UploadBackups(bk *[]byte) bool
	DownloadBackups() (bool, *[]byte)
}

//Six910BackupService Six910BackupService
type Six910BackupService struct {
	Store             ds.JSONDatastore
	TemplateStore     ds.JSONDatastore
	CarouselStore     ds.JSONDatastore
	CountryStore      ds.JSONDatastore
	CSSStore          ds.JSONDatastore
	MenuStore         ds.JSONDatastore
	StateStore        ds.JSONDatastore
	CarouselStorePath string
	ContentStorePath  string
	CountryStorePath  string
	CSSStorePath      string
	MenuStorePath     string
	StateStorePath    string
	TemplateStorePath string
	TemplateFilePath  string
	TemplateFullPath  string
	Log               *lg.Logger
	ImagePath         string
	ImageFullPath     string
}

//GetNew GetNew
func (c *Six910BackupService) GetNew() BackupService {
	return c
}

func (c *Six910BackupService) extractTarGzFile(tr *tar.Reader, h *tar.Header, dest string) error {
	var rtn error
	flname := h.Name
	c.Log.Debug("fname in extractTarGzFile: ", flname)
	switch h.Typeflag {
	case tar.TypeDir:
		derr := os.MkdirAll(dest+string(filepath.Separator)+flname, 0775)
		c.Log.Debug("MkdirAll in tar.TypeDir error in extractTarGzFile: ", derr)
		c.Log.Debug("MkdirAll in tar.TypeDir name in extractTarGzFile: ", dest+string(filepath.Separator)+flname)
		rtn = derr
	case tar.TypeReg:
		derr2 := os.MkdirAll(filepath.Dir(dest+string(filepath.Separator)+flname), 0775)
		rtn = derr2
		c.Log.Debug("MkdirAll in tar.TypeReg error in extractTarGzFile: ", derr2)
		c.Log.Debug("MkdirAll in tar.TypeReg dir name in extractTarGzFile: ", filepath.Dir(dest+string(filepath.Separator)+flname))
		if derr2 == nil {
			c.Log.Debug("MkdirAll in tar.TypeReg file name in extractTarGzFile: ", dest+string(filepath.Separator)+flname)
			writer, cerr := os.Create(dest + string(filepath.Separator) + flname)
			rtn = cerr
			c.Log.Debug("os.Create error in extractTarGzFile: ", cerr)
			if cerr == nil {
				io.Copy(writer, tr)
				err := os.Chmod(dest+string(filepath.Separator)+flname, 0664)
				c.Log.Debug("os.Chmod error in extractTarGzFile: ", err)
				rtn = err
				writer.Close()
			}
		}
	}
	return rtn
}
