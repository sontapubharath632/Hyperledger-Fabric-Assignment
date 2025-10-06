const express = require("express");
const bodyParser = require("body-parser");
const { Gateway, Wallets } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const config = require("./config.json");

const app = express();
app.use(bodyParser.json());

async function getContract() {
    const ccpPath = path.resolve(__dirname, config.connectionProfile);
    const ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));
    const wallet = await Wallets.newFileSystemWallet(config.walletPath);

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: config.userId,
        discovery: { enabled: true, asLocalhost: true }
    });

    const network = await gateway.getNetwork(config.channelName);
    return network.getContract(config.chaincodeName);
}

app.post("/asset", async (req, res) => {
    try {
        const contract = await getContract();
        const { dealerId, msisdn, mpin, balance, status } = req.body;
        await contract.submitTransaction("CreateAsset", dealerId, msisdn, mpin, balance, status);
        res.json({ message: `Asset ${dealerId} created successfully` });
    } catch (err) {
        res.status(500).json({ error: err.toString() });
    }
});

app.get("/asset/:id", async (req, res) => {
    try {
        const contract = await getContract();
        const result = await contract.evaluateTransaction("ReadAsset", req.params.id);
        res.json(JSON.parse(result.toString()));
    } catch (err) {
        res.status(500).json({ error: err.toString() });
    }
});

app.put("/asset/:id", async (req, res) => {
    try {
        const contract = await getContract();
        const { balance, status } = req.body;
        await contract.submitTransaction("UpdateAsset", req.params.id, balance, status);
        res.json({ message: `Asset ${req.params.id} updated successfully` });
    } catch (err) {
        res.status(500).json({ error: err.toString() });
    }
});

app.get("/asset/:id/history", async (req, res) => {
    try {
        const contract = await getContract();
        const result = await contract.evaluateTransaction("GetHistory", req.params.id);
        res.json(JSON.parse(result.toString()));
    } catch (err) {
        res.status(500).json({ error: err.toString() });
    }
});

app.listen(8080, () => console.log("REST API running on port 8080"));
