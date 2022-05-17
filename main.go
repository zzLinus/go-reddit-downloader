//			my very first go project !!!
//            --==============--
//   .-==-.===oooo=oooooo=ooooo===--===-
//  .==  =o=oGGGGGGo=oo=oGGGGGGGG=o=  oo-
//  -o= oo=G .=GGGGGo=o== .=GGGGG=ooo o=-
//   .-=oo=o==oGGGGG=oo=oooGGGGGo=oooo.
//    -ooooo=oooooo=.   .==ooo==oooooo-
//    -ooooooooooo====_====ooooooooooo=
//    -oooooooooooo==#.#==ooooooooooooo
//    -ooooooooooooo=#.#=oooooooooooooo
//    .oooooooooooooooooooooooooooooooo.
//     oooooooooooooooooooooooooooooooo.
//   ..oooooooooooooooooooooooooooooooo..
// -=o-=ooooooooooooooooooooooooooooooo-oo.
// .=- oooooooooooooooooooooooooooooooo-.-
//    .oooooooooooooooooooooooooooooooo-
//    -oooooooooooooooooooooooooooooooo-
//    -oooooooooooooooooooooooooooooooo-
//    -oooooooooooooooooooooooooooooooo-
//    .oooooooooooooooooooooooooooooooo
//     =oooooooooooooooooooooooooooooo-
//     .=oooooooooooooooooooooooooooo-
//       -=oooooooooooooooooooooooo=.
//      =oo====oooooooooooooooo==-oo=-
//     .-==-    .--=======---     .==-

package main

import (
	"fmt"
	"os"

	"github.com/zzLinus/GoTUITODOList/tuiapp"
)

func main() {
	if err := tuiapp.New().Start(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
