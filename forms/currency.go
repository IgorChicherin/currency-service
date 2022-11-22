package forms

type CurrencyRequest struct {
	fsyms string `form:"fsyms" binding:"required"`
	tsyms string `form:"tsyms" binding:"required"`
}
