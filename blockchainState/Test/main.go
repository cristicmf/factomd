package main

import (
	"fmt"
	"os"

	. "github.com/FactomProject/factomd/blockchainState"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/database/databaseOverlay"
	"github.com/FactomProject/factomd/database/hybridDB"
)

const level string = "level"
const bolt string = "bolt"

func main() {
	fmt.Println("Usage:")
	fmt.Println("Test level/bolt DBFileLocation")

	if len(os.Args) < 3 {
		fmt.Println("\nNot enough arguments passed")
		os.Exit(1)
	}
	if len(os.Args) > 3 {
		fmt.Println("\nToo many arguments passed")
		os.Exit(1)
	}

	levelBolt := os.Args[1]

	if levelBolt != level && levelBolt != bolt {
		fmt.Println("\nFirst argument should be `level` or `bolt`")
		os.Exit(1)
	}
	path := os.Args[2]

	var dbase *hybridDB.HybridDB
	var err error
	if levelBolt == bolt {
		dbase = hybridDB.NewBoltMapHybridDB(nil, path)
	} else {
		dbase, err = hybridDB.NewLevelMapHybridDB(path, false)
		if err != nil {
			panic(err)
		}
	}

	CheckDatabase(dbase)
}

func CheckDatabase(db interfaces.IDatabase) {
	if db == nil {
		return
	}

	dbo := databaseOverlay.NewOverlay(db)
	bs := new(BlockchainState)
	bs.Init()
	bl := new(BalanceLedger)
	bl.Init()

	dBlock, err := dbo.FetchDBlockHead()
	if err != nil {
		panic(err)
	}
	if dBlock == nil {
		panic("DBlock head not found")
	}

	fmt.Printf("\tStarting\n")

	max := int(dBlock.GetDatabaseHeight())
	if max > 10000 {
		//max = 10000
	}

	specialBlocks := []int{22880, 22882, 22938, 22946, 22972, 22973, 23261, 31451, 49225, 50145, 54339, 57198, 62763, 67791, 69064, 70411}
	nextSpacial := specialBlocks[0]
	specialBlocks = specialBlocks[1:]

	//for i := 0; i < int(dlock.GetDatabaseHeight()); i++ {
	for i := 0; i < max; i++ {
		set := FetchBlockSet(dbo, i)
		if i%1000 == 0 {
			fmt.Printf("\"%v\", //%v\n", set.DBlock.DatabasePrimaryIndex(), set.DBlock.GetDatabaseHeight())
		}
		if i == nextSpacial {
			if len(specialBlocks) > 0 {
				nextSpacial = specialBlocks[0]
				specialBlocks = specialBlocks[1:]
			}
			//ec := FetchFloatingBlockBefore(dbo, i)
			//bs.ProcessECBlock(ec)
		}

		err := bs.ProcessBlockSet(set.DBlock, set.ABlock, set.FBlock, set.ECBlock, set.EBlocks)
		if err != nil {
			panic(err)
		}
		/*
			err = bl.ProcessFBlock(set.FBlock)
			if err != nil {
				panic(err)
			}
		*/
	}
	fmt.Printf("\tFinished!\n")

	b, err := bs.MarshalBinaryData()
	if err != nil {
		panic(err)
	}
	fmt.Printf("BS size - %v\n", len(b))

	b, err = bl.MarshalBinaryData()
	if err != nil {
		panic(err)
	}
	fmt.Printf("BL size - %v\n", len(b))

	fmt.Printf("Expired - %v\n", Expired)
	fmt.Printf("LatestReveal - %v\n", LatestReveal)
	fmt.Printf("TotalEntries - %v\n", TotalEntries)
}

type BlockSet struct {
	ABlock  interfaces.IAdminBlock
	ECBlock interfaces.IEntryCreditBlock
	FBlock  interfaces.IFBlock
	DBlock  interfaces.IDirectoryBlock
	EBlocks []interfaces.IEntryBlock
}

func FetchFloatingBlockBefore(dbo interfaces.DBOverlay, index int) interfaces.IEntryCreditBlock {
	dBlock, err := dbo.FetchDBlockByHeight(uint32(index))
	if err != nil {
		panic(err)
	}
	ec := dBlock.GetDBEntries()[1] //EC Block
	ecBlock, err := dbo.FetchECBlock(ec.GetKeyMR())
	if err != nil {
		panic(err)
	}
	ecBlock, err = dbo.FetchECBlock(ecBlock.GetHeader().GetPrevHeaderHash())
	if err != nil {
		panic(err)
	}
	return ecBlock
}

func FetchBlockSet(dbo interfaces.DBOverlay, index int) *BlockSet {
	bs := new(BlockSet)

	dBlock, err := dbo.FetchDBlockByHeight(uint32(index))
	if err != nil {
		panic(err)
	}
	bs.DBlock = dBlock

	if dBlock == nil {
		return bs
	}
	entries := dBlock.GetDBEntries()
	for _, entry := range entries {
		switch entry.GetChainID().String() {
		case "000000000000000000000000000000000000000000000000000000000000000a":
			aBlock, err := dbo.FetchABlock(entry.GetKeyMR())
			if err != nil {
				panic(err)
			}
			bs.ABlock = aBlock
			break
		case "000000000000000000000000000000000000000000000000000000000000000c":
			ecBlock, err := dbo.FetchECBlock(entry.GetKeyMR())
			if err != nil {
				panic(err)
			}
			bs.ECBlock = ecBlock
			break
		case "000000000000000000000000000000000000000000000000000000000000000f":
			fBlock, err := dbo.FetchFBlock(entry.GetKeyMR())
			if err != nil {
				panic(err)
			}
			bs.FBlock = fBlock
			break
		default:
			eBlock, err := dbo.FetchEBlock(entry.GetKeyMR())
			if err != nil {
				panic(err)
			}
			bs.EBlocks = append(bs.EBlocks, eBlock)
			break
		}
	}

	return bs
}
