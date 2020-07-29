package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = `
	addBlock --data DATA	"add a block to blockchain"
	printChain				"print all blocks"
`

const AddBlockCmdString = "addBlock"
const PrintChainCmdString = "printChain"

type CLI struct {
	bc *BlockChain
}

func (cli *CLI) PrintUsage() {
	fmt.Println(usage)
	os.Exit(1)
}

//func (cli *CLI) ParamCheck() {
//	fmt.Println(usage)
//	os.Exit(1)
//}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("invalid input!")
		cli.PrintUsage()
	}
	addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)

	addBlockCmdPara := addBlockCmd.String("data", "", "block transaction info!")

	switch os.Args[1] {
	case AddBlockCmdString:
		// 添加动作
		err := addBlockCmd.Parse(os.Args[2:])
		CheckErr("Run()", err)
		if addBlockCmd.Parsed() {
			if *addBlockCmdPara == "" {
				fmt.Println("addBlock data should not be empty!")
				cli.PrintUsage()
			}
			cli.AddBlock(*addBlockCmdPara)
		}
	case PrintChainCmdString:
		// 打印输出
		err := printChainCmd.Parse(os.Args[2:])
		CheckErr("Run2()", err)
		if printChainCmd.Parsed() {
			cli.PrintChain()
		}
	default:
		cli.PrintUsage()

	}
}