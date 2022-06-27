package abi

const PostageStampABI = `[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_token",
				"type": "address"
			}
		],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "batchId",
				"type": "bytes32"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "totalAmount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "normalisedBalance",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "owner",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint8",
				"name": "depth",
				"type": "uint8"
			},
			{
				"indexed": false,
				"internalType": "uint8",
				"name": "bucketDepth",
				"type": "uint8"
			},
			{
				"indexed": false,
				"internalType": "bool",
				"name": "immutableFlag",
				"type": "bool"
			}
		],
		"name": "BatchCreated",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "batchId",
				"type": "bytes32"
			},
			{
				"indexed": false,
				"internalType": "uint8",
				"name": "newDepth",
				"type": "uint8"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "normalisedBalance",
				"type": "uint256"
			}
		],
		"name": "BatchDepthIncrease",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "batchId",
				"type": "bytes32"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "topupAmount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "normalisedBalance",
				"type": "uint256"
			}
		],
		"name": "BatchTopUp",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "from",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "Incentive",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "previousOwner",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "OwnershipTransferred",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "Paused",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "price",
				"type": "uint256"
			}
		],
		"name": "PriceUpdate",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "previousAdminRole",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "newAdminRole",
				"type": "bytes32"
			}
		],
		"name": "RoleAdminChanged",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "RoleGranted",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "RoleRevoked",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "Unpaused",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "DEFAULT_ADMIN_ROLE",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "PAUSER_ROLE",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "PRICE_ORACLE_ROLE",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"name": "batches",
		"outputs": [
			{
				"internalType": "address",
				"name": "owner",
				"type": "address"
			},
			{
				"internalType": "uint8",
				"name": "depth",
				"type": "uint8"
			},
			{
				"internalType": "bool",
				"name": "immutableFlag",
				"type": "bool"
			},
			{
				"internalType": "uint256",
				"name": "normalisedBalance",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_owner",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_initialBalancePerChunk",
				"type": "uint256"
			},
			{
				"internalType": "uint8",
				"name": "_depth",
				"type": "uint8"
			},
			{
				"internalType": "uint8",
				"name": "_bucketDepth",
				"type": "uint8"
			},
			{
				"internalType": "bytes32",
				"name": "_nonce",
				"type": "bytes32"
			},
			{
				"internalType": "bool",
				"name": "_immutable",
				"type": "bool"
			}
		],
		"name": "createBatch",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "currentTotalOutPayment",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getIncentive",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			}
		],
		"name": "getRoleAdmin",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "grantRole",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "hasRole",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_recipient",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_amount",
				"type": "uint256"
			}
		],
		"name": "incentive",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "_batchId",
				"type": "bytes32"
			},
			{
				"internalType": "uint8",
				"name": "_newDepth",
				"type": "uint8"
			}
		],
		"name": "increaseDepth",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "lastPrice",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "lastUpdatedBlock",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "pause",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "paused",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "_batchId",
				"type": "bytes32"
			}
		],
		"name": "remainingBalance",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "renounceOwnership",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "renounceRole",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "revokeRole",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_price",
				"type": "uint256"
			}
		],
		"name": "setPrice",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes4",
				"name": "interfaceId",
				"type": "bytes4"
			}
		],
		"name": "supportsInterface",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "token",
		"outputs": [
			{
				"internalType": "contract IERC20",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "_batchId",
				"type": "bytes32"
			},
			{
				"internalType": "uint256",
				"name": "_topupAmountPerChunk",
				"type": "uint256"
			}
		],
		"name": "topUp",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "totalOutPayment",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "transferOwnership",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "unPause",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

const PostageStampBin = "0x60806040523480156200001157600080fd5b506040516200362038038062003620833981810160405281019062000037919062000354565b6000600160006101000a81548160ff02191690831515021790555062000072620000666200010160201b60201c565b6200010960201b60201c565b80600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620000c86000801b33620001cc60201b60201c565b620000fa7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a33620001cc60201b60201c565b50620003ce565b600033905090565b600060018054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816001806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b620001de8282620001e260201b60201c565b5050565b620001f48282620002d360201b60201c565b620002cf57600160008084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550620002746200010160201b60201c565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b6000815190506200034e81620003b4565b92915050565b6000602082840312156200036757600080fd5b600062000377848285016200033d565b91505092915050565b60006200038d8262000394565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b620003bf8162000380565b8114620003cb57600080fd5b50565b61324280620003de6000396000f3fe608060405234801561001057600080fd5b50600436106101c45760003560e01c806391d14854116100f9578063d71ba7c411610097578063f2fde38b11610071578063f2fde38b146104be578063f7b188a5146104da578063f90ce5ba146104e4578063fc0c546a14610502576101c4565b8063d71ba7c414610452578063dde5aa6914610482578063e63ab1e9146104a0576101c4565b8063b67644b9116100d3578063b67644b9146103c9578063b998902f146103e5578063c81e25ab14610403578063d547741f14610436576101c4565b806391d148541461035d578063a217fddf1461038d578063b545ebc0146103ab576101c4565b80635239af71116101665780638456cb59116101405780638456cb59146102fd57806388eacbd8146103075780638da5cb5b1461032357806391b7f5ed14610341576101c4565b80635239af71146102b95780635c975abb146102d5578063715018a6146102f3576101c4565b80632f2ff15d116101a25780632f2ff15d1461024757806336568abe1461026357806347aab79b1461027f57806351b17cd01461029b576101c4565b806301ffc9a7146101c9578063053f14da146101f9578063248a9ca314610217575b600080fd5b6101e360048036038101906101de91906121c3565b610520565b6040516101f091906126db565b60405180910390f35b61020161059a565b60405161020e91906129ae565b60405180910390f35b610231600480360381019061022c91906120e6565b6105a0565b60405161023e91906126f6565b60405180910390f35b610261600480360381019061025c919061210f565b6105bf565b005b61027d6004803603810190610278919061210f565b6105e0565b005b61029960048036038101906102949190612187565b610663565b005b6102a36108d7565b6040516102b091906129ae565b60405180910390f35b6102d360048036038101906102ce9190612034565b61091f565b005b6102dd610d1e565b6040516102ea91906126db565b60405180910390f35b6102fb610d35565b005b610305610d49565b005b610321600480360381019061031c9190611ff8565b610dbc565b005b61032b610f73565b60405161033891906125f2565b60405180910390f35b61035b600480360381019061035691906121ec565b610f9b565b005b6103776004803603810190610372919061210f565b611065565b60405161038491906126db565b60405180910390f35b6103956110cf565b6040516103a291906126f6565b60405180910390f35b6103b36110d6565b6040516103c091906129ae565b60405180910390f35b6103e360048036038101906103de919061214b565b6110dc565b005b6103ed61135e565b6040516103fa91906126f6565b60405180910390f35b61041d600480360381019061041891906120e6565b611382565b60405161042d9493929190612696565b60405180910390f35b610450600480360381019061044b919061210f565b6113ec565b005b61046c600480360381019061046791906120e6565b61140d565b60405161047991906129ae565b60405180910390f35b61048a6114df565b60405161049791906129ae565b60405180910390f35b6104a8611526565b6040516104b591906126f6565b60405180910390f35b6104d860048036038101906104d39190611fcf565b61154a565b005b6104e26115ce565b005b6104ec611641565b6040516104f991906129ae565b60405180910390f35b61050a611647565b6040516105179190612711565b60405180910390f35b60007f7965db0b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916148061059357506105928261166d565b5b9050919050565b60055481565b6000806000838152602001908152602001600020600101549050919050565b6105c8826105a0565b6105d1816116d7565b6105db83836116eb565b505050565b6105e86117cb565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610655576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064c9061296e565b60405180910390fd5b61065f82826117d3565b5050565b61066b6118b4565b60006002600084815260200190815260200160002090503373ffffffffffffffffffffffffffffffffffffffff168160000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610714576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161070b906127ee565b60405180910390fd5b8060000160159054906101000a900460ff1615610766576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161075d9061282e565b60405180910390fd5b8060000160149054906101000a900460ff1660ff168260ff16116107bf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107b69061290e565b60405180910390fd5b6107c76108d7565b81600101541161080c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610803906128ee565b60405180910390fd5b60008160000160149054906101000a900460ff168361082b9190612bb8565b905060006108518260ff166001901b6108438761140d565b6118fe90919063ffffffff16565b9050838360000160146101000a81548160ff021916908360ff16021790555061088a8161087c6108d7565b61191490919063ffffffff16565b8360010181905550847faf27998ec15e9d3809edad41aec1b5551d8412e71bd07c91611a0237ead1dc8e8585600101546040516108c8929190612a53565b60405180910390a25050505050565b600080600654436108e89190612b84565b905060006109018260055461192a90919063ffffffff16565b90506109188160045461191490919063ffffffff16565b9250505090565b6109276118b4565b600073ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff161415610997576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161098e9061274e565b60405180910390fd5b60008360ff16141580156109b057508360ff168360ff16105b6109ef576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109e6906127ce565b60405180910390fd5b60003383604051602001610a04929190612644565b604051602081830303815290604052805190602001209050600073ffffffffffffffffffffffffffffffffffffffff166002600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610ac1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ab89061286e565b60405180910390fd5b6000610add8660ff166001901b8861192a90919063ffffffff16565b9050600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330846040518463ffffffff1660e01b8152600401610b3e9392919061260d565b602060405180830381600087803b158015610b5857600080fd5b505af1158015610b6c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b9091906120bd565b610bcf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bc69061294e565b60405180910390fd5b6000610beb88610bdd6108d7565b61191490919063ffffffff16565b905060405180608001604052808a73ffffffffffffffffffffffffffffffffffffffff1681526020018860ff1681526020018515158152602001828152506002600085815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548160ff021916908360ff16021790555060408201518160000160156101000a81548160ff02191690831515021790555060608201518160010155905050827f9b088e2c89b322a3c1d81515e1c88db3d386d022926f0e2d0b9b5813b7413d5883838c8b8b8a604051610d0b969594939291906129f2565b60405180910390a2505050505050505050565b6000600160009054906101000a900460ff16905090565b610d3d611940565b610d4760006119be565b565b610d737f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a33611065565b610db2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610da99061298e565b60405180910390fd5b610dba611a81565b565b610dc4611940565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff1660e01b8152600401610e2192919061266d565b602060405180830381600087803b158015610e3b57600080fd5b505af1158015610e4f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e7391906120bd565b610eb2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ea99061284e565b60405180910390fd5b80600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610f019190612aa3565b925050819055508060086000828254610f1a9190612aa3565b925050819055508173ffffffffffffffffffffffffffffffffffffffff167fd427e26a570fafcb4e8c2c61fde4ef99010612127e4bb6d5f5972eeb12e9f50882604051610f6791906129ae565b60405180910390a25050565b600060018054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b610fc57fdd24a0f121e5ab7c3e97c63eaaf859e0b46792c3e0edfd86e2b3ad50f63011d833611065565b611004576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ffb906128ae565b60405180910390fd5b60006005541461101d576110166108d7565b6004819055505b80600581905550436006819055507fae46785019700e30375a5d7b4f91e32f8060ef085111f896ebf889450aa2ab5a8160405161105a91906129ae565b60405180910390a150565b600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b6000801b81565b60045481565b6110e46118b4565b6000600260008481526020019081526020016000209050600073ffffffffffffffffffffffffffffffffffffffff168160000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561118f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111869061292e565b60405180910390fd5b6111976108d7565b8160010154116111dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111d3906128ee565b60405180910390fd5b60006112098260000160149054906101000a900460ff1660ff166001901b8461192a90919063ffffffff16565b9050600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330846040518463ffffffff1660e01b815260040161126a9392919061260d565b602060405180830381600087803b15801561128457600080fd5b505af1158015611298573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112bc91906120bd565b6112fb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112f29061294e565b60405180910390fd5b61131283836001015461191490919063ffffffff16565b8260010181905550837faf5756c62d6c0722ef9be1f82bef97ab06ea5aea7f3eb8ad348422079f01d88d8284600101546040516113509291906129c9565b60405180910390a250505050565b7fdd24a0f121e5ab7c3e97c63eaaf859e0b46792c3e0edfd86e2b3ad50f63011d881565b60026020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060000160149054906101000a900460ff16908060000160159054906101000a900460ff16908060010154905084565b6113f5826105a0565b6113fe816116d7565b61140883836117d3565b505050565b600080600260008481526020019081526020016000209050600073ffffffffffffffffffffffffffffffffffffffff168160000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156114b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114b09061292e565b60405180910390fd5b6114d76114c46108d7565b8260010154611ae390919063ffffffff16565b915050919050565b6000600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905090565b7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a81565b611552611940565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156115c2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016115b9906127ae565b60405180910390fd5b6115cb816119be565b50565b6115f87f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a33611065565b611637576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161162e906128ce565b60405180910390fd5b61163f611af9565b565b60065481565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60007f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b6116e8816116e36117cb565b611b5c565b50565b6116f58282611065565b6117c757600160008084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555061176c6117cb565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b600033905090565b6117dd8282611065565b156118b057600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506118556117cb565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a45b5050565b6118bc610d1e565b156118fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016118f39061280e565b60405180910390fd5b565b6000818361190c9190612af9565b905092915050565b600081836119229190612aa3565b905092915050565b600081836119389190612b2a565b905092915050565b6119486117cb565b73ffffffffffffffffffffffffffffffffffffffff16611966610f73565b73ffffffffffffffffffffffffffffffffffffffff16146119bc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119b39061288e565b60405180910390fd5b565b600060018054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816001806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b611a896118b4565b60018060006101000a81548160ff0219169083151502179055507f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258611acc6117cb565b604051611ad991906125f2565b60405180910390a1565b60008183611af19190612b84565b905092915050565b611b01611bf9565b6000600160006101000a81548160ff0219169083151502179055507f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa611b456117cb565b604051611b5291906125f2565b60405180910390a1565b611b668282611065565b611bf557611b8b8173ffffffffffffffffffffffffffffffffffffffff166014611c42565b611b998360001c6020611c42565b604051602001611baa9291906125b8565b6040516020818303038152906040526040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611bec919061272c565b60405180910390fd5b5050565b611c01610d1e565b611c40576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611c379061278e565b60405180910390fd5b565b606060006002836002611c559190612b2a565b611c5f9190612aa3565b67ffffffffffffffff811115611c9e577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280601f01601f191660200182016040528015611cd05781602001600182028036833780820191505090505b5090507f300000000000000000000000000000000000000000000000000000000000000081600081518110611d2e577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f780000000000000000000000000000000000000000000000000000000000000081600181518110611db8577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060006001846002611df89190612b2a565b611e029190612aa3565b90505b6001811115611eee577f3031323334353637383961626364656600000000000000000000000000000000600f861660108110611e6a577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b1a60f81b828281518110611ea7577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600485901c945080611ee790612cce565b9050611e05565b5060008414611f32576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611f299061276e565b60405180910390fd5b8091505092915050565b600081359050611f4b81613182565b92915050565b600081359050611f6081613199565b92915050565b600081519050611f7581613199565b92915050565b600081359050611f8a816131b0565b92915050565b600081359050611f9f816131c7565b92915050565b600081359050611fb4816131de565b92915050565b600081359050611fc9816131f5565b92915050565b600060208284031215611fe157600080fd5b6000611fef84828501611f3c565b91505092915050565b6000806040838503121561200b57600080fd5b600061201985828601611f3c565b925050602061202a85828601611fa5565b9150509250929050565b60008060008060008060c0878903121561204d57600080fd5b600061205b89828a01611f3c565b965050602061206c89828a01611fa5565b955050604061207d89828a01611fba565b945050606061208e89828a01611fba565b935050608061209f89828a01611f7b565b92505060a06120b089828a01611f51565b9150509295509295509295565b6000602082840312156120cf57600080fd5b60006120dd84828501611f66565b91505092915050565b6000602082840312156120f857600080fd5b600061210684828501611f7b565b91505092915050565b6000806040838503121561212257600080fd5b600061213085828601611f7b565b925050602061214185828601611f3c565b9150509250929050565b6000806040838503121561215e57600080fd5b600061216c85828601611f7b565b925050602061217d85828601611fa5565b9150509250929050565b6000806040838503121561219a57600080fd5b60006121a885828601611f7b565b92505060206121b985828601611fba565b9150509250929050565b6000602082840312156121d557600080fd5b60006121e384828501611f90565b91505092915050565b6000602082840312156121fe57600080fd5b600061220c84828501611fa5565b91505092915050565b61221e81612bec565b82525050565b61222d81612bfe565b82525050565b61223c81612c0a565b82525050565b61224b81612c77565b82525050565b600061225c82612a7c565b6122668185612a87565b9350612276818560208601612c9b565b61227f81612d56565b840191505092915050565b600061229582612a7c565b61229f8185612a98565b93506122af818560208601612c9b565b80840191505092915050565b60006122c8602083612a87565b91506122d382612d67565b602082019050919050565b60006122eb602083612a87565b91506122f682612d90565b602082019050919050565b600061230e601483612a87565b915061231982612db9565b602082019050919050565b6000612331602683612a87565b915061233c82612de2565b604082019050919050565b6000612354601483612a87565b915061235f82612e31565b602082019050919050565b6000612377600f83612a87565b915061238282612e5a565b602082019050919050565b600061239a601083612a87565b91506123a582612e83565b602082019050919050565b60006123bd601283612a87565b91506123c882612eac565b602082019050919050565b60006123e0601983612a87565b91506123eb82612ed5565b602082019050919050565b6000612403601483612a87565b915061240e82612efe565b602082019050919050565b6000612426602083612a87565b915061243182612f27565b602082019050919050565b6000612449602383612a87565b915061245482612f50565b604082019050919050565b600061246c602483612a87565b915061247782612f9f565b604082019050919050565b600061248f601583612a87565b915061249a82612fee565b602082019050919050565b60006124b2601483612a87565b91506124bd82613017565b602082019050919050565b60006124d5601783612a98565b91506124e082613040565b601782019050919050565b60006124f8601483612a87565b915061250382613069565b602082019050919050565b600061251b600f83612a87565b915061252682613092565b602082019050919050565b600061253e601183612a98565b9150612549826130bb565b601182019050919050565b6000612561602f83612a87565b915061256c826130e4565b604082019050919050565b6000612584602283612a87565b915061258f82613133565b604082019050919050565b6125a381612c60565b82525050565b6125b281612c6a565b82525050565b60006125c3826124c8565b91506125cf828561228a565b91506125da82612531565b91506125e6828461228a565b91508190509392505050565b60006020820190506126076000830184612215565b92915050565b60006060820190506126226000830186612215565b61262f6020830185612215565b61263c604083018461259a565b949350505050565b60006040820190506126596000830185612215565b6126666020830184612233565b9392505050565b60006040820190506126826000830185612215565b61268f602083018461259a565b9392505050565b60006080820190506126ab6000830187612215565b6126b860208301866125a9565b6126c56040830185612224565b6126d2606083018461259a565b95945050505050565b60006020820190506126f06000830184612224565b92915050565b600060208201905061270b6000830184612233565b92915050565b60006020820190506127266000830184612242565b92915050565b600060208201905081810360008301526127468184612251565b905092915050565b60006020820190508181036000830152612767816122bb565b9050919050565b60006020820190508181036000830152612787816122de565b9050919050565b600060208201905081810360008301526127a781612301565b9050919050565b600060208201905081810360008301526127c781612324565b9050919050565b600060208201905081810360008301526127e781612347565b9050919050565b600060208201905081810360008301526128078161236a565b9050919050565b600060208201905081810360008301526128278161238d565b9050919050565b60006020820190508181036000830152612847816123b0565b9050919050565b60006020820190508181036000830152612867816123d3565b9050919050565b60006020820190508181036000830152612887816123f6565b9050919050565b600060208201905081810360008301526128a781612419565b9050919050565b600060208201905081810360008301526128c78161243c565b9050919050565b600060208201905081810360008301526128e78161245f565b9050919050565b6000602082019050818103600083015261290781612482565b9050919050565b60006020820190508181036000830152612927816124a5565b9050919050565b60006020820190508181036000830152612947816124eb565b9050919050565b600060208201905081810360008301526129678161250e565b9050919050565b6000602082019050818103600083015261298781612554565b9050919050565b600060208201905081810360008301526129a781612577565b9050919050565b60006020820190506129c3600083018461259a565b92915050565b60006040820190506129de600083018561259a565b6129eb602083018461259a565b9392505050565b600060c082019050612a07600083018961259a565b612a14602083018861259a565b612a216040830187612215565b612a2e60608301866125a9565b612a3b60808301856125a9565b612a4860a0830184612224565b979650505050505050565b6000604082019050612a6860008301856125a9565b612a75602083018461259a565b9392505050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b6000612aae82612c60565b9150612ab983612c60565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115612aee57612aed612cf8565b5b828201905092915050565b6000612b0482612c60565b9150612b0f83612c60565b925082612b1f57612b1e612d27565b5b828204905092915050565b6000612b3582612c60565b9150612b4083612c60565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615612b7957612b78612cf8565b5b828202905092915050565b6000612b8f82612c60565b9150612b9a83612c60565b925082821015612bad57612bac612cf8565b5b828203905092915050565b6000612bc382612c6a565b9150612bce83612c6a565b925082821015612be157612be0612cf8565b5b828203905092915050565b6000612bf782612c40565b9050919050565b60008115159050919050565b6000819050919050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b6000612c8282612c89565b9050919050565b6000612c9482612c40565b9050919050565b60005b83811015612cb9578082015181840152602081019050612c9e565b83811115612cc8576000848401525b50505050565b6000612cd982612c60565b91506000821415612ced57612cec612cf8565b5b600182039050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000601f19601f8301169050919050565b7f6f776e65722063616e6e6f7420626520746865207a65726f2061646472657373600082015250565b7f537472696e67733a20686578206c656e67746820696e73756666696369656e74600082015250565b7f5061757361626c653a206e6f7420706175736564000000000000000000000000600082015250565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b7f696e76616c6964206275636b6574206465707468000000000000000000000000600082015250565b7f6e6f74206261746368206f776e65720000000000000000000000000000000000600082015250565b7f5061757361626c653a2070617573656400000000000000000000000000000000600082015250565b7f626174636820697320696d6d757461626c650000000000000000000000000000600082015250565b7f696e63656e74697665207472616e73666572206661696c656400000000000000600082015250565b7f626174636820616c726561647920657869737473000000000000000000000000600082015250565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b7f6f6e6c79207072696365206f7261636c652063616e207365742074686520707260008201527f6963650000000000000000000000000000000000000000000000000000000000602082015250565b7f6f6e6c79207061757365722063616e20756e70617573652074686520636f6e7460008201527f7261637400000000000000000000000000000000000000000000000000000000602082015250565b7f626174636820616c726561647920657870697265640000000000000000000000600082015250565b7f6465707468206e6f7420696e6372656173696e67000000000000000000000000600082015250565b7f416363657373436f6e74726f6c3a206163636f756e7420000000000000000000600082015250565b7f626174636820646f6573206e6f74206578697374000000000000000000000000600082015250565b7f6661696c6564207472616e736665720000000000000000000000000000000000600082015250565b7f206973206d697373696e6720726f6c6520000000000000000000000000000000600082015250565b7f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560008201527f20726f6c657320666f722073656c660000000000000000000000000000000000602082015250565b7f6f6e6c79207061757365722063616e2070617573652074686520636f6e74726160008201527f6374000000000000000000000000000000000000000000000000000000000000602082015250565b61318b81612bec565b811461319657600080fd5b50565b6131a281612bfe565b81146131ad57600080fd5b50565b6131b981612c0a565b81146131c457600080fd5b50565b6131d081612c14565b81146131db57600080fd5b50565b6131e781612c60565b81146131f257600080fd5b50565b6131fe81612c6a565b811461320957600080fd5b5056fea2646970667358221220a253eaf7f3697053bb911d55f978ac28e4fa9e35d139eeba4d4635763502295164736f6c63430008020033"
