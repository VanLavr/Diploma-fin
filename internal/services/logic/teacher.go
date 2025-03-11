package logic

import "context"

type TeacherUsecase interface {
	SetDate(context.Context, string, string, int64) error
}
