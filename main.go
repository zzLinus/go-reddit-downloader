//				My Very First Go Project !!!
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

// TODO: 真是你妈一拖屎现在。。。need to extract & download audio url if it's a mp4 content and use ffmpeg
// or something to join them
// TODO: different resolution option
// TODO: be able to let use chose what resolution they what to download
// TODO: fix the shitty UI
// 今天他妈的情人节我还在这跟这坨屎耗着我操你妈气死我了啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊

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
