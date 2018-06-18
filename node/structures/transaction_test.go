package structures

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"

	"time"

	"testing"

	"github.com/gelembjuk/democoin/lib/wallet"
)

func TestHash(t *testing.T) {
	PubKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}

	inputs := []TXInput{
		TXInput{[]byte{1, 2, 3}, 0, []byte{}, PubKey},
		TXInput{[]byte{4, 5, 6}, 1, []byte{}, PubKey},
	}

	outputs := []TXOutput{
		TXOutput{1, []byte{4, 3, 2, 1}},
		TXOutput{2, PubKey},
	}

	newTX := Transaction{nil, inputs, outputs, 0}

	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	time, _ := time.Parse(layout, str)
	newTX.Time = time.UnixNano()

	newTX.Hash()

	expected := "913acaf3f296c048df72565c3940b33afdbdba1b0683bd71969491a549a86ba3"

	expectedBytes, _ := hex.DecodeString(expected)

	if bytes.Compare(expectedBytes, newTX.ID) != 0 {
		t.Fatalf("Got \n%x\nexpected\n%x", newTX.ID, expected)
	}
}

func TestSignature(t *testing.T) {
	// wallet wallet address, wallets file, transaction, input transactions
	testSets := [][]string{
		[]string{
			"1yPg8JYsMepEBSyTFj1ZhBXcZXbrzrvtM",
			"wallet_1yPg8JYsMepEBSyTFj1ZhBXcZXbrzrvtM.dat",
			"3cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a000000ffacff820201012072d458be3467230355ce5d996ff980dc5922327c64bb80b3b16e348476bc8c460340f7a860a846d40e7cd834c612dedaeba6bd9bdba7fab343325e3b73f6a3811f3394e2800c00d91f08cba5cf22483c3cad93c37e31df6c0ad6a2120a99b6ab686000010201fef03f0114bdb034351ac4c903f4108ed6fc25d2079cbdf5440001fe224001140aaa38be11a67bed469208615f137beb4d5ed5fd0001f82a3d2e1a47018a8e00",
			"0fff8b040102ff8c00010401ff8200003cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a0000006eff8c000100012072d458be3467230355ce5d996ff980dc5922327c64bb80b3b16e348476bc8c46010102010222546869732069732074686520696e697469616c20626c6f636b20696e20636861696e00010101fe244001140aaa38be11a67bed469208615f137beb4d5ed5fd0000",
		},
		[]string{
			"1NUmoDL88h7JTukJJtf3AU6oJ78VfjWMi1",
			"wallet_1NUmoDL88h7JTukJJtf3AU6oJ78VfjWMi1.dat",
			"3cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a000000ffbaff8202010120cc5a25ccd4542f8f226eb927a662f37b403644369a621f38474350036fd1643701020240d76efc1748e2b2aa3203616752e5e3123fdcfba152de4e500bb37baff8bcd6ddfe45954df2cbe05958e0dbf033c918e4e9968caaceba53465e4efc0a0e78665a00010201f8682374236e46843f0114c1bc99e4005cd69c6f7fa571c189fc0e2f9005fc0001f87c0616b823b9ee3f0114eb9b4452b3aa30dac0c71b5a002b71220d76c6410001f82a3d314bf374fba000",
			"0fff8b040102ff8c00010401ff8200003cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a000000fe0121ff8c0001000120cc5a25ccd4542f8f226eb927a662f37b403644369a621f38474350036fd1643701010120187a47232a3af6583df42909e737c033002a8748e2898fbab354a53a673289ad01020140c198bcac044c69d83a35b8470e5c301de6b8dc521a8c9208c44b360cdc8a6093773d8114ddf5b25c2826b05a79371e1eeb5957b9423e46e97c6ef13f2464a4850140d76efc1748e2b2aa3203616752e5e3123fdcfba152de4e500bb37baff8bcd6ddfe45954df2cbe05958e0dbf033c918e4e9968caaceba53465e4efc0a0e78665a00010201f87b14ae47e17a843f0114c1bc99e4005cd69c6f7fa571c189fc0e2f9005fc0001f80ad7a3703d0aef3f0114eb9b4452b3aa30dac0c71b5a002b71220d76c6410001f82a3d3141470338fe00",
		},
	}

	for _, test := range testSets {
		ws := wallet.Wallets{}
		ws.WalletsFile = "testsdata/" + test[1]
		ws.LoadFromFile()

		w, e := ws.GetWallet(test[0])

		if e != nil {
			t.Fatalf("Error: %s", e.Error())
		}

		tb, _ := hex.DecodeString(test[2])
		tx := Transaction{}
		tx.DeserializeTransaction(tb)

		prevTXs := map[int]*Transaction{}
		tb, _ = hex.DecodeString(test[3])
		decoder := gob.NewDecoder(bytes.NewReader(tb))
		decoder.Decode(&prevTXs)

		signData, err := tx.PrepareSignData(prevTXs)

		if err != nil {
			t.Fatalf("Getting sign data Error: %s", err.Error())
		}

		err = tx.SignData(w.GetPrivateKey(), w.GetPublicKey(), signData)

		if err != nil {
			t.Fatalf("Signing Error: %s", err.Error())
		}

		err = tx.Verify(prevTXs)

		if err != nil {
			t.Fatalf("Verify Error: %s", err.Error())
		}

	}
}

