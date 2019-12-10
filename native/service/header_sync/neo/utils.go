/*

 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package neo

//func PutBlockHeader(native *native.NativeService, chainID uint64, blockHeader *neorpc.BlockHeader) error {
//	contract := utils.HeaderSyncContractAddress
//	headerBytes := blockHeader.ToBytes()
//	heightBytes := utils.GetUint32Bytes(blockHeader.Index)
//	chainIDBytes := utils.GetUint64Bytes(chainID)
//
//	blockHash := blockHeader.Hash
//	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.BLOCK_HEADER), chainIDBytes, blockHash.Bytes()),
//		cstates.GenRawStorageItem(headerBytes))
//	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.HEADER_INDEX), chainIDBytes, heightBytes),
//		cstates.GenRawStorageItem(blockHash.Bytes()))
//	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.CURRENT_HEADER_HEIGHT), chainIDBytes),
//		cstates.GenRawStorageItem(heightBytes))
//	notifyPutHeader(native, chainID, blockHeader.Index, blockHash.String())
//	return nil
//}

//func GetHeaderByHeight(native *native.NativeService, chainID uint64, height uint32) (*neorpc.BlockHeader, error) {
//	contract := utils.HeaderSyncContractAddress
//	heightBytes := utils.GetUint32Bytes(height)
//	chainIDBytes := utils.GetUint64Bytes(chainID)
//
//	blockHashStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(hscommon.HEADER_INDEX),
//		chainIDBytes, heightBytes))
//	if err != nil {
//		return nil, fmt.Errorf("GetHeaderByHeight, get blockHashStore error: %v", err)
//	}
//	if blockHashStore == nil {
//		return nil, fmt.Errorf("GetHeaderByHeight, can not find any index records")
//	}
//	blockHashBytes, err := cstates.GetValueFromRawStorageItem(blockHashStore)
//	if err != nil {
//		return nil, fmt.Errorf("GetHeaderByHeight, deserialize blockHashBytes from raw storage item err:%v", err)
//	}
//	header := &neorpc.BlockHeader{}
//	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(hscommon.BLOCK_HEADER),
//		chainIDBytes, blockHashBytes))
//	if err != nil {
//		return nil, fmt.Errorf("GetHeaderByHeight, get headerStore error: %v", err)
//	}
//	if headerStore == nil {
//		return nil, fmt.Errorf("GetHeaderByHeight, can not find any header records")
//	}
//	headerBytes, err := cstates.GetValueFromRawStorageItem(headerStore)
//	if err != nil {
//		return nil, fmt.Errorf("GetHeaderByHeight, deserialize headerBytes from raw storage item err:%v", err)
//	}
//	header.ToBlockHeader(headerBytes)
//	return header, nil
//}
//
//func GetHeaderByHash(native *native.NativeService, chainID uint64, hash common.Uint256) (*neorpc.BlockHeader, error) {
//	contract := utils.HeaderSyncContractAddress
//	chainIDBytes := utils.GetUint64Bytes(chainID)
//
//	header := &neorpc.BlockHeader{}
//	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(hscommon.BLOCK_HEADER), chainIDBytes, hash.ToArray()))
//	if err != nil {
//		return nil, fmt.Errorf("GetHeaderByHash, get headerStore error: %v", err)
//	}
//	if headerStore == nil {
//		return nil, fmt.Errorf("GetHeaderByHash, can not find any records")
//	}
//	headerBytes, err := cstates.GetValueFromRawStorageItem(headerStore)
//	if err != nil {
//		return nil, fmt.Errorf("GetHeaderByHash, deserialize from raw storage item err:%v", err)
//	}
//	header.ToBlockHeader(headerBytes)
//	return header, nil
//}
//
//func notifyPutHeader(native *native.NativeService, chainID uint64, height uint32, blockHash string) {
//	if !config.DefConfig.Common.EnableEventLog {
//		return
//	}
//	native.AddNotify(
//		&event.NotifyEventInfo{
//			ContractAddress: utils.HeaderSyncContractAddress,
//			States:          []interface{}{chainID, height, blockHash, native.GetHeight()},
//		})
//}
