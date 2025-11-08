import { ApiPromise, WsProvider } from "@polkadot/api";
import { ContractPromise } from "@polkadot/api-contract";
import { Keyring } from '@polkadot/keyring';

import BN from "bn.js";
import minimist from "minimist";

import crowdloanLendContractMetadata from "./polkadot_contracts/crowdloan_lend_contract.json" with { type: "json" };
import crowdloanContractMetadata from "./polkadot_contracts/crowdloan_contract.json" with { type: "json" };
import registry from "./polkadot_contracts/registry.json" with { type: "json" };
import settings from "./../settings.json" with { type: "json" };
import BigNumber from "bignumber.js";


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

function convertBalance(number) {
    let bigNumber = new BigNumber(number.replace(/,/g, ''));
    let dividedResult = bigNumber.dividedBy(new BigNumber(1e10));
    return dividedResult.toNumber();
}

const keyring = new Keyring({ type: 'sr25519' });
const argv = minimist(process.argv.slice(2))

async function mintCrowdloan(api, crowdloanContract, marketplaceId, goal, endDate, weeklyInterest, userWallet) {
    let res;

    await new Promise((resolve, reject) => {
        crowdloanContract.tx["crowdloan::mint"](
            {
                gasLimit: buildGasLimit(api),
            },
            { u16: marketplaceId }, goal, endDate, weeklyInterest,
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

async function crowdloanStatus(api, crowdloanContract, crowdloanId, from) {
    let { output } = await crowdloanContract.query["crowdloan::status"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: crowdloanId },
    );

    if (output?.isOk) {
        return output.asOk.toHuman();
    } else {
        return ""
    }
}

async function crowdloanGoalDate(api, crowdloanContract, crowdloanId, from) {
    let { output } = await crowdloanContract.query["crowdloan::goalDate"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: crowdloanId },
    );

    if (output?.isOk) {
        return parseInt(output.asOk.toHuman().replace(/,/g, '')) / 1000;
    } else {
        return null;
    }
}

async function crowdloanGoalValue(api, crowdloanContract, crowdloanId, from) {
    let { output } = await crowdloanContract.query["crowdloan::goalValue"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: crowdloanId },
    );

    if (output?.isOk) {
        return convertBalance(output.asOk.toHuman());
    } else {
        return null;
    }
}

async function crowdloanClosingAmount(api, crowdloanContract, crowdloanId, from) {
    let { output } = await crowdloanContract.query["crowdloan::closingAmount"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: crowdloanId },
    );

    if (output?.isOk) {
        return convertBalance(output.asOk.toHuman().Ok);
    } else {
        return null;
    }
}

async function crowdloanNumberOfLenders(api, crowdloanContract, crowdloanId, from) {
    let { output } = await crowdloanContract.query["crowdloan::numberOfLenders"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: crowdloanId },
    );

    if (output?.isOk) {
        return parseInt(output.asOk.toHuman().replace(/,/g, ''));
    } else {
        return null;
    }
}

async function crowdloanWeeklyInterest(api, crowdloanContract, crowdloanId, from) {
    let { output } = await crowdloanContract.query["crowdloan::weeklyInterest"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: crowdloanId },
    );

    if (output?.isOk) {
        return parseInt(output.asOk.toHuman());
    } else {
        return null;
    }
}

async function crowdloanCloseDate(api, crowdloanContract, crowdloanId, from) {
    let { output } = await crowdloanContract.query["crowdloan::closeDate"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: crowdloanId },
    );

    if (output?.isOk) {
        return parseInt(output.asOk.toHuman()) / 1000;
    } else {
        return null;
    }
}

async function crowdloanInfo(api, crowdloanContract, crowdloanId, from) {
    let cs = await crowdloanStatus(api, crowdloanContract, crowdloanId, from);
    let owner = await psp34OwnerOf(api, crowdloanContract, crowdloanId, from);
    let goalDate = await crowdloanGoalDate(api, crowdloanContract, crowdloanId, from);
    let goalValue = await crowdloanGoalValue(api, crowdloanContract, crowdloanId, from);
    let value = await crowdloanValue(api, crowdloanContract, crowdloanId, from);
    let weeklyInterest = await crowdloanWeeklyInterest(api, crowdloanContract, crowdloanId, from);
    let closeDate = await crowdloanCloseDate(api, crowdloanContract, crowdloanId, from);
    let closingAmount = await crowdloanClosingAmount(api, crowdloanContract, crowdloanId, from);
    let numberOfLenders = await crowdloanNumberOfLenders(api, crowdloanContract, crowdloanId, from);

    return {
        id: crowdloanId,
        owner: owner,
        status: {
            status: cs[0],
            timestamp: parseInt(cs[1].replace(/,/g, '')),
        },
        goal_date: goalDate,
        goal_value: goalValue,
        value: value,
        weekly_interest: weeklyInterest,
        close_date: closeDate,
        closing_amount: closingAmount,
        n_lenders: numberOfLenders,
    }
}

