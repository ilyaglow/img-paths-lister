package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ilyaglow/drw"
	"github.com/wagoodman/dive/dive/filetree"
	"github.com/wagoodman/dive/dive/image/docker"
)

func main() {
	uniq := flag.Bool("u", true, "Output only unique file paths")
	flag.Parse()
	name := flag.Arg(0)

	if name == "" {
		panic("specify image name to fetch")
	}

	engine := docker.NewResolverFromEngine()
	img, err := engine.Fetch(name)
	if err != nil {
		log.Fatal(err)
	}

	var w io.Writer
	w = os.Stdout
	if *uniq {
		drw := drw.NewWriter(os.Stdout, '\n', drw.NewMapCache(0))
		w = drw
	}

	visitor := func(node *filetree.FileNode) error {
		fmt.Fprintln(w, node.Path())
		return nil
	}

	visitEvaluator := func(node *filetree.FileNode) bool {
		return node.IsLeaf()
	}

	for _, t := range img.Trees {
		err := t.VisitDepthChildFirst(visitor, visitEvaluator)
		if err != nil {
			log.Fatal(err)
		}
	}
}
