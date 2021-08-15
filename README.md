go-ds-cluster
============
> gathering distributed key-value datastores to become a cluster


<!-- ABOUT THE PROJECT -->
## About The Project

This project is going to implement [go-datastore](https://github.com/ipfs/go-datastore) in a form of cluster.

[IPFS](https://github.com/ipfs/ipfs) is awsome, we like to use ipfs for data sharing. And [IPFS](https://github.com/ipfs/ipfs) using implementation of [go-datastore](https://github.com/ipfs/go-datastore) as base storage.

There are several implementation of [go-datastore](https://github.com/ipfs/go-datastore):
- [go-ds-flatfs](https://github.com/ipfs/go-ds-flatfs)
- [go-ds-badger](https://github.com/ipfs/go-ds-badger)
- [go-ds-leveldb](https://github.com/ipfs/go-ds-leveldb)
- ...

They are mainly focus on storing data within one PC. Single PC has limit on I/O, there lacks a way to take advantage of distributed system with several or more PCs.

We knew [ipfs-cluster](https://github.com/ipfs/ipfs-cluster), it offers a way to manage pin-set within multiple peers. But we thought it is more suitable for manage data backups and recovery. We expect that even only one ipfs peer could also take advantage of distributed storage.

## Architecture

- data node has two basic components
  - datastore offering key-value strorage
  - cluster instance maintains cluster related logic
- any node in the cluster can be used has ipfs datastore, sharding data to data nodes in the cluster
- a key-value table of all data in the datastore has been maintained by every node
- the relation between hash slots and data node also been maintained by every node

## Roadmap

- data sharding and hash slots maintaining 
- communication module based on libp2p between data nodes
- consensus module build up
- hash slots re-allocate and re-balance strategy 
- data migration after hash slots re-allocat or re-balance to support dynamically adding or removing nodes
- authentication and data management



<!-- GETTING STARTED -->
## Getting Started



### Prerequisites



### Installation





<!-- USAGE EXAMPLES -->
## Usage




<!-- ROADMAP -->
## Roadmap





<!-- CONTRIBUTING -->
## Contributing

PRs are welcome!



<!-- LICENSE -->
## License

Distributed under the MIT License. 



<!-- CONTACT -->
## Contact




<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements





