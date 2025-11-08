import { ApiPromise, WsProvider } from "@polkadot/api";
import { ContractPromise } from "@polkadot/api-contract";
import { Keyring } from '@polkadot/keyring';
import { mnemonicGenerate } from '@polkadot/util-crypto';

import BN from "bn.js";
import minimist from "minimist";

import escrowContractMetadata from "./polkadot_contracts/escrow_contract.json" with { type: "json" };
import dealContractMetadata from "./polkadot_contracts/deal_contract.json" with { type: "json" };
import registry from "./polkadot_contracts/registry.json" with { type: "json" };

import BigNumber from "bignumber.js";

function convertBalance(number) {
    let bigNumber = new BigNumber(number.replace(/,/g, ''));
    let dividedResult = bigNumber.dividedBy(new BigNumber(1e10));
    return dividedResult.toNumber();
}

/* Utility functions */
function buildGasLimit(api) {
    /**
     * Return WeightsV2 gasLimit identical to -1 (unlimited) in WeightsV1 
     */

    return api.registry.createType("SpWeightsWeightV2Weight", {
        refTime: new BN("4919333633"), 
        proofSize: new BN("126072"),
    });
}

const keyring = new Keyring({ type: 'sr25519' });
const argv = minimist(process.argv.slice(2))

async function mintEscrow(api, escrowContract, marketplaceId, sellerAddress, cid, _token, _value, nativeTokenValue, userWallet) {
    let res;
    await new Promise((resolve, reject) => {
        escrowContract.tx["escrow::mint"](
            {
                gasLimit: buildGasLimit(api),
                value: nativeTokenValue
            },
            { u16: marketplaceId }, sellerAddress, cid, null, null
        ).signAndSend(userWallet, (result) => {
            if (result.isInBlock) {
                var resultHash = result.status.asInBlock.toHex();
                var escrowId = -1;
                if (result.contractEvents) {
                    escrowId = parseInt(result.contractEvents[0].args[0].toHuman());
                }
                res = [resultHash, escrowId]
                resolve(res)
            } else if (result.isError) {
                reject(result.dispatchError);
            }
        });
    });

    return res;
}

async function escrowCancel(api, escrowContract, marketplaceWallet, escrowId) {
    let resultHash;
    await new Promise((resolve, reject) => {
        escrowContract.tx["escrow::cancel"](
            {
                gasLimit: buildGasLimit(api),
            },
            { u64: escrowId }
        ).signAndSend(marketplaceWallet, (result) => {
            if (result.isInBlock) {
                resultHash = result.status.asInBlock.toHex();
                resolve(resultHash);
            } else if (result.isError) {
                reject(result.dispatchError);
            }
        });
    });

    return resultHash;
}

async function escrowRelease(api, escrowContract, marketplaceWallet, escrowId) {
    let resultHash;

    await new Promise((resolve, reject) => {
        escrowContract.tx["escrow::release"](
            {
                gasLimit: buildGasLimit(api),
            },
            { u64: escrowId }
        ).signAndSend(marketplaceWallet, (result) => {
            if (result.isInBlock) {
                resultHash = result.status.asInBlock.toHex();
                resolve(resultHash);
            } else if (result.isError) {
                reject(result.dispatchError);
            }
        });
    });
    return resultHash;
}

async function psp34TotalSupply(api, escrowContract, from) {
    let { output } = await escrowContract.query["psp34::totalSupply"](
        from,
        { gasLimit: buildGasLimit(api) },
    );

    if (output?.isOk) {
        return output.asOk.toNumber();
    } else {
        return 0;
    }
}

async function psp34OwnersTokenByIndex(api, escrowContract, address, index, from) {
    let { output } = await escrowContract.query["psp34Enumerable::ownersTokenByIndex"](
        from,
        { gasLimit: buildGasLimit(api) },
        address, index
    );

    if (output?.isOk) {
        return parseInt(output.asOk.toHuman().Ok.U64)
    } else {
        return 0;
    }
}

async function psp34OwnerOf(api, escrowContract, escrowId, from) {
    let { output } = await escrowContract.query["psp34::ownerOf"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: escrowId },
    );

    if (output?.isOk) {
        return output.asOk.toHuman();
    } else {
        return ""
    }
}

