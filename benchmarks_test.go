//The idea of this benchmarks is to measure the runtime behavior when assigning
//a regular type variable to an interface variable
//
//When i was developing a package (gohl7) i faced a decision of how to model my data types
//and actually i need to choose between the following 2 options:
//		type SampleDataType []byte
//		or
//		type SampleDataType interface{
//			buffer []byte
//		}
//I choose the first one for simplicity and because it allows me to keep using the regular
//function on slices: len, cap, make
//But it came with a not small performace implications when you try to assing a variable of that type to
//a variable of interface type.
//
//I wrote the functions: CompareContainersLength and CompareContainersLengthReuse that
//simulate that process using the two posible data modeling shown above
//
//After running (note that im using the flag -gcflags=-l to tell the compiler to not inline functions):
//		go test -gcflags=-l -run=none -bench=CacheMixedLength -cpuprofile=cprof convertT2I
//and disasembling the function CompareContainersLengthReuse or CompareContainersLength
//you will see that for the slice type there is an extra call to runtime.convT2I(SB)
//that its causing the extra overhead, the call to runtime.convT2I is not present when the type
//implementing the interface is a pointer to struct, the disasm of the function CompareContainersLengthReuse
//is shown on the file output.disasm the line 38 shows the call
//
//
//There are 3 smiliar functions for benchmaks: MixedLength, SmallLength and MediumLength and LargeLength
//the only different between them is that i wanted to know if the size of the slice will
//affect the performance but the results proved that it DOES NOT affect the performance
//
//meaning that the reader can focus only the functions: BenchmarkMixedLength and BenchmarkCacheMixedLength
//
package convertT2I

import(
	"testing"
)

//function to initialize sample data

