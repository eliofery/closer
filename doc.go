/*
Package closer ensures a clean exit for your Go app.

The aim of this package is to provide an universal way to catch the event of application’s exit
and perform some actions before it’s too late.

Includes [Closer] that allows you to add functions to be closed and wait for them to be closed.

Demonstrates how to use closer.

	// Create new instance of Closer.
	clr := closer.New()

	// Some functions to be closed.
	clr.Add(func() error {
		time.Sleep(1 * time.Second)

		if rand.IntN(2) == 1 {
			return nil
		}

		return errors.New("error close something")
	})

	// Wait signals about closing.
	clr.Wait()

	// Close all functions.
	if err := clr.Close(); err != nil {
		log.Println(err)
	}
*/
package closer
