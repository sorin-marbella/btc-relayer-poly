package utils

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	USER = "test"
	PWD  = "test"

	BLOCK = "00000020b797de83cf1f8865be7e894ccc8116840bc5280d4d2ddd7c330c00f54e786e04330486c1ed6fcb9c163d46f0a986a701cc465fe3ac3a9faae13b84493c9c0797cfe03b5effff7f20000000001c020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff05029f040101ffffffff029aac5402000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000000000000266a24aa21a9edbff1389f1cd6097d4d9b587aaa0a52d6a6fbe1d0b3aa2c46597b2894c8c58e4a012000000000000000000000000000000000000000000000000000000000000000000000000001000000013df3f407d6eef47c9dc691aaa19a6478d44fc76878bfb18b8f1a772960acdf70020000006a473044022003d280c603b6e179854dab4ff65c07381cca13f17ffa67e438c0b1bf85ed64220220366f5a267ce4f2f0b970849fa4ecb396651aa2e26ba59985693fad2694b7336d012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03978101000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a48044fe29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000010000000128b8410b6d07b0f368e78efc4ec5d355d211a50bb6a06d23e41398a8848f6593020000006a473044022051a837efac5b8804fcf5d0fc6e961f79bd58f4e674a4a86f34e88225096757d302202f439f9703928c5a5838290cbde0bd7f4b2af323a2e59114a8f77ae1f858989f012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03f17e00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a4a795fe29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000014e287c4c02a3f996efa3d09f753e3d6290b483b729592becd4dde6c053e95b19020000006a47304402206108fa5abccce307f6b6c3174c61963110168c5c43cfc4a8b91d360d631573a702200b57135f0ae2fc0c07ce5703681e74431a1cbcbcfe15bea58f75b6111873c09d012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03f67f00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a44737ff29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000010000000110055bde86659522ed779da6f2031c4b28977282afda930cd3c5ba100562be9e020000006a47304402207557bcc71429ebbde0a88e197c6cc01381365c2ee5dcf66b6e3fbf64ff371fa502203d06f14b2484a4413182de08feb82c32fc5192bbc7e5efeadd12e49052b43ee4012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff038edb00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a40868fd29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000010000000191d7a91baa9accc174c19386d2b08a1f6365f1f9ce582e1de97033eac3ad6f26020000006a473044022057b2105a6ec7bcf3565902058c840b40dbcf80f18fcb22828ab087c348d38c790220261d939bbf70b6c7d2f69c6808eefa7c1e685ae967b18f4672b4db0ce4c8725f012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03f65c00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a4b4a3fd29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac000000000100000001bd2f9a1e6d19ed8cddbac31c459d8c56c3f8ac16121b32f63dd114dd24e66ee6020000006b483045022100f8c09bc5f1ac9094a2a579686d39e66763bd67daf8eb88e6017dc4cf1eb06d5902205caf297d578847e432595aa00a7e4ed809d116e6505556ffc11b4077d00d7aa8012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03247600000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a437defe29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac000000000100000001c819eaf6ed4197f55627d17566d083a721f96db82142d72605de30e59fb7024b020000006b483045022100fc6f5c1619f495d852a6c38ed41a7c11b09cb1956d83e2309945403f5675140002202e70a4046f466e3822a61b825c0e5c342ca764d9177ce784690526e92ef2aee5012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03b45b01000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a47690fd29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac000000000100000001af818776cbd3e63b9125e23d34613c11d7c421cb93426c454d878cf32bd0deeb020000006b48304502210084eb5695349ddea7050e61f1d7d0f91aec5ef077577354ace2049128c5d1255f0220798d211311e7d7ec4f23a31ca03c270f0a96fdddbc2572e14e489594dfca8669012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff033d5b01000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a49709f029010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac000000000100000001c12980b44ef668cc2509e9e7e6e0de5523ccfbbe04d2234073514e979856a43c020000006b483045022100a95dba335518b4e5aa69c69b898874598bcc8d4321cbdd4b6e9ebddb37ef947d0220025e67ef3b91173029ef5538588032b23b6f079f8cfd7106bb1cfef93a0fe8a8012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff036a0d00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a42503f029010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac000000000100000001bdf05ddf8dceb59f7b520e9580d0a9ed6a59c158a14fa7a1226f564c3166e3b9020000006a47304402207170e9a2ddfca39506653c11731d4896c477a7acf685e90fb03d65d71bdfd960022051fb05f68782fb929197c59dd5dd366f31701aca94c397f5035da8edc37bb07a012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff0317de00000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a46ca1fc29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000017383fe6fea7b97c15c928d9f48e4ed1e0cbfef3eae831387ce1402722650f58b020000006a4730440220107811c2c28a5f8f39e6feeaf499c43c769a7ebcad65b4b3ac7b7f884072a57b02205be2c47753147b1d3c9f24e4ebaf49e304a4c7620302a94f1404dbbacd2c24c1012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff0358a500000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a4a41af129010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000012e903485e67691edd5cf257f8a74787bfcf23765f6ee01472c73cda625d46f0a020000006a4730440220154f4619785fde888d1d36a6388fb2ad12e51865299550fe61566a7cdf69c2ab022063282f25a3444fbdc03711b2a46d73ca3a46e94ec2b415f678463f51d0e06d5f012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03644100000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a44258ff29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000019c4c0bca3b44c3126bfb01ad36359d17b5a09417e586baf6f50510cdedbdab0a020000006a47304402203950fadcca927474f149d371f2a5e0901d637357d112757cd442e6c635d5b89402200cf98bd2016d06de05b83380d2191e99db03f7091529f994d5993f7619130f33012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03fd2201000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a4a241fe29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000010000000122c54e204104b8bb7582e28220267dde7c5369259772c4995f79a7e7f4b659a3020000006a473044022062e0d1c063d906f9f0acc6a20e8018aa50dfa196aa10c1b075ee286e672aa504022023fa0085f2be766bb980ac32e8b33febed7efebff7b9f9d706bec1f11cd1c7c0012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03f0a000000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a43dd7f329010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000015d526edd27d7f766885c721c6c4c5a1ba2f3b9ae73c3b1a678dea2515b9c1835020000006a47304402204fa0d40d8a72a03a24e2d49825e74ce317964e7e4589ede58965685c15fef15d02202e001f772b6140043bc04e9c6cf8fea5b25c580de944639aa4c6b32f7ed41f36012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff0335ee00000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a4ee5cfe29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000012af0f6bd93e5346ea956fccc29c1c4e2f8c8944e6a8ebdf2f3a200bcdf7ba68f020000006a4730440220346650f43ded0f053981b643fd15503375e51f5f38e72cf2b32d6e474d8eb7c20220483473960ee4bdc9a79e61783dfcfad2373755a955e0843026ed271369f952f1012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff030c3e01000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a40cb6f029010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000010000000110c4a068a1572f562489b0ca2c960e52ca380b3269c3f43cee42424cc86be857020000006a4730440220397c8fa00a2b341d66383c041adce2664151c46651b42f2adcfb1aa8cbec18f7022059cad392ba0bd6e7202c11d87dab8c7aba50741b48cfb12105a06f9a8e69489e012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03d02c01000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a4a8adfd29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000010041434660bade294101bcd2946bbf877f7d0c27fa3b7e163b28dd0f7c13b842020000006a47304402200fb57380e0a5d0700b2254ef4932f4684ae87862c6f5133d4e154d2199b9cbce02201a5a89d10221c2ed371e08dcb2325e19037b6786592303e73d84834cd577302e012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff039d8e00000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a48a4efe29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000018f67eb6dabd8e4fdaa0dc5a1ac18b944d567e3513e964d864683be93d3a8a4d6020000006a47304402201f36864bf5b241778e60b8a49d565663050f3bbd29baeaefff09f3a6b4516a2802205c35c871b9678ba1a021160c3472d3decaea5cd6fffbd71650604be4975a91f4012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff039a6c00000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a434c8fe29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000015db9e5e6e5ad420b4ce05d36999f1291242ae40f63c3bae13cfc9a0423d11a64020000006b483045022100ccdaad2fd59066c47731246562ce58f0fdd3372f184e959d8a93e1a4611e41c902206dd4759f0093192ad2d4274ccbc46f164f1c7b3153e4a1efdcefb14c146aba86012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03040900000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a4b4e6ff29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac000000000100000001ccf481c3cdd1d027da1ca9f86000a4597b70d41047ec2f6f3663b5263e05588b020000006b483045022100abcc712c4b9e47c2fca93252c7808733db700b88acd3505d53365598ed738f3b022058f8f79449baf40e82406dd6cec90aaa4fdcbbb050437d14912c69b20f259bf0012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03057e00000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a4faf7fd29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac000000000100000001c0f43054706c5ed6d7aee7789f306ef6f4aa99abc6d34ee1822a6f9d353afc1f020000006b483045022100baf8df9d5a3c3868522432a2ff6059e7e0b4aa64df6348fcbc97559747afbde80220033ad1d800523eeeb0a9b03d95352ad53191fffb53c2ad6077a810c893597c1c012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff034b6700000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a47aadff29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000011b8abcb18afc5a240181e294cf1fc37a91872f54ff9ab37cb4ecbebd83b340f1020000006b483045022100ea68e5f150ef95a14d248d4e7ac7492cf3bba3a40156467a8469cd5bb1e9a0b502201637639e5cd9c086e6f8e96da4f7a7011238223a51226e4f6b4bbe2f0c5d8d7f012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03a77900000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000000000003d6a3b660300000000000000000000000000000014e553510e9f7eb980b3f225056d144529525cc5ec14f3b8a17f1f957f60c88f105e32ebff3f022e56a44c88fe29010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000001000000000104bc2c43486cafbaae40ff7f6ab13d68173967c7cc3bd732f465de2ecd98c2466f0000000000ffffffff1983915b8cce815a1dfd6310e9de38b49de5e7b1c6104905cd36075f0ec9dcac00000000fd5e0200483045022100c855aa6843246008452c68bbabc20e36995a257b36ffc8670e8c24a0453943bb022034f468b9d6e96543ed5b3edf559f0befcd66b9e981f8d7a47eead2d0782cfdfe0147304402206ea52077c3f131d06d183d979948ad025fe63f7b5c50c34e46c66f54bb08bfff02203e71ca4f03877dee1adf86fec3acdc0844121e21194ec85ab55aace3817a61370148304502210092f3872327c794b7527251325fd3aae48542babc2511ac9244c78eb0e8da100502201a07d17086ad2d4de89c2fb8d8addc84ec2424aee349343adf0531da66e6b5b20147304402201350c09eca30319525f9746bb1ba26d353d57a50e913619cc114705c1c95ecd10220412c9f21e7e33e7330e04ef178e4b02e1a4a903d2f5a5a8a96c7cd9ae33a2b21014730440220690e5509a489e5880a0ab4157d0371a121a3697504e466896451a5a6d82371bb0220052a3c9530d35f8103fa594b64e86d830936dfce25d05a4f2e86f4f2aa6922f1014cf15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57aeffffffff42a91ab1a9469a7571d45f094803e5f7aee95f1c767b95e514f92a02db898fe30000000000ffffffff107351cdfbea0f709e74b2db056fe6eac4c18a59d5266deb7ee50a2d10f6299e00000000fd5f0200473044022025c75d237d5b86530663d208a1116c9dc55099b9f54726e8bff7862f0ba496fc02207d260ed7b470c89672b36cc527313619fb99350dfab2512b22162483c300a35b01483045022100a20501477b754a0456d5d3ee785b53815f8be5943b2436bdab6aa9997d800aab02204f9089b3125817e665a31ae450d08efcf8f6a7e4a7713353c160ae4c31dc755301483045022100e70d1f8f02274df604767d1c8666c2c39a519c92e0ea549ae36d0cd043e9df4c02201bafd3eaa2ca34fa2dd9a698c36e1332fa3600b2954e165f2117ed4cfe8bfae60147304402202fe6730423d0fd05ab2ba23a26145d0455c215f8742f024bec82bede8965c386022060a1204e6370bbf7132290412add94542b649e651a0c7ecae8e48b9fe690333801483045022100bb54f645afc7448bbb4088c2ffa8e9743bc514c866c07705230dd8ca151328890220695dfedd0ed778bd959f0c7d09a6b961ee6ce3a7a5ce826587d0bead209d6d00014cf15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57aeffffffff0268660000000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188acc91000000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c27430700483045022100b2bdb931dc2fd98ba2f85e88761a6d421a270f9186ed321ea46eedf408fff5750220758a764dc13f30663a525a1d58ccf661ba61d290acd5e9e01000c27d4c0fe164014830450221008558ae2f9ad102fedf31802f7e2d4bf5f84ed2c70fcf791767d9066830d6175d0220750c888c803a2dbd8d51ab6ca15e8c009a5f8967bdeecfef0902fe1ce3d369c60147304402200400ea4daf6b26a084403fc7920517fdd379ba9b1e1a2c0eb0be547ea1b4c12b0220408a11486693a3dd9d532a0d664aa7ca4c980f9e9dfc797c26831af31836ef97014830450221009b60e59a3f25b541276d932dbd8be432985f667eb498162edd3fad6d706145960220523072467cf82cd1084c13ce3c24fc08f7ec803417bdf8543ebd97de99e9d06001473044022069c974249621045929d5016dd16b8517673d70d39f4dceb85cf1762fe6cfdffc0220630bc41306104c7dbd137ae30759d63e973e011fdd2f9797557f1cecb4a5991b01f15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae000700473044022079096901c22dc16c3577fc279154e9d744b3aa9883b3cc6650c50481452f3e5c022062700d98b1a5ca25fe218f10647d4ab2f9f707e134bc0efeb39dcab8560449fd0147304402204763ca6ee7b5f47db5c0d670e5053b05d2121ad76e03c37eef5e0d2020f2bd4c0220272e35ad92ad22ba78ce65e55efbef8894143f36f1984140d20031758ae2b80f01473044022026a1871268c54bc1087a0a49cfb0af746e3a8365b797a40eff3599a1aa10ac7e02204d57e9494c40e8c414be49df52179282c437c8c8dced6607249710718a00f9ea01483045022100c35738577782a1ea7d9d5f46540e09605fb6a403a0e46da9be9c12f60b50b8f0022001a3507fbde77bec354ea7e12fd795afebc01e1987368717f9026dbeee12606d0147304402205c5b290046833e0ff3970274ea94701f224bc75d5a7aa978c1bfeb69f25c55df02204274cd4908d3a086a089323a81804b850a9e9ed8b6698d53d45080578c3e207d01f15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae00000000000100000002f0c48c6ee256cda61c7f0a074feb09edb18c06b0eb3062f6a67e4dfffba509be00000000fd600200483045022100965e00db44860a7f42aae09362a212fed8fc80551f81270983975478dcfc6387022049fe9e05fd926c12e12cff88cf6d8b95edd35fc8b50656a8512060c08d869e9801483045022100b81591330b56ad490006f19719e5ef385edf869fda195f1efd01ecab1c72519102200797092883d036dfc9c09d9eb18496c7bdfad4953f5cc34d888fc61bcd4f9f1f014830450221009a7b21d5c97f1551e3255832988a1068a58eafe50bedc06fa2889e1823582ad902205b38a7d050f6b258fa7375a4199dbb0d431d1d27a4f8da76954469fc23f42a4f014730440220108035ad70bb5bbb1e8c13b3dcfd2b61d49ea97edf570292c4c15222dce7db740220551d802fb67956971fdb39f4ca98b4b9a7402b8bb476c1faf82785a42bb70f5c01483045022100b23414c189a91de419ecc75492c6624f777456e709771e6b09834706c0d8d80d022031d62ad5ee603ee64ea0a59c9b4931c48e3f58fd6923c7aa713092a4c0e1938c014cf15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57aeffffffff73156203b9c4b48c18725ef662d0e10e8eea8f4c0de114c780a5fc2f286dd45c00000000fd6002004830450221008ca51d3e439184bd53ea3d77a1ffd1bc90f1fc5b987260e20b30d41af0c39ad302207fa0644bbc1742d0feacc7c211c3eb282b1b2ebb887b4a7960ad52738a4e9f2301483045022100eca8a530a9e3b3cc430baf798fb9a13ebb8b1faca555e12e8bfe8e78c533e05f02200947165c9bf52191a42fd7dc13b2ad356a51629f56bf4fdeb84b4999f39dfcfb014730440220186ec60d2d1f75d5e96b614cb11ae6b66b8476b2a335d529613434e0453a2c8e02206e0343ee3f7e71d6549eaf0b90e6ce3f1120dc7d126efede4b2a049d90881f6e01483045022100e106510eb72cdd1258cdad9a7aecf26635e03bf5c381c8a367797ee1f4c64ffe02203118a440f5f476cecb9a96db95adb67be4532fd744d4f1510729066f00fb357301483045022100ac49e9317de20ea6cebbf53e9be6a7dc7f6248adb3236ccfd8a046c1ae6b47b9022075bd427071a72418ee8d93f06d4c28a6aaf21fedebba30cdbc16e2460253403e014cf15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57aeffffffff02fd8f0000000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac012a00000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c274300000000010000000001046ca7c2a1a180eab6ded6a51397dbcba9d5c189ad3811b4d6944c2128f38bdfcb0000000000ffffffffdccf67d24791e8a5a039bcbe94bd4d0d6f00c29226f488c0724289412d6ed0460000000000ffffffff8e855370d74eb327f2b6fa8e20e91b3a53cb0b439afe7baad8287e5c8ef5221900000000fd5f0200483045022100c46f3b11c080c6930b8dcdbcafaecf3e0dd6bada902accccb18c599cc887d6cb022027e85575b7d293a298eacf69a0f73099f570edd842ad0c79facc900335d75e2201483045022100f7f988ee55bd9825415ab164835c62d6fae9e2c60009723b156fe5717b9aac8b0220141132cfd4f814eba469f41a1ad9d70cdccbcb45a141e1c9bb2ae33729e90af00148304502210089221802621b8d17effa83edab855784801ab12cca49a75668cd34bc14e1e6b802205d581aa1186b164cdfb189552f0b9e3fb4600479a2d541599e75b26452222a5f0147304402201bff68315eb5f0b813c582e5517480dd8b3469ae5dbbb931e82a7b9b198066f5022029be19b843e657ac1f259abe44cc47cd6cbeb0d5d447bbcdea5105c1f6b6d16a014730440220680831bafdc86cdcec708af675a3b0562505236abc66f3fcba8b1403ef1c1588022057d8f061dcf998d8a7c9139d54faebe4bfbcca43067c0a7556cb289c4ca3958c014cf15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57aeffffffff4797b1f930b026ade7e1654789415ee834bf7591941c1a0929350fbc2c4478fe00000000fd5f020047304402206e7847775ca340a26faf11a661ba42fe63908110c56be4c239f2be766e5ef10b02200298a390b9bd3a8ab5becfd72b0d141414815cee48b24462fc4b98c6c71771f701483045022100f1314bcd38f42015ba898f4a65d1b3726afe557babfc030e835693095c24bdd0022068ca841d6e2932db8c64439a758c0863921793a5c3dbbcdf74467e45d81a79900147304402206d0238af52fefd14b0b8391c474bc0d753567e254f1155241bcb14e253e5e8040220438996c901424750abe9128cb67aef1f40ad0b10010cced0b92e1c1d56e4d3a40148304502210097986fc25b4010056bdab2d394bb3b0ed913d22a61cb117854f534618edc13f60220673c8b8d4a9bd392116fcce8048c418a95e986481ccb16a93eea8f7dec24bcb601483045022100b348c3a9e20a330e822485250b3da7053d6ce6380b6625da23e3bf6bcac8280f022073c233938575062d2a6dad72bb31c68fe1a7951c7321d2c56ac1356e1206d78d014cf15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57aeffffffff02d94b0100000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0d0301000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c27430700483045022100ea0f9175284af169f60e1fcb79accaf2dff2d51a4f30b413aa2e2dc8b861b0cb02201efeb693db3de82106adfa1ff3306341144753fe2786e0de01a6b38e7a4115a90147304402205871c0a8f209c55c4e8675d6f5666908086ceb5dcc9ac37b03c1fe74e4dd90ac02205976315ae147eff9761c17dd5d83f103ae092067d86afe32e7298cbf91f4966101483045022100c585a8516dcc087859378bc804ea69aa56e7b9cf04d195578a29759cffc364840220626e78b7c9f2cd929961c6e598f4179fe0b7a8aad2b4a6ef3a5fd0ddccfc6f0e01483045022100d1ff04e6404b6c900fcf53328711e9537e772f7feb1d9ee2bffef6b67babd4f90220378d5f1f782d7628719e5b350304e8faa2692e8982697eb49c9738a0840564c201483045022100c63f577d7d9ce565491083c32db04ecea96389f7bfacf90913d9ec1d4b5773be02200b6c07739ea90feff7238437552f68be0ed2fc65a41b619b2024c5ba8506211101f15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae07004830450221009aac812eeae9cc658355be275ad13f6666706ebe79363841ec5b8f07179c85b50220148f76b36400f19805a967efd408754e3e2dd511b1b943aaa218e564cd61415601483045022100f92a27d532d5a1e1d7b3b1018c93b61ba0f63fce7f3d07792cfcc7ed1e82b04e022074352ad7066ee0efafcbca88701d003dab8044f996139f02b17292833c9a1ab201483045022100942cfceb2de4e949862c66f05a8942712ec5fe7a53890e5d2a939e973975b960022054b9ac547efd7d0c20af1f068db7814fb0e21422649ea0310440c6f329a2827101483045022100c3a67a1b9ea3b898ef6ea35d0be57a7f0161651c2d3990fd295c043b5c40f5f20220702a31dc5936ebba91b7af22e8753344d67864c68b94859055d36381cfa4bf9701483045022100deef82066d023395673f9c79db443410778a9e0bc28a9d617ebfc39db3175283022039421787c9c06604812876aedf786b8b33ee0cca43c290f82738352482cfb59301f15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae00000000000001000000000102f13323104dddac70bb32ee589e043511019a940e50b0def957583457166510410000000000ffffffff70e60ef9406e4b151e01c95e02d83b5ea50114c5f13ea406a8e00c7045be80d50000000000ffffffff027b960000000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ace23100000000000022002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c2743070047304402201d2a6562d5e01ccffc5924db6408f6df95f3ac1821bd243d19bee2458b8b2f1a022041fd6236b1db18b9bdb8414fe5d0fecb8e56d8ff4d4574d5a0ae80d3b4d4a494014730440220449242aa366334a10b530b9e131ccc2fae2f28118aa2b7d9ec1238ec39f729a60220298c4f0250d117ea2d4a0e282127a28b59b83b0100fb4065d46fcae63ed8ae9201483045022100c2c8b0f09895425d865c4810e9d329cdb387e1689d3c32492aa8c4d03ee6a02602206f99807c26f9d5296cdc1280c2ed161f5b454e9655806bd974730d769299f70301483045022100d4215859adda1b0178450bfd484db2111125302f6a0b7f4e9837f36c29712cb002202424f3f713c3c96b0053cb78ce75d23914cd49b54e1c6d430f88d1242c0385a60148304502210096cb52526c43b329c10420663bb1a2c6ef17b880fdb0946047b5a4886301ef2a022072d2721ccb10da1d006c147a48cd73ecf55a1b7d309ff1e91b0ee56a348630a201f15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae07004830450221008c96ba86676a7e653ee8b3964c0db554cc43ec1a758e866896f72d913149a57102201a9b8887b85bcb257c12b9de794a92f017187e3ddc4749b3c509faeeaf29d7b90147304402200d3350f4c3bc9b0b66f6655d80c30bb616b705192cc0b51ada473c41e807ee640220519c898ada9d3d30e1736e0af33764988babd5d6d79c26e6d689ad891641bc1101483045022100f3fa7300cd3b5e7d0d951929a02ae547fead73a19553be049ee8a9b5b174d3660220668eed549031f66eea1bfe316084f0635fbc92789a400abf607c321bd40a7411014730440220589a6fa1d7ae330faae1bdad8a2055258511c2e09b33cfe0709200f1abf200d502206c081bfae13a06e1d2cd6d1caaf6ffbfc3e2794053af81d3b1d7ad762644868e01483045022100b9458c7de776a8c475cc76f5af27dfc5b46f80c420457232083e71de85fc68b50220769c3ac0cc35372a7602fd516edc24e1619dfaba24277c72e45c3d8a41ac9cba01f15521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae00000000"
)

