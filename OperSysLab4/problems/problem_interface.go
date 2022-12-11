package problems

type Problem interface {
	Init(args interface{})
	Run()
	Description() string
}
