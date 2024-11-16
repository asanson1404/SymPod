import { Brevis, ErrCode, ProofRequest, Prover, ReceiptData, Field } from 'brevis-sdk-typescript';
import { ethers } from 'ethers';
// @ts-ignore
import { sympodAbi } from '../abi/sympodAbi.js';
import dotenv from 'dotenv';

dotenv.config();

async function main() {

    const provider = new ethers.providers.JsonRpcProvider(process.env.SEPOLIA_RPC_URL);
    const sympodContract = new ethers.Contract(process.env.SYMPOD_CONTRACT_ADDRESS!, sympodAbi, provider);

    const prover = new Prover('localhost:33247');
    const brevis = new Brevis('appsdkv3.brevis.network:443');

    // Listen to an event to verify withdrawal credentials
    sympodContract.on("VerificationRequired", async (pubkey, depositDataRoot, sympodAddr, event) => {
        console.log("VerificationRequired event received:");

        try {            
            await processVerification(prover, brevis, pubkey, depositDataRoot, sympodAddr, event.transactionHash);

        } catch (error) {
            console.error("Error processing verification:", error);
        }
    });

    console.log('Listening for VerificationRequired events...');

}

async function processVerification(prover: Prover, brevis: Brevis, pubkey: string, depositDataRoot: string, sympodAddr: string, transactionHash: string) {
    console.log('Processing verification with arguments:');
    console.log('prover:', prover);
    console.log('brevis:', brevis); 
    console.log('pubkey:', pubkey);
    console.log('depositDataRoot:', depositDataRoot);
    console.log('sympodAddr:', sympodAddr);
    console.log('transactionHash:', transactionHash);

    const proofReq = new ProofRequest();

    proofReq.addReceipt(
        new ReceiptData({
            tx_hash: transactionHash,
            fields: [
                // sympod address
                new Field({
                    log_pos: 0,
                    is_topic: true,
                    field_index: 1,
                }),
                // pubkey
                new Field({
                    log_pos: 0,
                    is_topic: false,
                    field_index: 0,
                }),
            ],
        }),
    );

    console.log('proofReq', proofReq);

    const proofRes = await prover.prove(proofReq);
    // error handling
    if (proofRes.has_err) {
        const err = proofRes.err;
        switch (err.code) {
            case ErrCode.ERROR_INVALID_INPUT:
                console.error('invalid receipt/storage/transaction input:', err.msg);
                break;

            case ErrCode.ERROR_INVALID_CUSTOM_INPUT:
                console.error('invalid custom input:', err.msg);
                break;

            case ErrCode.ERROR_FAILED_TO_PROVE:
                console.error('failed to prove:', err.msg);
                break;
        }
        return;
    }
    console.log('proof', proofRes.proof);

    try {
        const brevisRes = await brevis.submit(proofReq, proofRes, 11155111, 11155111, 0, "", process.env.SYMPOD_CONTRACT_ADDRESS!);
        console.log('brevis res', brevisRes);

        await brevis.wait(brevisRes.queryKey, 11155111);
    } catch (err) {
        console.error(err);
    }

}

await main();