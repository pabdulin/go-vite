package nodemanager

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"

	"github.com/vitelabs/go-vite/common"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/node"
)

type CheckChainNodeManager struct {
	ctx  *cli.Context
	node *node.Node
	log  log15.Logger
}

func NewCheckChainNodeManager(ctx *cli.Context, maker NodeMaker) (*CheckChainNodeManager, error) {
	node, err := maker.MakeNode(ctx)
	if err != nil {
		return nil, err
	}

	// single mode
	node.Config().Single = true
	node.ViteConfig().Net.Single = true

	// no miner
	node.Config().MinerEnabled = false
	node.ViteConfig().Producer.Producer = false

	// no ledger gc
	ledgerGc := false
	node.Config().LedgerGc = &ledgerGc
	node.ViteConfig().Chain.LedgerGc = ledgerGc

	return &CheckChainNodeManager{
		ctx:  ctx,
		node: node,
		log:  log15.New("module", "checkChainCMD"),
	}, nil
}

func (nodeManager *CheckChainNodeManager) Start() error {
	node := nodeManager.node

	err := StartNode(nodeManager.node)
	if err != nil {
		return err
	}

	c := node.Vite().Chain()
	fmt.Println("start check.")
	// check recent blocks
	nodeManager.log.Info("start check recent blocks")
	if err := c.CheckRecentBlocks(); err != nil {
		common.Crit(err.Error(), "check_chain", "recent_blocks")
	}
	nodeManager.log.Info("finish checking recent blocks")
	fmt.Println("check recent blocks success.")
	// check redo
	nodeManager.log.Info("start check redo")
	if err := c.CheckRedo(); err != nil {
		common.Crit(err.Error(), "check_chain", "redo")
	}
	nodeManager.log.Info("finish checking redo")
	fmt.Println("check redo success.")

	// check onroad
	nodeManager.log.Info("start check onroad")
	if err := c.CheckOnRoad(); err != nil {
		common.Crit(err.Error(), "check_chain", "onroad")
	}
	nodeManager.log.Info("finish checking onroad")

	fmt.Println("check onroad success.")

	fmt.Println("check success.")

	return nil
}

func (nodeManager *CheckChainNodeManager) Stop() error {

	StopNode(nodeManager.node)

	return nil
}

func (nodeManager *CheckChainNodeManager) Node() *node.Node {
	return nodeManager.node
}
