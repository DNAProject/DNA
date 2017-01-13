package testutil

func RunSuite(t GinkgoTestingT, name string) {
	RunDefer()

	SynchronizedBeforeSuite(func() []byte {
		RunDefer("setup")
		return nil
	}, func(data []byte) {})
	SynchronizedAfterSuite(func() {
		RunDefer("teardown")
	}, func() {})

	RegisterFailHandler(Fail)
	RunSpecs(t, name)
}
