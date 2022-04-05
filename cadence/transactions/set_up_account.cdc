import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import RedSquirrelNFT from "../contracts/RedSquirrelNFT.cdc"

transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&RedSquirrelNFT.Collection>(from: RedSquirrelNFT.CollectionStoragePath) == nil {
            // create a new empty collection
            let collection <- RedSquirrelNFT.createEmptyCollection()
            
            // save it to the account
            signer.save(<- collection, to: RedSquirrelNFT.CollectionStoragePath)

            // Creates a public capability for the collection so that other users can publicly access electable attributes.
            // The pieces inside of the brackets specify the type of the linked object, and only expose the fields and
            // functions on those types.
            signer.link<&RedSquirrelNFT.Collection{NonFungibleToken.CollectionPublic, RedSquirrelNFT.RedSquirrelNFTCollectionPublic}>(
                RedSquirrelNFT.CollectionPublicPath, target: RedSquirrelNFT.CollectionStoragePath
            )
        }
    }
}