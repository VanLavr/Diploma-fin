package logic

import (
	"context"

	"github.com/xuri/excelize/v2"
)

type FileUsecase interface {
	ParseFile(context.Context, *excelize.File) error
}
