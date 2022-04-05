import MetadataViews from "../contracts/standard/MetadataViews.cdc"
import RedSquirrelNFT from "../contracts/RedSquirrelNFT.cdc"

pub struct NFT {
    pub let redSquirrelID: UInt64
    pub let resourceID: UInt64
    pub let owner: Address
    pub let type: String
    pub let name: String
    pub let description: String
    pub let thumbnail: String

    init(
        redSquirrelID: UInt64,
        resourceID: UInt64,
        owner: Address,
        nftType: String,
        name: String,
        description: String,
        thumbnail: String,
    ) {
        self.redSquirrelID = redSquirrelID
        self.resourceID = resourceID
        self.owner = owner
        self.type = nftType
        self.name = name
        self.description = description
        self.thumbnail = thumbnail
    }
}

pub fun main(address: Address, id: UInt64): NFT {
    let account = getAccount(address)

    let collection = account.getCapability(RedSquirrelNFT.CollectionPublicPath)
        .borrow<&{RedSquirrelNFT.RedSquirrelNFTCollectionPublic}>()
        ?? panic("Could not borrow a reference to the collection")

    let nft = collection.borrowRedSquirrelNFT(id: id)!

    // Get the basic display information for this NFT
    let view = nft.resolveView(Type<MetadataViews.Display>())!

    let display = view as! MetadataViews.Display
    
    let owner: Address = nft.owner!.address!
    let nftType = nft.getType()

    return NFT(
        redSquirrelID: nft.id,
        resourceID: nft.uuid,
        owner: owner,
        nftType: nftType.identifier,
        name: display.name,
        description: display.description,
        thumbnail: display.thumbnail.uri(),
    )
}