async function crowdloanLendCrowdloanId(api, crowdloanLendContract, lendId, from) {

    let { output } = await crowdloanLendContract.query["crowdloanLend::crowdloanId"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: lendId },
    );

    if (output?.isOk) {
        return parseInt(output.asOk.toHuman().U64);
    } else {
        return null;
    }
}

async function crowdloandLendValue(api, crowdloanLendContract, lendId, from) {

    let { output } = await crowdloanLendContract.query["crowdloanLend::lendValue"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: lendId },
    );

    if (output?.isOk) {
        return convertBalance(output.asOk.toHuman());
    } else {
        return null;
    }
}

async function crowdloanValue(api, crowdloanContract, crowdloanId, from) {
    let { output } = await crowdloanContract.query["crowdloan::value"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: crowdloanId },
    );

    if (output?.isOk) {   
        return convertBalance(output.asOk.toHuman());
    } else {
        return null;
    }
}

async function crowdloandLendIndex(api, crowdloanLendContract, lendId, from) {
    let { output } = await crowdloanLendContract.query["crowdloanLend::index"](
        from,
        { gasLimit: buildGasLimit(api) },
        { u64: lendId },
    );

    if (output?.isOk) {
        return parseInt(output.asOk.toHuman());
    } else {
        return null;
    }
}

async function crowdloanLendInfo(api, crowdloanLendContract, lendId, from) {

    let crowdloandId = await crowdloanLendCrowdloanId(api, crowdloanLendContract, lendId, from);
    let lendValue = await crowdloandLendValue(api, crowdloanLendContract, lendId, from);
    let index = await crowdloandLendIndex(api, crowdloanLendContract, lendId, from);

    
    return {
        id: lendId,
        crowdloan_id: crowdloandId,
        lend_value: lendValue,
        index: index,
    }
}

async function crowdloanFund(api, escrowContract, crowdloanId, value, userWallet) {
    let resultHash;

    await new Promise((resolve, reject) => {
        escrowContract.tx["crowdloan::fund"](
            {
                gasLimit: buildGasLimit(api),
                value: value,
            },
            { u64: crowdloanId }
        ).signAndSend(userWallet, (result) => {
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

async function crowdloanWithdraw(api, escrowContract, crowdloanId, userWallet) {
    let resultHash;

    await new Promise((resolve, reject) => {
        escrowContract.tx["crowdloan::withdraw"](
            {
                gasLimit: buildGasLimit(api),
            },
            { u64: crowdloanId }
        ).signAndSend(userWallet, (result) => {
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

async function crowdloanPayback(api, escrowContract, crowdloanId, value, userWallet) {
    let resultHash;

    await new Promise((resolve, reject) => {
        escrowContract.tx["crowdloan::payback"](
            {
                gasLimit: buildGasLimit(api),
                value: value,
            },
            { u64: crowdloanId }
        ).signAndSend(userWallet, (result) => {
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

async function crowdloanPayout(api, escrowContract, lendId, userWallet) {
    let resultHash;

    await new Promise((resolve, reject) => {
        escrowContract.tx["crowdloan::payout"](
            {
                gasLimit: buildGasLimit(api),
            },
            { u64: lendId }
        ).signAndSend(userWallet, (result) => {
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

// CLI Methods

if (argv?.mode == "crowdloan_info" && argv?.crowdloan_id) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])
    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );

    let ei = await crowdloanInfo(api, crowdloanContract, argv?.crowdloan_id, marketplaceWallet.address);
    console.log(JSON.stringify(ei));
}

if (argv?.mode == "mint_crowdloan" && argv?.goal_value && argv?.goal_date && argv?.weekly_interest && argv?.user_wallet_mnemonic) {

    const wsProvider = new WsProvider(registry.endpoint);
    const api = await ApiPromise.create({ provider: wsProvider });

    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );
    const userWallet = keyring.addFromUri(argv?.user_wallet_mnemonic, { name: 'user pair' }, 'ed25519')

    let [transactionHash, crowdloanId] = await mintCrowdloan(
        api,
        crowdloanContract, 
        settings["polkadot_escrow_contract_marketplace_id"], 
        argv?.goal_value,
        argv?.goal_date,
        argv?.weekly_interest,
        userWallet,
    );

    console.log(JSON.stringify({
        hash: transactionHash,
        address: userWallet.address,
        crowdloan_id: crowdloanId,
    }))

}

if (argv?.mode == "crowdloan_total_supply") {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])
    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );
    const totalSupply = await psp34TotalSupply(api, crowdloanContract, marketplaceWallet.address);

    console.log(totalSupply);
}

if (argv?.mode == "crowdloan_owner_of" && argv?.crowdloan_id) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])
    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );
    const ownerOf = await psp34OwnerOf(api, crowdloanContract, argv?.crowdloan_id, marketplaceWallet.address);
    console.log(ownerOf);
}

if (argv?.mode == "list_wallet_crowdloans" && argv?.wallet) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])

    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );

    let numCrowdloans = await psp34BalanceOf(api, crowdloanContract, argv?.wallet, marketplaceWallet.address);
    var lends = []

    for (var i = 0; i < numCrowdloans; i++){
        var lendId = await psp34OwnersTokenByIndex(
            api, 
            crowdloanContract, 
            argv?.wallet, 
            i, 
            marketplaceWallet.address,
        );
        let escrow = await crowdloanInfo(api, crowdloanContract, lendId, marketplaceWallet.address);
        lends.push(escrow)
    }

    console.log(JSON.stringify(lends));
}

if (argv?.mode == "list_wallet_lends" && argv?.wallet) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    const marketplaceWallet = keyring.addFromUri(settings["polkadot_phrase"])

    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );

    let crowdloanLendContract = new ContractPromise(
        api,
        crowdloanLendContractMetadata,
        registry.contracts.crowdloan_lend
    );

    let numLends = await psp34BalanceOf(api, crowdloanLendContract, argv?.wallet, marketplaceWallet.address);
    var lends = []

    for (var i = 0; i < numLends; i++){
        var lendId = await psp34OwnersTokenByIndex(
            api, 
            crowdloanLendContract, 
            argv?.wallet, 
            i, 
            marketplaceWallet.address,
        );
        let lend = await crowdloanLendInfo(api, crowdloanLendContract, lendId, marketplaceWallet.address);
        let loan = await crowdloanInfo(api, crowdloanContract, lend.crowdloan_id, marketplaceWallet.address)
        lend["loan"] = loan
        lends.push(lend)
    }

    console.log(JSON.stringify(lends));
}

if (argv?.mode == "fund_crowdloan" && argv?.crowdloan_id && argv?.value && argv?.user_wallet_mnemonic) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );

    const userWallet = keyring.addFromUri(argv?.user_wallet_mnemonic, { name: 'user pair' }, 'ed25519');
    const transactionHash = await crowdloanFund(
        api, 
        crowdloanContract, 
        argv?.crowdloan_id, 
        argv?.value, 
        userWallet,
    )
    
    console.log(JSON.stringify({
        hash: transactionHash,
        address: userWallet.address,
    }))
}

