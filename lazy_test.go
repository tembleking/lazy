package lazy_test

import (
	"errors"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/lazy"
)

var _ = Describe("Lazy", func() {
	var lazyValue lazy.Lazy[int]

	BeforeEach(func() {
		lazyValue = lazy.Lazy[int]{}
	})

	It("initializes the component only once", func() {
		value := lazyValue.GetOrInit(func() int { return 42 })
		Expect(value).To(Equal(42))

		value = lazyValue.GetOrInit(func() int { panic("unreachable") })
		Expect(value).To(Equal(42))

		value, err := lazyValue.GetOrTryInit(func() (int, error) { panic("unreachable") })
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(Equal(42))
	})

	It("initializes the component only once", func() {
		value, err := lazyValue.GetOrTryInit(func() (int, error) { return 42, nil })
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(Equal(42))

		value = lazyValue.GetOrInit(func() int { panic("unreachable") })
		Expect(value).To(Equal(42))

		value, err = lazyValue.GetOrTryInit(func() (int, error) { panic("unreachable") })
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(Equal(42))
	})

	When("the first initialization returns an error", func() {
		It("doesn't initialize it with error", func() {
			_, err := lazyValue.GetOrTryInit(func() (int, error) { return 0, errors.New("some error") })
			Expect(err).To(MatchError("some error"))

			value, err := lazyValue.GetOrTryInit(func() (int, error) { return 42, nil })
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(Equal(42))
		})
	})

	When("there are competing initializers in different goroutines", MustPassRepeatedly(100), func() {
		It("initializes the value with the first one that succeeds", func() {
			valuesObserved := launchInitializersThatWillSucceedInTheRange(100, 50, 60, &lazyValue)

			Expect(valuesObserved).To(HaveLen(100))
			Expect(valuesObserved).To(HaveEach(Equal(valuesObserved[0])))
		})
	})
})

func launchInitializersThatWillSucceedInTheRange(numberOfInitializers int, rangeStart int, rangeEnd int, lazyValue *lazy.Lazy[int]) (valuesObservedByInitializers []int) {
	valuesObservedByInitializers = make([]int, 0, numberOfInitializers)
	mutex := &sync.Mutex{}

	initFunc := func(i int) func() (int, error) {
		return func() (int, error) {
			if i >= rangeStart && i <= rangeEnd {
				return i, nil
			}

			return 0, errors.New("some error")
		}
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < numberOfInitializers; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			Eventually(func() error {
				_, err := lazyValue.GetOrTryInit(initFunc(i))
				return err
			}).Should(Succeed())

			value, _ := lazyValue.GetOrTryInit(initFunc(i))
			mutex.Lock()
			defer mutex.Unlock()
			valuesObservedByInitializers = append(valuesObservedByInitializers, value)
		}()
	}

	wg.Wait()
	return valuesObservedByInitializers
}