/*
func TestSignatureAndVerify(t *testing.T) {
	// wallet wallet address, wallets file, transaction, input transactions, tx before verify
	testSets := [][]string{
		[]string{
			"12Dk8TwRi196Lpy77qoz82UnqexhnwWtkm",
			"wallet_12Dk8TwRi196Lpy77qoz82UnqexhnwWtkm.dat",
			"3cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a000000ffb8ff8202010120722bf6ef370589e35f1ccff33608e9a075d564877a15c80b42520f98e611d0490340cf0c4157ebb2cf8565de6690e9f0ade8e26ec45b9e39ed2c4b0e6b7a55e24a3f9ff7465bacc5ecec0c828cbe1d4dd07c6c16bd409f77834465eea52a24f6e4bd00010201f826ffd5db4d8d833f0114022a5faad64ef73a4987af7ce5b84a146efaf9b80001f8a0aa027b6db23d3f01140d611e8da5c528f2f4a75119931677dd3cca5fba0001f82a3db819db5b1eec00",
			"0fff8b040102ff8c00010401ff8200003cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a000000fe011fff8c0001000120722bf6ef370589e35f1ccff33608e9a075d564877a15c80b42520f98e611d049010101206d43b7845b2569b036d3c2ceeb97c415fe9841ed116f6ec4872120e0ac89021e0240727e28a66858da898593fd25cc9450427055f8da67512e7cbe6ab2eb5596ba0f1d98ed814bca0d999582d28d9e7bf753fab37728512a7d7f4c07c03f5084678a0140cbb9f2a8b29ec00c4ea20f1d9a1d8c49238b1a01d2ec3dd4001a7f27afe8b6065a647ac967abbded6cb0762cf34b3b17f552d42cb22041f2e6ebcac7cb04ccf500010201f87b14ae47e17a843f01140d611e8da5c528f2f4a75119931677dd3cca5fba0001f8ae47e17a14aeef3f01140b8cd4d5171a2db9d6877171ee708bfb3fd2a81d0001f82a3db80dd11e949200",
			"3cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a000000fe011bff820120506913fcc9bce34f3ee9ed74c75c52b413936d5c17bf8181ba4c5277fc9fde1701010120722bf6ef370589e35f1ccff33608e9a075d564877a15c80b42520f98e611d049023fa6a4ad6466bd42e4a62e089ec863c5f9f7bdb7fc7de08f0694fefcdbf400944d1fbf84ababcd58c6f025a6062525a29b09ea6a8749d910abe25de28fe52d7c0140cf0c4157ebb2cf8565de6690e9f0ade8e26ec45b9e39ed2c4b0e6b7a55e24a3f9ff7465bacc5ecec0c828cbe1d4dd07c6c16bd409f77834465eea52a24f6e4bd00010201f826ffd5db4d8d833f0114022a5faad64ef73a4987af7ce5b84a146efaf9b80001f8a0aa027b6db23d3f01140d611e8da5c528f2f4a75119931677dd3cca5fba0001f82a3db819db5b1eec00",
		},
	}

	for _, test := range testSets {
		ws := wallet.Wallets{}
		ws.WalletsFile = "testsdata/" + test[1]
		ws.LoadFromFile()

		w, e := ws.GetWallet(test[0])

		if e != nil {
			t.Fatalf("Error: %s", e.Error())
		}

		tb, _ := hex.DecodeString(test[2])
		tx := Transaction{}
		tx.DeserializeTransaction(tb)

		prevTXs := map[int]*Transaction{}
		tb, _ = hex.DecodeString(test[3])
		decoder := gob.NewDecoder(bytes.NewReader(tb))
		decoder.Decode(&prevTXs)

		signData, err := tx.PrepareSignData(prevTXs)

		if err != nil {
			t.Fatalf("Getting sign data Error: %s", err.Error())
		}

		err = tx.SignData(w.GetPrivateKey(), signData)

		if err != nil {
			t.Fatalf("Signing Error: %s", err.Error())
		}

		//txdata, _ := tx.Serialize()
		//txstr := hex.EncodeToString(txdata)

		tb, _ = hex.DecodeString(test[4])
		tx2 := Transaction{}
		tx2.DeserializeTransaction(tb)

		err = tx.Verify(prevTXs)

		if err != nil {
			t.Fatalf("Verify 1 Error: %s", err.Error())
		}

		err = tx2.Verify(prevTXs)

		if err != nil {
			t.Fatalf("Verify 2 Error: %s", err.Error())
		}
	}
}
*/
/*
func TestVerify(t *testing.T) {
	// wallet wallet address, wallets file, transaction, input transactions
	testSets := [][]string{
		[]string{
			"3cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a000000fffaff820120902e5dde6bd4215c9e795ca0f273f7ea1bb94f2ed963544514ef6b4e891dfde501010120754eec4c5b83f6ff4a1600254d9d1c46cfbc33db9b1ba800484cefd73e1dae0e023fe004de1df4afb86680de18ae872dc14a2a00bd7a39265b73eed526ceac8ffd330994ecdb217e85d03e0e3484a26bb02b65455c7faae1c8f0110568590cd390014069ed4eb77cf4615d294394b22e072d8df4580ce3005bdb3a4c4ebba77c3e5c20fe5201de547332f05e88c89e9bc6ab9e6f09a47c4b8b65c1aa41d3d61018b0c600010101f87b14ae47e17a843f0114d43ac7037342df5cb1a3f7c855bb9c58c3a5aeab0001f82a3db639c37f657600",
			"0fff8b040102ff8c00010401ff8200003cff810301010b5472616e73616374696f6e01ff8200010401024944010a00010356696e01ff86000104566f757401ff8a00010454696d65010400000024ff85020101155b5d7472616e73616374696f6e2e5458496e70757401ff860001ff84000040ff83030101075458496e70757401ff84000104010454786964010a000104566f757401040001095369676e6174757265010a0001065075624b6579010a00000025ff89020101165b5d7472616e73616374696f6e2e54584f757470757401ff8a0001ff8800002fff870301010854584f757470757401ff88000102010556616c7565010800010a5075624b657948617368010a000000fffeff8c0001000120754eec4c5b83f6ff4a1600254d9d1c46cfbc33db9b1ba800484cefd73e1dae0e010101205a8ce6222ce533b3056f20eae7743ff1ffb484cb6754f1edd5bdd900c1af0c4b0240e70500b4e113214e19a51c8f2aa22392ac90ded4dc2bb62ef860e74d8bd97277a000dd20928882c570a41dff092cceed602584a8c540c3af4a074215d1c9b4790140a62b2d7c1b9c3e1b854f19a4a8fa1e334037ece72d786cd7c293a1337c96807b19b775e1ce25e8d23b866070e7ddbac3848e736d2251e05c101f018f9a312b6800010101f87b14ae47e17a843f01142c053a6bde73afea1820315acfd7efcd8a683cc90001f82a3db6399904faf400",
		},
	}

	for _, test := range testSets {

		tb, _ := hex.DecodeString(test[0])
		tx := Transaction{}
		tx.DeserializeTransaction(tb)

		prevTXs := map[int]*Transaction{}
		tb, _ = hex.DecodeString(test[1])
		decoder := gob.NewDecoder(bytes.NewReader(tb))
		decoder.Decode(&prevTXs)

		err := tx.Verify(prevTXs)

		if err != nil {
			t.Fatalf("Verify Error: %s", err.Error())
		}

	}
}
*/