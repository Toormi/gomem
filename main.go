package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	_ "go.uber.org/automaxprocs"
	"gomem/gccompare"
	_map "gomem/map"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("cpuprofile")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	//gobuntdb.TestBuntDb()
	//gomem.TestMemDb()
	//memory.TestMemory()
	gomem := cli.NewApp()
	gomem.Commands = cli.Commands{
		{
			Name:  "testmap",
			Usage: "./gomem testmap",
			Action: func(c *cli.Context) error {
				n := c.Int("n")
				m := c.Int("m")
				_map.TestMap(n, m)
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "num,n",
					Value: 50000,
					Usage: "num",
				},
				cli.IntFlag{
					Name:  "m",
					Value: 1000,
					Usage: "m",
				},
			},
		},
		{
			Name:  "gccompare",
			Usage: "./gomem gccompare",
			Action: func(c *cli.Context) error {
				n := c.Int("entry")
				gccompare.GcCompare(n)
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "entry,e",
					Value: 50000,
					Usage: "entry",
				},
			},
		},
	}
	if err = gomem.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
