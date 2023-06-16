# Test network - Fabric native setup

# Prereqs

- Follow the Fabric documentation for the [Prereqs](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html)
- Follow the Fabric documentation for [downloading the Fabric samples and binaries](https://hyperledger-fabric.readthedocs.io/en/latest/install.html). You can skip the docker image downloads by using `./install-fabric.sh binary samples`

# Instructions for starting network
## Running each component separately

Open terminal windows for 3 ordering nodes or 4 if running BFT Consensus, 4 peer nodes, and 4 peer admins as seen in the following terminal setup. The first two peers and peer admins belong to Org1, the latter two peer and peer admins belong to Org2.
Note, you can start with two ordering nodes and a single Org1 peer node and single Org1 peer admin terminal if you would like to keep things even more minimal (two ordering nodes are required to achieve Raft consensus (2 of 3), while a single peer from Org1 can be utilized since the endorsement policy is set as any single organization).
![Terminal setup](terminal_setup.png)

The following instructions will have you run simple bash scripts that set environment variable overrides for a component and then runs the component.
The scripts contain only simple single-line commands so that they are easy to read and understand.
If you have trouble running bash scripts in your environment, you can just as easily copy and paste the individual commands from the script files instead of running the script files.

- cd to the `test-network-nano-bash` directory in each terminal window
- In the first orderer terminal, run `./generate_artifacts.sh` to generate crypto material (calls cryptogen) and application channel genesis block and configuration transactions (calls configtxgen). The artifacts will be created in the `crypto-config` and `channel-artifacts` directories. If you are running BFT consensus then run `./generate_artifacts.sh BFT`.
- In the three orderer terminals, run `./orderer1.sh`, `./orderer2.sh`, `./orderer3.sh` respectively. If you are running BFT consensus then run `./orderer4.sh` in the fourth orderer terminal also.
- In the four peer terminals, run `./peer1.sh`, `./peer2.sh`, `./peer3.sh`, `./peer4.sh` respectively.
- Note that each orderer and peer write their data (including their ledgers) to their own subdirectory under the `data` directory
- Open a different terminal and run `./join_orderers.sh`. If you are running BFT Consensus then run `./join_orderers.sh BFT` instead.
- In the four peer admin terminals, run `source peer1admin.sh && ./join_channel.sh`, `source peer2admin.sh && ./join_channel.sh`, `source peer3admin.sh && ./join_channel.sh`, `source peer4admin.sh && ./join_channel.sh` respectively.

Note the syntax of running the scripts. The peer admin scripts set the admin environment variables and must be run with the `source` command in order that the exported environment variables can be utilized by any subsequent user commands.

## Starting the network with one command

Using the individual scripts above gives you more control of the process of starting a Fabric network and demonstrates how all the required components fit together, however the same network can also be started using a single script for convenience.

For Raft consensus type:
```shell
./network.sh start
```

For BFT consensus type:
```shell
./network.sh start -o BFT
```

After the network has started, use seperate terminals to run peer commands.
You will need to configure the peer environment for each new terminal.
For example to run against peer1, use:

```shell
source peer1admin.sh
```
## 1. Using a chaincode container

Package and install the chaincode on peer1:

```shell
peer lifecycle chaincode package chaincode1.tar.gz --path ./chaincode-go --lang golang --label chaincode1_1

peer lifecycle chaincode install chaincode1.tar.gz
```

The chaincode install may take a minute since the `fabric-ccenv` chaincode builder docker image will be downloaded if not already available on your machine.

Copy the returned chaincode package ID into a `CHAINCODE_ID` environment variable for use in subsequent commands, or better yet use the peer `calculatepackageid` command to set the environment variable:

```shell
export CHAINCODE_ID=$(peer lifecycle chaincode calculatepackageid basic.tar.gz) && echo $CHAINCODE_ID
```

## Activate the chaincode

Using the peer1 admin, approve and commit the chaincode (only a single approver is required based on the lifecycle endorsement policy of any organization):

```shell
peer lifecycle chaincode approveformyorg -o 127.0.0.1:6050 --channelID mychannel --name basic --version 1 --package-id $CHAINCODE_ID --sequence 1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt

peer lifecycle chaincode commit -o 127.0.0.1:6050 --channelID mychannel --name basic --version 1 --sequence 1 --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

**Note:** after following the instructions above, the chaincode will only be installed on peer1 and will only be available in the peer1admin shell.
Rerun the `peer lifecycle chaincode install` command in other peer admin shells to install it on the corresponding peer.
You will also need to rerun the `peer lifecycle chaincode approveformyorg` command to use the chaincode on peers in another organisation, e.g. using the peer3admin shell.

## Interact with the chaincode

For your convenience you can run `chaincode_interaction.sh` from peer1admin terminal to make this simple transaction. The ouput of the script is redirected to the logs folder.\
Congratulations, you have deployed a minimal Fabric network! Inspect the scripts if you would like to see the minimal set of commands that were required to deploy the network.

# Stopping the network

If you started the Fabric componentes individually, utilize `Ctrl-C` in the orderer and peer terminal windows to kill the orderer and peer processes. You can run the scripts again to restart the components with their existing data, or run `./generate_artifacts` again to clean up the existing artifacts and data if you would like to restart with a clean environment.

If you used the `network.sh` script, utilize `Ctrl-C` to kill the orderer and peer processes. You can restart the network with the existing data, or run `./network.sh clean` to remove old data before restarting.