func ValuesList() [][]byte{

	return [][]byte{
		[]byte("small slice"),
		//1K slice length
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam quis erat quis sapien rutrum mattis. Integer lobortis at neque at maximus. Vestibulum condimentum enim at imperdiet auctor. Nunc dictum ante quis sodales ultricies. Sed a dui suscipit, consectetur libero id, feugiat mauris. Nam ultricies ultricies massa, non sodales dui rutrum nec. Etiam fermentum arcu eu turpis maximus, nec rhoncus odio tincidunt. Duis semper condimentum leo, vel elementum tortor sollicitudin in. Vestibulum vel volutpat massa. Morbi accumsan tempus consequat. Fusce leo augue, venenatis vitae ultrices ac, euismod nec libero. Fusce posuere lorem lacus, in accumsan dui porta id. Nunc lacinia tincidunt eros sit amet placerat.Praesent at bibendum purus. Etiam convallis ex vitae libero convallis, sit amet tristique neque tincidunt. Praesent nec fermentum augue. Nullam magna libero, luctus et pellentesque accumsan, lacinia nec ipsum. Sed quis mattis nunc, et malesuada urna. Aenean pulvinar magna id erat rhoncus, et vulputate elit posuere"),
		//4K slice length
		[]byte("Praesent suscipit urna dolor, quis convallis est rutrum et. Cras ut augue odio. Suspendisse quis nulla sed magna vehicula imperdiet. Vestibulum pulvinar at magna et iaculis. Vivamus at euismod enim, id lacinia sem. Integer eu nisi vel dolor interdum vehicula. Ut dignissim, dui vitae maximus blandit, est orci dignissim diam, eget finibus tellus dui vitae velit. Ut commodo purus sit amet mollis accumsan. Ut aliquam, erat sed pharetra fermentum, ipsum neque maximus sem, sed rutrum nibh lacus a metus. Sed non mattis ex. Aenean porta urna quis risus fermentum tincidunt. Phasellus tincidunt malesuada justo, eget vulputate est finibus quis.Duis bibendum sed arcu at tempus. Vivamus pharetra posuere sem eu imperdiet. Donec et quam diam. Cras sagittis imperdiet libero ac faucibus. Etiam dapibus risus ut purus cursus, nec dignissim est bibendum. Donec eu magna ac neque feugiat lobortis. Sed fermentum blandit ante a hendrerit. Phasellus aliquet aliquam laoreet. Maecenas eget tempor libero. Integer finibus velit vel ullamcorper commodo. Vivamus vitae dolor in nibh imperdiet aliquet non a est. Nullam bibendum est sed nibh ullamcorper tempus eget sit amet nisl. Mauris auctor mi id nisl gravida, non elementum nulla accumsan. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vitae fringilla velit. Donec tristique bibendum nisi eget elementum.Aliquam vel massa blandit, commodo ipsum vel, tincidunt mauris. Cras mattis consectetur lorem. Proin tincidunt nisi bibendum accumsan fermentum. Ut dapibus est at erat sagittis vehicula. Donec et lectus nisl. Morbi vel finibus tortor. Sed mattis vel ex ac ullamcorper. Suspendisse vel sodales libero. Vivamus eu risus elementum, ultrices diam a, consectetur eros. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla fringilla diam eget velit pharetra, eget volutpat nulla feugiat. Etiam nulla mi, faucibus quis porta faucibus, posuere at diam. Curabitur id augue placerat, ullamcorper mi eget, vulputate lectus. Praesent sit amet eros porta, posuere est eu, ultrices diam. Proin iaculis odio viverra ligula vehicula, sed hendrerit turpis pharetra.Nunc a nunc scelerisque arcu suscipit maximus. Suspendisse nec finibus neque. Aliquam dapibus tincidunt justo in condimentum. Proin pellentesque, neque viverra tempor sollicitudin, eros neque sodales orci, id condimentum nulla justo vitae erat. Quisque sollicitudin sollicitudin hendrerit. Etiam eleifend pharetra est sed commodo. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.Nullam vitae consectetur erat. Nulla posuere ornare auctor. Etiam pellentesque risus elit, vel eleifend elit porta quis. Proin sed ultrices ligula. Suspendisse id sem non mauris euismod commodo eu quis quam. Suspendisse lectus ligula, volutpat vitae lorem id, pellentesque rhoncus ex. Suspendisse scelerisque metus auctor blandit tristique. Nam id mauris dui. Integer commodo leo vitae dapibus luctus. Suspendisse quis mi in leo mollis faucibus vel sed tellus. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Sed sodales velit vel erat ultrices tincidunt. In hac metus"),
	}
}

func BenchmarkMixedLength(b *testing.B) {

        values := ValuesList()
        //skip the creation of the values from the profile
       	b.ResetTimer()
        l := len(values)

        for n := 0; n < b.N; n++ {
        	for i := 0; i < l; i++{
        		CompareContainersLength(values[i])
        	}
        }
}

func BenchmarkSmallLength(b *testing.B) {

        values := ValuesList()
        //skip the creation of the values from the profile
       	b.ResetTimer()
        l := len(values)

        for n := 0; n < b.N; n++ {
        	for i := 0; i < l; i++{
        		CompareContainersLength(values[0])
        	}
        }
}

func BenchmarkMediumLength(b *testing.B) {

        values := ValuesList()
        //skip the creation of the values from the profile
       	b.ResetTimer()
        l := len(values)

        for n := 0; n < b.N; n++ {
        	for i := 0; i < l; i++{
        		CompareContainersLength(values[1])
        	}
        }
}

func BenchmarkLargeLength(b *testing.B) {

        values := ValuesList()
        //skip the creation of the values from the profile
       	b.ResetTimer()
        l := len(values)

        for n := 0; n < b.N; n++ {
        	for i := 0; i < l; i++{
        		CompareContainersLength(values[2])
        	}
        }
}


func BenchmarkCacheMixedLength(b *testing.B) {

        values := ValuesList()
        //skip the creation of the values from the profile
       	b.ResetTimer()
        l := len(values)

        for n := 0; n < b.N; n++ {
        	for i := 0; i < l; i++{
        		CompareContainersLengthReuse(values[i])
        	}
        }
}