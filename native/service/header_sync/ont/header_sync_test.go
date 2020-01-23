package ont

import (
	"encoding/hex"
	"testing"

	"github.com/ontio/multi-chain/account"
	"github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/core/genesis"
	"github.com/ontio/multi-chain/core/store/leveldbstore"
	"github.com/ontio/multi-chain/core/store/overlaydb"
	"github.com/ontio/multi-chain/core/types"
	"github.com/ontio/multi-chain/native"
	scom "github.com/ontio/multi-chain/native/service/header_sync/common"
	"github.com/ontio/multi-chain/native/storage"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/stretchr/testify/assert"
)

var (
	acct          *account.Account = account.NewAccount("")
	getNativeFunc                  = func() *native.NativeService {
		store, _ := leveldbstore.NewMemLevelDBStore()
		cacheDB := storage.NewCacheDB(overlaydb.NewOverlayDB(store))
		service := native.NewNativeService(cacheDB, nil, 0, 200, common.Uint256{}, 0, nil, false, nil)
		return service
	}
	getBtcHanderFunc = func() *ONTHandler {
		return NewONTHandler()
	}
	setBKers = func() {
		genesis.GenesisBookkeepers = []keypair.PublicKey{acct.PublicKey}
	}
)

func init() {
	setBKers()
}

func NewNative(args []byte, tx *types.Transaction, db *storage.CacheDB) *native.NativeService {
	if db == nil {
		store, _ := leveldbstore.NewMemLevelDBStore()
		db = storage.NewCacheDB(overlaydb.NewOverlayDB(store))
	}
	return native.NewNativeService(db, tx, 0, 0, common.Uint256{0}, 0, args, false, nil)
}

func TestSyncGenesisHeader(t *testing.T) {
	genesisHeader, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000038f121e642af587e34202511f911701ad81725a5fcc83eb6e4fc52d3f9480f0e000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fdfb047b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a362c226964223a22303366313039353238396537666464623838326631636233653135386163633163333064396465363036616632316339376261383531383231653862366561353335227d2c7b22696e646578223a352c226964223a22303338626663353062306533663065356466366434353130363930363563626661376162356433383261353833396363653832653063393633656462303236653934227d2c7b22696e646578223a322c226964223a22303335656236353462616436633634303938393462396234323238396134333631343837346337393834626465366230336161663666633164303438366439643435227d2c7b22696e646578223a332c226964223a22303238316431393863306464333733376139633339313931626332643161663764363561343432363161386136346436656637346436336632376366623565643932227d2c7b22696e646578223a312c226964223a22303235336363666434333962323965636130666539306361376336656161316639383537326130353461613264316435366537326164393663343636313037613835227d2c7b22696e646578223a342c226964223a22303233393637626261333036306266386164653036643962616434356430323835336636633632336534643466353264373637656235366466346433363461393966227d2c7b22696e646578223a372c226964223a22303231353836356261616237303630376634613234313361376139626139356162326333633032303264356237373331633638323465656634386538393966633930227d5d2c22706f735f7461626c65223a5b322c372c322c362c342c332c332c342c322c342c342c352c372c322c342c362c362c312c312c332c332c372c352c362c372c312c372c352c322c342c332c362c352c362c312c362c352c312c342c342c322c312c372c322c362c352c342c352c372c372c332c372c312c332c362c312c352c312c362c332c332c372c322c372c342c342c332c352c332c312c322c372c322c372c362c322c342c362c312c352c352c322c362c372c332c342c352c312c352c342c352c322c312c372c332c342c362c312c312c322c332c352c332c362c325d2c224d6178426c6f636b4368616e676556696577223a31303030307d7de0b053244d70b2c5622a4786d4922c2991e1fe880000")

	param := new(scom.SyncGenesisHeaderParam)
	param.ChainID = 3
	param.GenesisHeader = genesisHeader
	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	tx := &types.Transaction{
		SignedAddr: []common.Address{acct.Address},
	}

	native := NewNative(sink.Bytes(), tx, nil)
	ontHandler := NewONTHandler()
	err := ontHandler.SyncGenesisHeader(native)
	assert.NoError(t, err)
}