func GetCCI() (CrossChainItem, int) {
	txid, _ := chainhash.NewHashFromStr("fd285cf687d0759f215e63b3234e1e2a010cb8060d2b4775793ec4447b8385c1")
	proof, _ := hex.DecodeString("01000030c64cb1478e2a2a1774a55ea35b272a397730422ec3e47244f9eb0a062d5a3f1c88a8701698569f51645ff81ae2a8f45c1dbcf369aa9e1f8600f02b20d2c8a156dfb30e5effff7f2000000000080000000490eb94f4609ad4d50547ea65730da8d276777a2c90056f9978a5a2df47eabb6c1786f80a919e1531433f0be9469987602193c5abe5ffa9520dc0942958e6c521c185837b44c43e7975472b0d06b80c012a1e4e23b3635e219f75d087f65c28fd51634411a436e65609dce1f0eb21da7109f01ab468927f81270e90684499b9a7012b")
	rawTx, _ := hex.DecodeString("0100000001633f586f140397287ee87b41c1feb4aad5447bbdc66bc4a11f3a400a637ec0e3020000006b483045022100c67f2ad9b598134c9554490bc411349a122529e9b052ee2673d4326021c69cfb02205f8f63ccab64de8244fecc8727afa38869fddda10770f0cff62227cc10e212da012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff039d8f00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660200000000000000000000000000000014dc68bcc275bf668129c6d214202f6d6ee77e309214f3b8a17f1f957f60c88f105e32ebff3f022e56a45a37052a010000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000")
	return CrossChainItem{
		Height: 100,
		Txid:   *txid,
		Proof:  proof,
		Tx:     rawTx,
	}, 36 + 8 + len(proof) + len(rawTx)
}

