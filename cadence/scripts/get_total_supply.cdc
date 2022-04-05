import RedSquirrelNFT from "../contracts/RedSquirrelNFT.cdc"

pub fun main(): UInt64 {
  return RedSquirrelNFT.totalSupply
}