func TestSyncBlockHeader(t *testing.T) {
	ontHandler := NewONTHandler()
	var native *native.NativeService
	{
		header7152785, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000038f121e642af587e34202511f911701ad81725a5fcc83eb6e4fc52d3f9480f0e000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fdfb047b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a362c226964223a22303366313039353238396537666464623838326631636233653135386163633163333064396465363036616632316339376261383531383231653862366561353335227d2c7b22696e646578223a352c226964223a22303338626663353062306533663065356466366434353130363930363563626661376162356433383261353833396363653832653063393633656462303236653934227d2c7b22696e646578223a322c226964223a22303335656236353462616436633634303938393462396234323238396134333631343837346337393834626465366230336161663666633164303438366439643435227d2c7b22696e646578223a332c226964223a22303238316431393863306464333733376139633339313931626332643161663764363561343432363161386136346436656637346436336632376366623565643932227d2c7b22696e646578223a312c226964223a22303235336363666434333962323965636130666539306361376336656161316639383537326130353461613264316435366537326164393663343636313037613835227d2c7b22696e646578223a342c226964223a22303233393637626261333036306266386164653036643962616434356430323835336636633632336534643466353264373637656235366466346433363461393966227d2c7b22696e646578223a372c226964223a22303231353836356261616237303630376634613234313361376139626139356162326333633032303264356237373331633638323465656634386538393966633930227d5d2c22706f735f7461626c65223a5b322c372c322c362c342c332c332c342c322c342c342c352c372c322c342c362c362c312c312c332c332c372c352c362c372c312c372c352c322c342c332c362c352c362c312c362c352c312c342c342c322c312c372c322c362c352c342c352c372c372c332c372c312c332c362c312c352c312c362c332c332c372c322c372c342c342c332c352c332c312c322c372c322c372c362c322c342c362c312c352c352c322c362c372c332c342c352c312c352c342c352c322c312c372c332c342c362c312c312c322c332c352c332c362c325d2c224d6178426c6f636b4368616e676556696577223a31303030307d7de0b053244d70b2c5622a4786d4922c2991e1fe880000")
		param := new(scom.SyncGenesisHeaderParam)
		param.ChainID = 3
		param.GenesisHeader = header7152785
		sink := common.NewZeroCopySink(nil)
		param.Serialization(sink)

		tx := &types.Transaction{
			SignedAddr: []common.Address{acct.Address},
		}

		native = NewNative(sink.Bytes(), tx, nil)
		err := ontHandler.SyncGenesisHeader(native)
		assert.NoError(t, err)
	}
	header1, _ := hex.DecodeString("00000000945d33e0aef7e6f8df67bbb42a22f306ce2d5f59e2f5974aabe23a7a6d7dfba10000000000000000000000000000000000000000000000000000000000000000be10c2305e2788739342a6d5be6ca4445d5680cd3b07e26c9af2e9b7fe794007eb1d285e0100000066d9abb328e3217ffd0c017b226c6561646572223a322c227672665f76616c7565223a224249433147524a685346307047655073304466522b44752b686351736b2f696f2b31682b4e2b4a656d394f7835476c584447555645774a3747665456636d38764b4351587448597273455338754b3669594a4d696f636f3d222c227672665f70726f6f66223a2258336f78614e624c4364477147612f435474507357717449776648336c2f37333646476a655870743331514a542b584b61362f30566b474a4c7a6c38753252686a6539702b74466e5059366a654b43746374574b30513d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000721035eb654bad6c6409894b9b42289a43614874c7984bde6b03aaf6fc1d0486d9d4521038bfc50b0e3f0e5df6d451069065cbfa7ab5d382a5839cce82e0c963edb026e9421023967bba3060bf8ade06d9bad45d02853f6c623e4d4f52d767eb56df4d364a99f210253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a852103f1095289e7fddb882f1cb3e158acc1c30d9de606af21c97ba851821e8b6ea535210281d198c0dd3737a9c39191bc2d1af7d65a44261a8a64d6ef74d63f27cfb5ed92210215865baab70607f4a2413a7a9ba95ab2c3c0202d5b7731c6824eef48e899fc900740cd2e4014e59799333541f330fb6f56fdb8b7ecf7845eb2fb173b1b2b33f4e7aae09fce3be84d17c89ee747c96f726f0b9e8fb2ab695094688e61f5ea64c2467c40754531cddec4affbbe5d210c77ef0e387bf9f236cd9016efd2c80f8eb3628ebdd0b79e24fad7a023dbd6049f9f4c5429afd23ef54797016b7812579dfa63c03d405028848c8be0d8ce2f9fea23dabd91f57a88d9eec1a8d76c2aade2353eb98ad94d089fb5a88014c27eff140bf79d8e0de8f14ad802e9732c6788028ce512ef364003fd8151c226103c89bc4750b35106b1ac8a4cfad4ec3f6565cdf1e442c806d882a040807f72b06f10f49a669baec4304a9f29f0b3814983689e8587dda5b733403b50f01bd376aaf5f0e2a9058e157264946160430061d7c4cb9d63669382fe72e17a182c1db8494d42888589eb9baff6d70ba3ff8316ea65cdded3468427ebfc403778eec68d959cb0f8ddfc8f939bd4814eec8167d5a7b7d726797c974a3bc03e97a4d6cd10dc7fe6fb17aacdd0066a5e15c9b69e9eac96b5ffaad6cfc13e57c9409054fff6d98474a5e3069042f9e1e127962779b284e32249e0a6c43ae81acb1d30726ce07112fba8a2641b13d31250014ecabcd64a616eb68126f968df7a77a4")
	header2, _ := hex.DecodeString("00000000edb82394c06326f54bc9027c916209f2f4d395097f3a1845fe4276f8404c53098367d628d7bdadc49984f436ae7c84e03035e87b1534256d9b6b6b7282cbec65805c5c040211fda21c600476c7d35b8f67d5cbf992864cd35af1d4d042bc5961071e285e02000000a5229c3a73e18053fd0c017b226c6561646572223a332c227672665f76616c7565223a22424770355144486c6258317a7a396265335131614a37584857573065723462656a5431336c4a57646f70784d7177745a673258386a336b466950566433504d584567482f6e6a61624852703634644351594f456b6242343d222c227672665f70726f6f66223a22717a55587837477a7369304552554765326d5a53576f34454354596d6f62504f34616948437a2b797a47563150714641474c6a444a7275597436674b505a692f5a455531614c374f30466474576d6c696c4f674238413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d000000000000000000000000000000000000000007210281d198c0dd3737a9c39191bc2d1af7d65a44261a8a64d6ef74d63f27cfb5ed9221038bfc50b0e3f0e5df6d451069065cbfa7ab5d382a5839cce82e0c963edb026e942103f1095289e7fddb882f1cb3e158acc1c30d9de606af21c97ba851821e8b6ea535210253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a8521035eb654bad6c6409894b9b42289a43614874c7984bde6b03aaf6fc1d0486d9d45210215865baab70607f4a2413a7a9ba95ab2c3c0202d5b7731c6824eef48e899fc9021023967bba3060bf8ade06d9bad45d02853f6c623e4d4f52d767eb56df4d364a99f0740782801c7e4dc73b4c223c5fcf2c49b6795aa545542e6ee9abf48f69eb820845766776ab9da46bdee52ce443dbeb7a3c27d0f3b774cd721ca07af2f1b5bb69508407a10be27550ba8a2c6789eef76585acefe3875530781d591a29cd7d0158583caf1aba069a986fa521bb78db54eba9cb3cd06a7cef44fb0fbceeffd194ed31a9240074e85056b1d2e66862b083d899fe952b8ecce67878abc579ca1b6be88104a284c4ee7de67371fcedf1c8917cf0c868af3284225bedad9716b55fe437abd8f264014fca561dd180f5c7b774128f6dfbd9b4ea5e4ab420e7057b40f11161a206cdb224b39da7121e0f954c74d473dd27f04ad3b3cf9d235677eef5e3425613e3a1c403bb7f8a9253eaa4629cc6f88d244f5dfacca63168e5070c2253423a636f7f96154e4227e57d0d5029bc90519765601d6d863d8e20971a6c655f4d1ca87922d9c40cd5337ee62459dc4f40b6f3278a73d0e1093aca4933036edc380304ed8f483d97928fbe6be8a0433cd0df3ac342639e9455bba2733a21d01bb6201271c4dcd05409fd6303ca41c406b4b4fc9912c0a65b6c727ec959e7a0cf7dff31da1cbd552041e00836687df511a2436441ca908ac456968bad465edcb39a2bb3fbb1dfb86c9")

	param := new(scom.SyncBlockHeaderParam)
	param.ChainID = 3
	param.Address = acct.Address
	param.Headers = append(param.Headers, header1)
	param.Headers = append(param.Headers, header2)
	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	tx := &types.Transaction{
		SignedAddr: []common.Address{acct.Address},
	}

	native = NewNative(sink.Bytes(), tx, native.GetCacheDB())
	err := ontHandler.SyncBlockHeader(native)
	assert.NoError(t, err)
}