func TestCrossChainItem_Serialize(t *testing.T) {
	cci, l := GetCCI()
	b, err := cci.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != l {
		t.Fatalf("%d is not right and should be %d", len(b), l)
	}
}

func TestCrossChainItem_Deserialize(t *testing.T) {
	cci, _ := GetCCI()
	b, err := cci.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	var cciNew CrossChainItem
	err = cciNew.Deserialize(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(cci.Tx, cciNew.Tx) || !bytes.Equal(cci.Proof, cciNew.Proof) || !cci.Txid.IsEqual(&cciNew.Txid) ||
		cci.Height != cciNew.Height {
		t.Fatal("val is not equal")
	}
}

func GetCCIArr(n int) (CrossChainItemArr, int) {
	var res CrossChainItemArr
	length := 0
	for i := 1; i <= n; i++ {
		cci, l := GetCCI()
		cci.Height = uint32(i)
		length += l
		res = append(res, &cci)
	}

	return res, length + 4*(n+1)
}

func TestCrossChainItemArr_Serialize(t *testing.T) {
	arr, l := GetCCIArr(10)
	arrb, err := arr.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if len(arrb) != l {
		t.Fatal("wrong length")
	}
}

func TestCrossChainItemArr_Deserialize(t *testing.T) {
	arr, _ := GetCCIArr(10)
	arrb, err := arr.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	var arrNew CrossChainItemArr
	err = arrNew.Deserialize(arrb)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(arr); i++ {
		cci := arr[i]
		cciNew := arrNew[i]
		if !bytes.Equal(cci.Tx, cciNew.Tx) || !bytes.Equal(cci.Proof, cciNew.Proof) || !cci.Txid.IsEqual(&cciNew.Txid) ||
			cci.Height != cciNew.Height {
			t.Fatalf("no %d val is not equal", i)
		}
	}
}

func TestRestCli_GetProof(t *testing.T) {
	cli := NewRestCli(startMockBtcServer(), USER, PWD)
	proof, err := cli.GetProof([]string{""})
	assert.NoError(t, err)
	assert.Equal(t, "proof", proof)
}

func TestRestCli_GetCurrentHeight(t *testing.T) {
	cli := NewRestCli(startMockBtcServer(), USER, PWD)
	h, hash, err := cli.GetCurrentHeightAndHash()
	assert.NoError(t, err)
	assert.Equal(t, uint32(1403), h)
	assert.Equal(t, "4ea571b949996d9380df6b13b62084d921986d9b65ba3872c937bcb091fdcdc9", hash)
}

func TestRestCli_GetTxsInBlock(t *testing.T) {
	cli := NewRestCli(startMockBtcServer(), USER, PWD)
	blk, err := cli.GetTxsInBlock("")
	assert.NoError(t, err)
	assert.Equal(t, blk.BlockHash().String(), "033926b8bbf6457952af9accafe829e2e8a058595966f8565ef912b8d348b0aa")
}

func TestRestCli_GetScriptPubKey(t *testing.T) {
	cli := NewRestCli(startMockBtcServer(), USER, PWD)
	s, err := cli.GetScriptPubKey("", 0)
	assert.NoError(t, err)
	assert.Equal(t, "002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c2743", s)
}

func TestRestCli_BroadcastTx(t *testing.T) {
	rawtx := "01000000015a93813ac8d05a5a36168d3383ebb23d9c833443132ef4fda34242c7c74b966a020000006a4730440220793143bf61db374c268239646a386af25a25cdf24523089bb5c11c5ed177e3a902200aeb8230a8941ccde68c4cc915bbf657f6cbf58b18665736746db1b162e2152e012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03102700000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000dab47e816313a79c9459b544720c90a725264e0d10684a1f000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000"
	cli := NewRestCli(startMockBtcServer(), USER, PWD)
	_, err := cli.BroadcastTx(rawtx)
	assert.NoError(t, err)
}

func startMockBtcServer() string {
	fmt.Println("starting mock server")
	ms := httptest.NewServer(http.HandlerFunc(handleReq))
	return ms.URL
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	rb, _ := ioutil.ReadAll(r.Body)
	req := new(Request)
	_ = json.Unmarshal(rb, req)

	switch req.Method {
	case "gettxoutproof":
		res, _ := btcjson.MarshalResponse(1, "proof", nil)
		w.Write(res)
	case "getblock":
		res, _ := btcjson.MarshalResponse(1, BLOCK, nil)
		w.Write(res)
	case "getblockhash":
		res, _ := btcjson.MarshalResponse(1, "033926b8bbf6457952af9accafe829e2e8a058595966f8565ef912b8d348b0aa", nil)
		w.Write(res)
	case "getblockheader":
		res, _ := btcjson.MarshalResponse(1, "00000020b797de83cf1f8865be7e894ccc8116840bc5280d4d2ddd7c330c00f"+
			"54e786e04330486c1ed6fcb9c163d46f0a986a701cc465fe3ac3a9faae13b84493c9c0797cfe03b5effff7f2000000000", nil)
		w.Write(res)
	case "getchaintips":
		resp := make(map[string]interface{})
		resp["height"] = 1403
		resp["hash"] = "4ea571b949996d9380df6b13b62084d921986d9b65ba3872c937bcb091fdcdc9"
		resp["branchlen"] = 0
		resp["status"] = "active"
		res, _ := btcjson.MarshalResponse(1, []interface{}{resp}, nil)
		w.Write(res)
	case "getrawtransaction":
		resp := make(map[string]interface{})
		outVal := make(map[string]interface{})
		outVal["hex"] = "002044978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c2743"
		out := make(map[string]interface{})
		out["scriptPubKey"] = outVal
		resp["vout"] = []interface{}{out}
		res, _ := btcjson.MarshalResponse(1, resp, nil)
		w.Write(res)
	case "sendrawtransaction":
		res, _ := btcjson.MarshalResponse(1, "", nil)
		w.Write(res)
	default:
		fmt.Fprint(w, "wrong method")
	}
}