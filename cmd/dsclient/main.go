package main

import (
	"context"
	"fmt"
	"os"

	"github.com/filedrive-team/filehelper"
	"github.com/filedrive-team/filehelper/dataset"
	"github.com/filedrive-team/go-ds-cluster/clusterclient"
	"github.com/filedrive-team/go-ds-cluster/config"
	"github.com/ipfs/go-blockservice"
	ds "github.com/ipfs/go-datastore"
	dsmount "github.com/ipfs/go-datastore/mount"
	dss "github.com/ipfs/go-datastore/sync"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	dshelp "github.com/ipfs/go-ipfs-ds-help"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	log "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var logging = log.Logger("dsclient")

func init() {
	log.SetLogLevel("*", "INFO")
}

func main() {
	local := []*cli.Command{
		addCmd,
		importDatasetCmd,
	}

	app := &cli.App{
		Name:     "dsclient",
		Commands: local,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

var importDatasetCmd = &cli.Command{
	Name:  "import-dataset",
	Usage: "import files from the specified dataset",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "dscluster-cfg",
			Required: true,
			Usage:    "specify the dscluster config path",
		},
		&cli.IntFlag{
			Name:  "retry",
			Value: 5,
			Usage: "retry write file to datastore",
		},
		&cli.IntFlag{
			Name:  "retry-wait",
			Value: 1,
			Usage: "sleep time before a retry",
		},
	},
	Action: func(c *cli.Context) (err error) {
		ctx := context.Background()
		dscluster := c.String("dscluster-cfg")

		targetPath := c.Args().First()
		targetPath, err = homedir.Expand(targetPath)
		if err != nil {
			return err
		}
		if !filehelper.ExistDir(targetPath) {
			return xerrors.Errorf("Unexpected! The path to dataset does not exist")
		}

		return dataset.Import(ctx, targetPath, dscluster, c.Int("retry"), c.Int("retry-wait"))
	},
}

var addCmd = &cli.Command{
	Name:  "add",
	Usage: "import single file to ds-cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "conf",
			Usage: "specify the dscluster config path",
			Value: "config.json",
		},
	},
	Action: func(c *cli.Context) error {
		cfg, err := config.ReadConfig(c.String("conf"))
		if err != nil {
			return err
		}
		target := c.Args().First()
		target, err = homedir.Expand(target)
		if err != nil {
			return err
		}
		var ds ds.Datastore
		ds, err = clusterclient.NewClusterClient(context.Background(), cfg)
		if err != nil {
			return err
		}
		ds = dsmount.New([]dsmount.Mount{
			{
				Prefix:    bstore.BlockPrefix,
				Datastore: ds,
			},
		})
		bs2 := bstore.NewBlockstore(dss.MutexWrap(ds))
		dagServ := merkledag.NewDAGService(blockservice.New(bs2, offline.Exchange(bs2)))

		// cidbuilder
		cidBuilder, err := merkledag.PrefixForCidVersion(0)
		if err != nil {
			return err
		}

		files := filehelper.FileWalkAsync([]string{target})

		for item := range files {
			fileNode, err := filehelper.BuildFileNode(item, dagServ, cidBuilder)
			if err != nil {
				return err
			}
			k := dshelp.MultihashToDsKey(fileNode.Cid().Hash())
			logging.Infof("imported file: %s, root: %s, key: %s", item.Path, fileNode, k)
		}

		return nil
	},
}
