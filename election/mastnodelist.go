// Copyright 2018 The Matrix Authors
// This file is part of the Matrix library.
//
// The Matrix library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Matrix library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Matrix library. If not, see <http://www.gnu.org/licenses/>.

package election

type ListNodeInfo struct {
	AnnonceRate     uint32 //Declared computing power (TPS)
	Ip              string //External IP
	NodeId          uint32 //NodeID
	MortgageAccount uint64 //Deposits
	Upime           uint64 //uptime: the consensus algorithm provides functions, and the masternode verification process in charge of maintaining
	TxHash          uint64 //Transacction Hash
}

type MasterListInfo struct {
	preHash         uint64
	height          uint32
	masterNum       int
	masterNodeList  [MaxMasterNodeNum]ListNodeInfo
	offlineNodeList []ListNodeInfo
}

func (mn *MasterListInfo) SetPreHash(hashval uint64) {
	mn.preHash = hashval
}

func (mn *MasterListInfo) SetBlockHeight(height uint32) {
	mn.height = height
}

func (mn *MasterListInfo) SetMasterNodeNum(masterNum int) {
	mn.masterNum = masterNum
}

func (mn *MasterListInfo) AddMasterNodeInfo(nodeinfo *ListNodeInfo) bool{
	if mn.masterNum >= (MaxMasterNodeNum-1){
		return false
	}
	mn.masterNodeList[mn.masterNum] = *nodeinfo
	mn.masterNum++

	return true
}

func (mn *MasterListInfo) GetMsaterListLen() int {
	return mn.masterNum
}

func masterNodeStateUpdate(electionList *MasterListInfo, updateList *MasterListInfo) {
	var key float64

	mapNewMasterNodeList := make(map[float64]int)
	for i := 0; i < updateList.masterNum; i++ {
		key = float64(updateList.masterNodeList[i].NodeId) + float64(updateList.masterNodeList[i].TxHash)*float64(1<<32)
		mapNewMasterNodeList[key] = i
	}

	electionList.offlineNodeList = make([]ListNodeInfo, 0)
	for i := 0; i < electionList.masterNum; i++ {
		key = float64(electionList.masterNodeList[i].NodeId) + float64(electionList.masterNodeList[i].TxHash)*float64(1<<32)
		if _, state := mapNewMasterNodeList[key]; !state {
			electionList.offlineNodeList = append(electionList.offlineNodeList, electionList.masterNodeList[i])
		}
	}
}

func (list *MasterListInfo) masterNodeSearch(txid uint64, nodeid uint32) (int, bool) {
	for i := 0; i < list.masterNum; i++ {
		if (list.masterNodeList[i].NodeId == nodeid) && (list.masterNodeList[i].NodeId == nodeid) {
			return i, true
		}
	}
	return 0, false
}

func (nodeinfo *ListNodeInfo) IsMinner(elenet *ElectionNetInfo) bool {
	if elenet.Active == false {
		return false
	}

	if index, state := elenet.MasterList.masterNodeSearch(nodeinfo.TxHash, nodeinfo.NodeId); state {
		if elenet.ListMapTreeTbl[index].treeType == CalcullatorNetFlag {
			tree := elenet.CalcullatorNet[elenet.ListMapTreeTbl[index].regionIndex]
			treeNodeInfo := tree.treeNode[elenet.ListMapTreeTbl[index].treeNodeIndex]
			if (elenet.ListMapTreeTbl[index].treeNodeIndex == 0) && treeNodeInfo.localIdx != NullNodeIdx {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}