async function psp34BalanceOf(api, escrowContract, address, from) {
    let { output } = await escrowContract.query["psp34::balanceOf"](
        from,
        { gasLimit: buildGasLimit(api) },
        address,
    );

    if (output?.isOk) {
        return output.asOk.toHuman();
    } else {
        return 0
    }
}

async function escrowShippingStatus(api, escrowContract, escrowId, from) {
    let { output } = await escrowContract.query["escrow::shippingStatus"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: escrowId },
    );

    if (output?.isOk) {
        return output.asOk.toHuman();
    } else {
        return ""
    }
}

async function escrowStatus(api, escrowContract, escrowId, from) {
    let { output } = await escrowContract.query["escrow::escrowStatus"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: escrowId },
    );

    if (output?.isOk) {
        return output.asOk.toHuman();
    } else {
        return ""
    }
}

async function escrowValue(api, escrowContract, escrowId, from) {
    let { output } = await escrowContract.query["escrow::value"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: escrowId },
    );

    if (output?.isOk) {
        return output.asOk.toHuman()
    } else {
        return ""
    }
}

async function escrowTokenUri(api, escrowContract, escrowId, from) {
    let { output } = await escrowContract.query["escrow::tokenUri"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: escrowId },
    );

    if (output?.isOk) {
        return output.asOk.toHuman()
    } else {
        return ""
    }
}

async function escrowSellerAddress(api, escrowContract, escrowId, from) {
    let { output } = await escrowContract.query["escrow::sellerAddress"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: escrowId },
    );

    if (output?.isOk) {
        return output.asOk.toHuman()
    } else {
        return ""
    }
}

async function escrowMarketplaceId(api, escrowContract, escrowId, from) {
    let { output } = await escrowContract.query["escrow::marketplaceId"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: escrowId },
    );

    if (output?.isOk) {
        return parseInt(output.asOk.toHuman().U16)
    } else {
        return ""
    }
}

async function escrowInfo(api, escrowContract, escrowId, from) {
    let es = await escrowStatus(api, escrowContract, escrowId, from);
    let ss = await escrowShippingStatus(api, escrowContract, escrowId, from);
    let value = await escrowValue(api, escrowContract, escrowId, from);
    let sellerAddress = await escrowSellerAddress(api, escrowContract, escrowId, from);
    let marketplaceId = await escrowMarketplaceId(api, escrowContract, escrowId, from);
    let uri = await escrowTokenUri(api, escrowContract, escrowId, from);
    let owner = await psp34OwnerOf(api, escrowContract, escrowId, from);


    return {
        id: escrowId,
        owner: owner,
        shipping_status: {
            status: ss[0],
            timestamp: parseInt(ss[1].replace(/,/g, '')),
        },
        escrow_status: {
            status: es[0],
            timestamp: parseInt(es[1].replace(/,/g, '')),
        },
        value: convertBalance(value),
        uri: uri,
        marketplace_id: marketplaceId,
        seller_address: sellerAddress,
    }
}

// CLI Methods

if (argv?.mode == "escrow_info" && argv?.escrow_id) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])
    let escrowContract = new ContractPromise(
        api,
        escrowContractMetadata,
        registry.contracts.escrow
    );

    let ei = await escrowInfo(api, escrowContract, argv?.escrow_id, marketplaceWallet.address);
    console.log(JSON.stringify(ei));
}

if (argv?.mode == "escrow_total_supply") {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])
    let escrowContract = new ContractPromise(
        api,
        escrowContractMetadata,
        registry.contracts.escrow
    );
    const totalSupply = await psp34TotalSupply(api, escrowContract, marketplaceWallet.address);

    console.log(totalSupply);
}

if (argv?.mode == "escrow_owner_of" && argv?.escrow_id) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])
    let escrowContract = new ContractPromise(
        api,
        escrowContractMetadata,
        registry.contracts.escrow
    );
    const ownerOf = await psp34OwnerOf(api, escrowContract, argv?.escrow_id, marketplaceWallet.address);
    console.log(ownerOf);
}

if (argv?.mode == "create_wallet") {
    const mnemonic = mnemonicGenerate();
    const keypair = keyring.addFromUri(mnemonic, { name: 'user pair' }, 'ed25519');
    console.log(JSON.stringify({
        mnemonic: mnemonic,
        public_key: keypair.address
    }));
}

