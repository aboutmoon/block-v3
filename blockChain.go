package main

import (
	"github.com/boltdb/bolt"
	"os"
)

const dbFile = "blockChain.db"
const blockBucket = "bucket"
const lastHashKey = "key"

type BlockChain struct {
	//blocks []*Block
	// 数据库操作句柄
	db *bolt.DB
	// 尾巴, 表示最后一个区块的哈希值
	tail []byte
}

func NewBlockChain() *BlockChain  {
	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("NewBlockChain", err)
	var lastHash []byte
	// db.View
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil {
			// 取出最后区块的Hash值
			lastHash = bucket.Get([]byte(lastHashKey))
		} else {
			// 创建创世块
			genesis := NewGenesisBlock()
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			CheckErr("NewBlockChain2", err)
			bucket.Put(genesis.Hash, genesis.Serialize()) // Todo
			CheckErr("NewBlockChain3", err)
			bucket.Put([]byte(lastHashKey), genesis.Hash)
			CheckErr("NewBlockChain4", err)
			lastHash = genesis.Hash
		}
		return nil
	})

	return &BlockChain{db, lastHash}
}

func (bc *BlockChain)AddBlock(data string)  {
	var prevBlockHash []byte

	bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}

		prevBlockHash =  bucket.Get([]byte(lastHashKey))
		return nil
	})
	block := NewBlock(data, prevBlockHash)

	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}

		err := bucket.Put(block.Hash, block.Serialize()) // Todo
		CheckErr("AddBlock1", err)
		err = bucket.Put([]byte(lastHashKey), block.Hash)
		CheckErr("AddBlock2", err)
		bc.tail = block.Hash
		return nil
	})

	CheckErr("AddBlock2", err)
}

// 迭代器，就是一个对象， 它里面包含了一个游标，一直向前（向后）移动

type BlockChainIterator struct {
	currHash []byte
	db *bolt.DB
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{currHash: bc.tail, db: bc.db}
}

func (it *BlockChainIterator)Next()  (block *Block) {
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}

		data := bucket.Get(it.currHash)
		block = Deserialize(data)
		it.currHash = block.PrevBlockHash
		return nil
	})
	CheckErr("Next()", err)
	return
}

