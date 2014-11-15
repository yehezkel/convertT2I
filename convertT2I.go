package convertT2I

//A dummy interface
type Container interface{
	Length() int
}


//type that implement the interface
//note that the inner type is an slice
type SliceContainer []byte

//type that implement the interface
//note here that this is an struct
type StructContainer struct{
	buffer []byte
}

//SliceContainer Length method implementation
func (c SliceContainer) Length() int{
	return len(c)
}

func MakeSliceContainer(value []byte) SliceContainer{

	var s SliceContainer
	s = value

	return s
}


//StructContainer LengthMethod Implementation
func (c *StructContainer) Length() int{
	return len(c.buffer)
}

func NewStructContainer(value []byte) *StructContainer{

	return &StructContainer{
		buffer: value,
	}
}


func CompareContainersLength(value []byte) bool{
	var c1 Container
	var c2 Container

	//this are the points to benchmark
	//the idea is to see the runtime behavior when
	//the type implementing the interface has and slice as inner type
	//or a pointer to an struct
	c1 = MakeSliceContainer(value)

	c2 = NewStructContainer(value)

	return (c1.Length() == c2.Length())
}

//simulate a cache of StructContainer to avoid the memory allocation
//on function NewStructContainer when the struct is created
var (
	//this will be use a poll to avoid
	reuse = new(StructContainer)
)

//define a new way to create the StructContainers using
//the cache variable
func ReuseStructContainer(value []byte) *StructContainer{

	reuse.buffer = value
	return reuse
}


//rewrite the CompareContainersLength but this time using
//the ReuseStructContainer function instead of NewStructContainer
func CompareContainersLengthReuse(value []byte) bool{
	var c1 Container
	var c2 Container

	c1 = MakeSliceContainer(value)

	c2 = ReuseStructContainer(value)

	return (c1.Length() == c2.Length())
}