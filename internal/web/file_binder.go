package web

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	typeMultipartFileHeader      = reflect.TypeOf((*multipart.FileHeader)(nil))
	typeMultipartSliceFileHeader = reflect.TypeOf(([]*multipart.FileHeader)(nil))
)

type BindFunc func(interface{}, echo.Context) error

func (fn BindFunc) Bind(i interface{}, ctx echo.Context) error {
	return fn(i, ctx)
}

func NewFileBinder(b echo.Binder) echo.Binder {
	return BindFunc(func(i interface{}, ctx echo.Context) error {
		err := b.Bind(i, ctx)
		if err == nil {
			ctype := ctx.Request().Header.Get(echo.HeaderContentType)
			// if bind form
			if strings.HasPrefix(ctype, echo.MIMEApplicationForm) || strings.HasPrefix(ctype, echo.MIMEMultipartForm) {
				// get form files
				var form *multipart.Form
				form, err = ctx.MultipartForm()
				if err == nil {
					err = BindFile(i, ctx, form.File)
				}
			}
		}
		return err
	})
}

func BindFile(i interface{}, ctx echo.Context, files map[string][]*multipart.FileHeader) error {
	iValue := reflect.Indirect(reflect.ValueOf(i))
	// check bind type is struct pointer
	if iValue.Kind() != reflect.Struct {
		return fmt.Errorf("BindFile input not is struct pointer, indirect type is %s", iValue.Type().String())
	}

	iType := iValue.Type()
	for i := 0; i < iType.NumField(); i++ {
		fType := iType.Field(i)
		// check canset field
		fValue := iValue.Field(i)
		if !fValue.CanSet() {
			continue
		}
		// revc type must *multipart.FileHeader or []*multipart.FileHeader
		switch fType.Type {
		case typeMultipartFileHeader:
			file := getFiles(files, fType.Name, fType.Tag.Get("form"))
			if len(file) > 0 {
				fValue.Set(reflect.ValueOf(file[0]))
			}
		case typeMultipartSliceFileHeader:
			file := getFiles(files, fType.Name, fType.Tag.Get("form"))
			if len(file) > 0 {
				fValue.Set(reflect.ValueOf(file))
			}
		}
	}
	return nil
}

func getFiles(files map[string][]*multipart.FileHeader, names ...string) []*multipart.FileHeader {
	for _, name := range names {
		file, ok := files[name]
		if ok {
			return file
		}
	}
	return nil
}