if (argv?.mode == "withdraw_crowdloan" && argv?.crowdloan_id && argv?.user_wallet_mnemonic) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );

    const userWallet = keyring.addFromUri(argv?.user_wallet_mnemonic, { name: 'user pair' }, 'ed25519');
    const transactionHash = await crowdloanWithdraw(
        api, 
        crowdloanContract, 
        argv?.crowdloan_id,
        userWallet,
    )
    
    console.log(JSON.stringify({
        address: userWallet.address,
        hash: transactionHash,
    }))
}

if (argv?.mode == "payback_crowdloan" && argv?.crowdloan_id && argv?.user_wallet_mnemonic) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );

    const userWallet = keyring.addFromUri(argv?.user_wallet_mnemonic, { name: 'user pair' }, 'ed25519');
    
    let closingAmount = await crowdloanClosingAmount(
        api, 
        crowdloanContract, 
        argv?.crowdloan_id, 
        userWallet.address,
    ) * 1e10;

    const transactionHash = await crowdloanPayback(
        api, 
        crowdloanContract, 
        argv?.crowdloan_id, 
        closingAmount, 
        userWallet,
    )
    
    console.log(JSON.stringify({
        address: userWallet.address,
        hash: transactionHash
    }))
}

if (argv?.mode == "payout_crowdloan" && argv?.crowdloan_id && argv?.lend_id && argv?.user_wallet_mnemonic) {
    const wsProvider = new WsProvider(registry["endpoint"]);
    const api = await ApiPromise.create({ provider: wsProvider });
    let crowdloanContract = new ContractPromise(
        api,
        crowdloanContractMetadata,
        registry.contracts.crowdloan
    );

    const userWallet = keyring.addFromUri(argv?.user_wallet_mnemonic, { name: 'user pair' }, 'ed25519');

    const transactionHash = await crowdloanPayout(
        api, 
        crowdloanContract,
        argv?.lend_id,
        userWallet,
    )
    
    console.log(JSON.stringify({
        address: userWallet.address,
        hash: transactionHash
    }))
}


process.exit()
