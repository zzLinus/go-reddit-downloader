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

// TODO: 真是你妈一拖屎现在。。。extractor返回的类结构体没写，现在就返回一个url，视频大小，音频url也没有,
// 还需要用正则提取一下h1 tag用来做文件路径的一部分
// 下载视频的路经放到一个特定的文件夹里，
// TODO: support gif download
// TODO: different resolution option

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
