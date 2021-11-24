// ExpressJS Setup
const express = require('express');
const app = express();

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Constants
const PORT = 8080;
const HOST = '0.0.0.0';

// use static file
app.use(express.static(path.join(__dirname, 'views')));

// configure app to use body-parser
app.use(express.json());
app.use(express.urlencoded({ extended: false }));

// main page routing
app.get('/', (req, res)=>{
    res.sendFile(__dirname + '/index.html');
})

async function cc_call(fn_name, args){
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('myart');
    const contract = network.getContract('ArtWork');

    var result;
    
    if(fn_name == 'addUser')
        result = await contract.submitTransaction('addUser', args);
    else if( fn_name == 'addWork')
    {
        UCIcode = args[0];
        title = args[1];
        artist = args[2];
        initialPrice = args[3];
        properties = args[4];

        result = await contract.submitTransaction('addWork', UCIcode, title, artist, initialPrice, properties);
    }
    else if (fn_name == 'tradeProps')
    {
        seller = args[0];
        buyer = args[1];
        work = args[2];
        propsNum = args[3];

        result = await contract.submitTransaction('tradeProps', seller, buyer, work, propsNum);
    }
    else if (fn_name == 'endTradeProps')
    {
        work = args[0];
        finalPrice = args[1];

        result = await contract.submitTransaction('endTradeProps', work, finalPrice);
    }
    else if(fn_name == 'getHistory')
        result = await contract.evaluateTransaction('getHistory', args);
    else
        result = 'not supported function'

    return result;
}

// add user
app.post('/user', async(req, res)=>{
    const userID = req.body.userID;
    console.log("add user: " + userID);

    result = cc_call('addUser', userID)

    const myobj = {result: "success"}
    res.status(200).json(myobj)
})

// add work
app.post('/work', async(req, res)=>{
    const UCIcode = req.body.UCIcode;
    const title = req.body.title;
    const artist = req.body.artist;
    const initialPrice = req.body.initialPrice;
    const totalProperty = req.body.totalProperty;
    console.log("add work UCIcode: " + UCIcode);
    console.log("add work title: " + title);
    console.log("add work artist: " + artist);
    console.log("add work initialPrice: " + initialPrice);
    console.log("add work totalProperty: " + totalProperty);

    var args=[UCIcode, title, artist, initialPrice, totalProperty];
    result = cc_call('addWork', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// trade Property
app.post('/trade', async(req, res)=>{
    const sellerID = req.body.sellerID;
    const buyerID = req.body.buyerID;
    const workID = req.body.workID;
    const propsNum = req.body.propsNum;
    console.log("start trade with sellerID: " + sellerID);
    console.log("start trade with buyerID: " + buyerID);
    console.log("start trade with workID: " + workID);
    console.log("start trade with propsNum: " + propsNum);

    var args=[sellerID, buyerID, workID, propsNum];
    result = cc_call('tradeProps', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// end trade Property
app.post('/end', async(req, res)=>{
    const workID = req.body.workID;
    const finalPrice = req.body.finalPrice;
    console.log("end trade workID: " + workID);
    console.log("ende trade workID with final price: " + finalPrice);

    var args=[workID, finalPrice];
    result = cc_call('endTradeProps', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// history
app.get('/user', async(req, res)=>{

    try {
        const id = req.query.pid;
        console.log(`${id}`);

        //result = cc_call('history', [id])

        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);
    
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log(`cc_call`);
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('myart');
        const contract = network.getContract('ArtWork');
        
        result = await contract.evaluateTransaction('history', id);
        
        gateway.disconnect();
        
        console.log(`${result}`);
        const myobj = JSON.parse(result)

        
        res.status(200).json(myobj)

    }
    catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        //process.exit(1);
    }
});

// // find mate
// app.post('/mate/:email', async (req,res)=>{
//     const email = req.body.email;
//     console.log("email: " + req.body.email);
//     const walletPath = path.join(process.cwd(), 'wallet');
//     const wallet = new FileSystemWallet(walletPath);
//     console.log(`Wallet path: ${walletPath}`);

//     // Check to see if we've already enrolled the user.
//     const userExists = await wallet.exists('user1');
//     if (!userExists) {
//         console.log('An identity for the user "user1" does not exist in the wallet');
//         console.log('Run the registerUser.js application before retrying');
//         return;
//     }
//     const gateway = new Gateway();
//     await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
//     const network = await gateway.getNetwork('mychannel');
//     const contract = network.getContract('teamate');
//     const result = await contract.evaluateTransaction('readRating', email);
//     const myobj = JSON.parse(result)
//     res.status(200).json(myobj)
//     // res.status(200).json(result)

// });

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
