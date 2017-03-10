package progressbar

type ProgressBar interface {
	Start(int)
	Increment()
	Finish()
}
