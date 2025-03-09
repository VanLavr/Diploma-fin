package student

type Group struct {
	ID   int64
	Name string
}

const (
	EmailText string = "Здравствуйте,\n\nСтудент %s %s %s (Группа: %s) запрашивает дату пересдачи по: %s.\n\nВсего наилучшего, администрация РТУ МИРЭА,\nДату можно назначить на сайте"
)