if (argv?.mode == "get_balance" && argv?.wallet) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const { _, data: balance } = await api.query.system.account(argv?.wallet);
   
    console.log(JSON.stringify({
        free: convertBalance(balance.free.toLocaleString('fullwide', {useGrouping:false})),
        reserved: convertBalance(balance.reserved.toLocaleString('fullwide', {useGrouping:false})),
        frozen: convertBalance(balance.frozen.toLocaleString('fullwide', {useGrouping:false})),
    }));
}

if (argv?.mode == "send_dot" && argv?.user_wallet_mnemonic && argv?.amount && argv?.to) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });

    const senderAccount = keyring.addFromUri(argv?.user_wallet_mnemonic, { name: 'user pair' }, 'ed25519');

    const transfer = api.tx.balances.transferKeepAlive(argv?.to, argv?.amount);
    const hash = await transfer.signAndSend(senderAccount);

    console.log(JSON.stringify({
        address: senderAccount.address,
        hash: hash
    }));
}

if (argv?.mode == "list_wallet_escrows" && argv?.wallet) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])

    let escrowContract = new ContractPromise(
        api,
        escrowContractMetadata,
        registry.contracts.escrow
    );

    let numEscrows = await psp34BalanceOf(api, escrowContract, argv?.wallet, marketplaceWallet.address);
    var deals = []

    for (var i = 0; i < numEscrows; i++){
        var dealId = await psp34OwnersTokenByIndex(
            api, 
            escrowContract, 
            argv?.wallet, 
            i, 
            marketplaceWallet.address,
        );
        let escrow = await escrowInfo(api, escrowContract, dealId, marketplaceWallet.address);
        deals.push(escrow)
    }

    console.log(JSON.stringify(deals));
}

if (argv?.mode == "list_wallet_deals" && argv?.wallet) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])

    let dealContract = new ContractPromise(
        api,
        dealContractMetadata,
        registry.contracts.deal
    );

    let escrowContract = new ContractPromise(
        api,
        escrowContractMetadata,
        registry.contracts.escrow
    );

    let numDeals = await psp34BalanceOf(api, dealContract, argv?.wallet, marketplaceWallet.address);
    var escrows = []

    for (var i = 0; i < numDeals; i++){
        var dealId = await psp34OwnersTokenByIndex(
            api, 
            dealContract, 
            argv?.wallet, 
            i, 
            marketplaceWallet.address,
        );
        let escrow = await escrowInfo(api, escrowContract, dealId, marketplaceWallet.address);
        escrows.push(escrow)
    }

    console.log(JSON.stringify(escrows));
}

if (argv?.mode == "release_escrow" && argv?.escrow_id) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    let escrowContract = new ContractPromise(
        api,
        escrowContractMetadata,
        registry.contracts.escrow
    );

    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"]);
    const transactionHash = await escrowRelease(api, escrowContract, marketplaceWallet, argv?.escrow_id)
    
    console.log(JSON.stringify({
        address: marketplaceWallet.address,
        hash: transactionHash
    }))
}

if (argv?.mode == "cancel_escrow" && argv?.escrow_id) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });

    let escrowContract = new ContractPromise(
        api,
        escrowContractMetadata,
        registry.contracts.escrow
    );

    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"]);
    const transactionHash = await escrowCancel(api, escrowContract, marketplaceWallet, argv?.escrow_id)

    console.log(JSON.stringify({
        address: marketplaceWallet.address,
        hash: transactionHash
    }))
}

if (argv?.mode == "mint_escrow" && argv?.seller_address && argv?.value && argv?.user_wallet_mnemonic) {

    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });

    let escrowContract = new ContractPromise(
        api,
        escrowContractMetadata,
        registry.contracts.escrow
    );
    const userWallet = keyring.addFromUri(argv?.user_wallet_mnemonic, { name: 'user pair' }, 'ed25519')

    let [transactionHash, escrowId] = await mintEscrow(
        api,
        escrowContract, 
        settings["polkadot_escrow_contract_marketplace_id"], 
        argv?.seller_address,
        argv?.cid || "",
        null, null, 
        argv?.value, 
        userWallet,
    );

    console.log(JSON.stringify({
        hash: transactionHash,
        address: userWallet.address,
        escrow_id: escrowId,
    }))

}

process.exit()