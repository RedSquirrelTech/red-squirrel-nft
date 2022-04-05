import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import RedSquirrelNFT from "../contracts/RedSquirrelNFT.cdc"

transaction(recipient: Address, name: String, description: String, thumbnail: String) {
    let minter: &RedSquirrelNFT.NFTMinter

    prepare(signer: AuthAccount) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.borrow<&RedSquirrelNFT.NFTMinter>(from: RedSquirrelNFT.MinterStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")
    }

    execute {
        // get the public account object for the recipient
        let recipient = getAccount(recipient)

        // borrow the recipient's public NFT collection reference
        let receiver = recipient
            .getCapability(RedSquirrelNFT.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        // mint the NFT and deposit it to the recipient's collection
        self.minter.mintNFT(
            recipient: receiver,
            name: name,
            description: description,
            thumbnail: thumbnail,
        )
    }
}