func TestSyncBlockHeaderTwice(t *testing.T) {
	ontHandler := NewONTHandler()
	var native *native.NativeService
	{
		header0, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000038f121e642af587e34202511f911701ad81725a5fcc83eb6e4fc52d3f9480f0e000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fdfb047b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a362c226964223a22303366313039353238396537666464623838326631636233653135386163633163333064396465363036616632316339376261383531383231653862366561353335227d2c7b22696e646578223a352c226964223a22303338626663353062306533663065356466366434353130363930363563626661376162356433383261353833396363653832653063393633656462303236653934227d2c7b22696e646578223a322c226964223a22303335656236353462616436633634303938393462396234323238396134333631343837346337393834626465366230336161663666633164303438366439643435227d2c7b22696e646578223a332c226964223a22303238316431393863306464333733376139633339313931626332643161663764363561343432363161386136346436656637346436336632376366623565643932227d2c7b22696e646578223a312c226964223a22303235336363666434333962323965636130666539306361376336656161316639383537326130353461613264316435366537326164393663343636313037613835227d2c7b22696e646578223a342c226964223a22303233393637626261333036306266386164653036643962616434356430323835336636633632336534643466353264373637656235366466346433363461393966227d2c7b22696e646578223a372c226964223a22303231353836356261616237303630376634613234313361376139626139356162326333633032303264356237373331633638323465656634386538393966633930227d5d2c22706f735f7461626c65223a5b322c372c322c362c342c332c332c342c322c342c342c352c372c322c342c362c362c312c312c332c332c372c352c362c372c312c372c352c322c342c332c362c352c362c312c362c352c312c342c342c322c312c372c322c362c352c342c352c372c372c332c372c312c332c362c312c352c312c362c332c332c372c322c372c342c342c332c352c332c312c322c372c322c372c362c322c342c362c312c352c352c322c362c372c332c342c352c312c352c342c352c322c312c372c332c342c362c312c312c322c332c352c332c362c325d2c224d6178426c6f636b4368616e676556696577223a31303030307d7de0b053244d70b2c5622a4786d4922c2991e1fe880000")
		param := new(scom.SyncGenesisHeaderParam)
		param.ChainID = 3
		param.GenesisHeader = header0
		sink := common.NewZeroCopySink(nil)
		param.Serialization(sink)

		tx := &types.Transaction{
			SignedAddr: []common.Address{acct.Address},
		}

		native = NewNative(sink.Bytes(), tx, nil)
		err := ontHandler.SyncGenesisHeader(native)
		assert.NoError(t, err)
	}
	{
		header1, _ := hex.DecodeString("00000000945d33e0aef7e6f8df67bbb42a22f306ce2d5f59e2f5974aabe23a7a6d7dfba10000000000000000000000000000000000000000000000000000000000000000be10c2305e2788739342a6d5be6ca4445d5680cd3b07e26c9af2e9b7fe794007eb1d285e0100000066d9abb328e3217ffd0c017b226c6561646572223a322c227672665f76616c7565223a224249433147524a685346307047655073304466522b44752b686351736b2f696f2b31682b4e2b4a656d394f7835476c584447555645774a3747665456636d38764b4351587448597273455338754b3669594a4d696f636f3d222c227672665f70726f6f66223a2258336f78614e624c4364477147612f435474507357717449776648336c2f37333646476a655870743331514a542b584b61362f30566b474a4c7a6c38753252686a6539702b74466e5059366a654b43746374574b30513d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000721035eb654bad6c6409894b9b42289a43614874c7984bde6b03aaf6fc1d0486d9d4521038bfc50b0e3f0e5df6d451069065cbfa7ab5d382a5839cce82e0c963edb026e9421023967bba3060bf8ade06d9bad45d02853f6c623e4d4f52d767eb56df4d364a99f210253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a852103f1095289e7fddb882f1cb3e158acc1c30d9de606af21c97ba851821e8b6ea535210281d198c0dd3737a9c39191bc2d1af7d65a44261a8a64d6ef74d63f27cfb5ed92210215865baab70607f4a2413a7a9ba95ab2c3c0202d5b7731c6824eef48e899fc900740cd2e4014e59799333541f330fb6f56fdb8b7ecf7845eb2fb173b1b2b33f4e7aae09fce3be84d17c89ee747c96f726f0b9e8fb2ab695094688e61f5ea64c2467c40754531cddec4affbbe5d210c77ef0e387bf9f236cd9016efd2c80f8eb3628ebdd0b79e24fad7a023dbd6049f9f4c5429afd23ef54797016b7812579dfa63c03d405028848c8be0d8ce2f9fea23dabd91f57a88d9eec1a8d76c2aade2353eb98ad94d089fb5a88014c27eff140bf79d8e0de8f14ad802e9732c6788028ce512ef364003fd8151c226103c89bc4750b35106b1ac8a4cfad4ec3f6565cdf1e442c806d882a040807f72b06f10f49a669baec4304a9f29f0b3814983689e8587dda5b733403b50f01bd376aaf5f0e2a9058e157264946160430061d7c4cb9d63669382fe72e17a182c1db8494d42888589eb9baff6d70ba3ff8316ea65cdded3468427ebfc403778eec68d959cb0f8ddfc8f939bd4814eec8167d5a7b7d726797c974a3bc03e97a4d6cd10dc7fe6fb17aacdd0066a5e15c9b69e9eac96b5ffaad6cfc13e57c9409054fff6d98474a5e3069042f9e1e127962779b284e32249e0a6c43ae81acb1d30726ce07112fba8a2641b13d31250014ecabcd64a616eb68126f968df7a77a4")
		header2, _ := hex.DecodeString("00000000edb82394c06326f54bc9027c916209f2f4d395097f3a1845fe4276f8404c53098367d628d7bdadc49984f436ae7c84e03035e87b1534256d9b6b6b7282cbec65805c5c040211fda21c600476c7d35b8f67d5cbf992864cd35af1d4d042bc5961071e285e02000000a5229c3a73e18053fd0c017b226c6561646572223a332c227672665f76616c7565223a22424770355144486c6258317a7a396265335131614a37584857573065723462656a5431336c4a57646f70784d7177745a673258386a336b466950566433504d584567482f6e6a61624852703634644351594f456b6242343d222c227672665f70726f6f66223a22717a55587837477a7369304552554765326d5a53576f34454354596d6f62504f34616948437a2b797a47563150714641474c6a444a7275597436674b505a692f5a455531614c374f30466474576d6c696c4f674238413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d000000000000000000000000000000000000000007210281d198c0dd3737a9c39191bc2d1af7d65a44261a8a64d6ef74d63f27cfb5ed9221038bfc50b0e3f0e5df6d451069065cbfa7ab5d382a5839cce82e0c963edb026e942103f1095289e7fddb882f1cb3e158acc1c30d9de606af21c97ba851821e8b6ea535210253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a8521035eb654bad6c6409894b9b42289a43614874c7984bde6b03aaf6fc1d0486d9d45210215865baab70607f4a2413a7a9ba95ab2c3c0202d5b7731c6824eef48e899fc9021023967bba3060bf8ade06d9bad45d02853f6c623e4d4f52d767eb56df4d364a99f0740782801c7e4dc73b4c223c5fcf2c49b6795aa545542e6ee9abf48f69eb820845766776ab9da46bdee52ce443dbeb7a3c27d0f3b774cd721ca07af2f1b5bb69508407a10be27550ba8a2c6789eef76585acefe3875530781d591a29cd7d0158583caf1aba069a986fa521bb78db54eba9cb3cd06a7cef44fb0fbceeffd194ed31a9240074e85056b1d2e66862b083d899fe952b8ecce67878abc579ca1b6be88104a284c4ee7de67371fcedf1c8917cf0c868af3284225bedad9716b55fe437abd8f264014fca561dd180f5c7b774128f6dfbd9b4ea5e4ab420e7057b40f11161a206cdb224b39da7121e0f954c74d473dd27f04ad3b3cf9d235677eef5e3425613e3a1c403bb7f8a9253eaa4629cc6f88d244f5dfacca63168e5070c2253423a636f7f96154e4227e57d0d5029bc90519765601d6d863d8e20971a6c655f4d1ca87922d9c40cd5337ee62459dc4f40b6f3278a73d0e1093aca4933036edc380304ed8f483d97928fbe6be8a0433cd0df3ac342639e9455bba2733a21d01bb6201271c4dcd05409fd6303ca41c406b4b4fc9912c0a65b6c727ec959e7a0cf7dff31da1cbd552041e00836687df511a2436441ca908ac456968bad465edcb39a2bb3fbb1dfb86c9")
		param := new(scom.SyncBlockHeaderParam)
		param.ChainID = 3
		param.Address = acct.Address
		param.Headers = append(param.Headers, header1)
		param.Headers = append(param.Headers, header2)
		sink := common.NewZeroCopySink(nil)
		param.Serialization(sink)

		tx := &types.Transaction{
			SignedAddr: []common.Address{acct.Address},
		}

		native = NewNative(sink.Bytes(), tx, native.GetCacheDB())
		err := ontHandler.SyncBlockHeader(native)
		assert.NoError(t, err)
	}
	{
		header1, _ := hex.DecodeString("00000000945d33e0aef7e6f8df67bbb42a22f306ce2d5f59e2f5974aabe23a7a6d7dfba10000000000000000000000000000000000000000000000000000000000000000be10c2305e2788739342a6d5be6ca4445d5680cd3b07e26c9af2e9b7fe794007eb1d285e0100000066d9abb328e3217ffd0c017b226c6561646572223a322c227672665f76616c7565223a224249433147524a685346307047655073304466522b44752b686351736b2f696f2b31682b4e2b4a656d394f7835476c584447555645774a3747665456636d38764b4351587448597273455338754b3669594a4d696f636f3d222c227672665f70726f6f66223a2258336f78614e624c4364477147612f435474507357717449776648336c2f37333646476a655870743331514a542b584b61362f30566b474a4c7a6c38753252686a6539702b74466e5059366a654b43746374574b30513d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000721035eb654bad6c6409894b9b42289a43614874c7984bde6b03aaf6fc1d0486d9d4521038bfc50b0e3f0e5df6d451069065cbfa7ab5d382a5839cce82e0c963edb026e9421023967bba3060bf8ade06d9bad45d02853f6c623e4d4f52d767eb56df4d364a99f210253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a852103f1095289e7fddb882f1cb3e158acc1c30d9de606af21c97ba851821e8b6ea535210281d198c0dd3737a9c39191bc2d1af7d65a44261a8a64d6ef74d63f27cfb5ed92210215865baab70607f4a2413a7a9ba95ab2c3c0202d5b7731c6824eef48e899fc900740cd2e4014e59799333541f330fb6f56fdb8b7ecf7845eb2fb173b1b2b33f4e7aae09fce3be84d17c89ee747c96f726f0b9e8fb2ab695094688e61f5ea64c2467c40754531cddec4affbbe5d210c77ef0e387bf9f236cd9016efd2c80f8eb3628ebdd0b79e24fad7a023dbd6049f9f4c5429afd23ef54797016b7812579dfa63c03d405028848c8be0d8ce2f9fea23dabd91f57a88d9eec1a8d76c2aade2353eb98ad94d089fb5a88014c27eff140bf79d8e0de8f14ad802e9732c6788028ce512ef364003fd8151c226103c89bc4750b35106b1ac8a4cfad4ec3f6565cdf1e442c806d882a040807f72b06f10f49a669baec4304a9f29f0b3814983689e8587dda5b733403b50f01bd376aaf5f0e2a9058e157264946160430061d7c4cb9d63669382fe72e17a182c1db8494d42888589eb9baff6d70ba3ff8316ea65cdded3468427ebfc403778eec68d959cb0f8ddfc8f939bd4814eec8167d5a7b7d726797c974a3bc03e97a4d6cd10dc7fe6fb17aacdd0066a5e15c9b69e9eac96b5ffaad6cfc13e57c9409054fff6d98474a5e3069042f9e1e127962779b284e32249e0a6c43ae81acb1d30726ce07112fba8a2641b13d31250014ecabcd64a616eb68126f968df7a77a4")
		header2, _ := hex.DecodeString("00000000edb82394c06326f54bc9027c916209f2f4d395097f3a1845fe4276f8404c53098367d628d7bdadc49984f436ae7c84e03035e87b1534256d9b6b6b7282cbec65805c5c040211fda21c600476c7d35b8f67d5cbf992864cd35af1d4d042bc5961071e285e02000000a5229c3a73e18053fd0c017b226c6561646572223a332c227672665f76616c7565223a22424770355144486c6258317a7a396265335131614a37584857573065723462656a5431336c4a57646f70784d7177745a673258386a336b466950566433504d584567482f6e6a61624852703634644351594f456b6242343d222c227672665f70726f6f66223a22717a55587837477a7369304552554765326d5a53576f34454354596d6f62504f34616948437a2b797a47563150714641474c6a444a7275597436674b505a692f5a455531614c374f30466474576d6c696c4f674238413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d000000000000000000000000000000000000000007210281d198c0dd3737a9c39191bc2d1af7d65a44261a8a64d6ef74d63f27cfb5ed9221038bfc50b0e3f0e5df6d451069065cbfa7ab5d382a5839cce82e0c963edb026e942103f1095289e7fddb882f1cb3e158acc1c30d9de606af21c97ba851821e8b6ea535210253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a8521035eb654bad6c6409894b9b42289a43614874c7984bde6b03aaf6fc1d0486d9d45210215865baab70607f4a2413a7a9ba95ab2c3c0202d5b7731c6824eef48e899fc9021023967bba3060bf8ade06d9bad45d02853f6c623e4d4f52d767eb56df4d364a99f0740782801c7e4dc73b4c223c5fcf2c49b6795aa545542e6ee9abf48f69eb820845766776ab9da46bdee52ce443dbeb7a3c27d0f3b774cd721ca07af2f1b5bb69508407a10be27550ba8a2c6789eef76585acefe3875530781d591a29cd7d0158583caf1aba069a986fa521bb78db54eba9cb3cd06a7cef44fb0fbceeffd194ed31a9240074e85056b1d2e66862b083d899fe952b8ecce67878abc579ca1b6be88104a284c4ee7de67371fcedf1c8917cf0c868af3284225bedad9716b55fe437abd8f264014fca561dd180f5c7b774128f6dfbd9b4ea5e4ab420e7057b40f11161a206cdb224b39da7121e0f954c74d473dd27f04ad3b3cf9d235677eef5e3425613e3a1c403bb7f8a9253eaa4629cc6f88d244f5dfacca63168e5070c2253423a636f7f96154e4227e57d0d5029bc90519765601d6d863d8e20971a6c655f4d1ca87922d9c40cd5337ee62459dc4f40b6f3278a73d0e1093aca4933036edc380304ed8f483d97928fbe6be8a0433cd0df3ac342639e9455bba2733a21d01bb6201271c4dcd05409fd6303ca41c406b4b4fc9912c0a65b6c727ec959e7a0cf7dff31da1cbd552041e00836687df511a2436441ca908ac456968bad465edcb39a2bb3fbb1dfb86c9")

		param := new(scom.SyncBlockHeaderParam)
		param.ChainID = 3
		param.Address = acct.Address
		param.Headers = append(param.Headers, header1)
		param.Headers = append(param.Headers, header2)
		sink := common.NewZeroCopySink(nil)
		param.Serialization(sink)

		tx := &types.Transaction{
			SignedAddr: []common.Address{acct.Address},
		}

		native = NewNative(sink.Bytes(), tx, native.GetCacheDB())
		err := ontHandler.SyncBlockHeader(native)
		assert.Error(t, err)
	}
}
