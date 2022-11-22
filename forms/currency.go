package forms

type CurrencyRequest struct {
	Fsyms string `form:"fsyms" binding:"required"`
	Tsyms string `form:"tsyms" binding:"required"`
}
