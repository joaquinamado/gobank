package services

// whtmltopdf needs to run on the main thread
// so that's why this package is necesary.

var Run = make(chan func())

// CallFunc calls the provided function on the main thread.
func CallFunc(f func() error) error {
	err := make(chan error, 1)
	Run <- func() {
		err <- f()
	}
	return <-err
}
