//        My Very First Go Project !!!
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

// TODO: add a list UI component to be able to let user choose what resolution they what to download
// TODO: current logging is fake,link it with zhe download state
// TODO: be able to download from mutible at the same time(play around with the "goroutine" stuff)
// TODO: the process bar when quiting is stupid as fuck,has to fix it sometimes...

package main

import (
	"log"
	"os"

	_ "github.com/zzLinus/GoRedditDownloader/extractor/reddit"
	"github.com/zzLinus/GoRedditDownloader/tuiapp"
)

func main() {
	if err := tuiapp.New().Start(); err != nil {
		log.Fatalf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
