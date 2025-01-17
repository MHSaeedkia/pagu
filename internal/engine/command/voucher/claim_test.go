package voucher

import (
	"testing"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClaim(t *testing.T) {
	td := setup(t)

	voucherCode := "12345678"
	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}
	cmd := &command.Command{}

	t.Run("Invalid Voucher Code", func(t *testing.T) {
		args := map[string]string{
			"code":    "0",
			"address": "pc1p...",
		}
		result := td.voucherCmd.claimHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "An error occurred: voucher code is not valid, length must be 8")
	})

	t.Run("Voucher Code Not Issued Yet", func(t *testing.T) {
		args := map[string]string{
			"code":    voucherCode,
			"address": "pc1p...",
		}
		result := td.voucherCmd.claimHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "An error occurred: voucher code is not valid, no voucher found")
	})

	t.Run("Claim a Voucher", func(t *testing.T) {
		testVoucher := td.createTestVoucher(t, WithCode(voucherCode))
		validatorAddr := "pc1p..."

		td.mockClientMgr.EXPECT().GetValidatorInfo(validatorAddr).Return(
			nil, nil,
		).AnyTimes()

		td.mockClientMgr.EXPECT().FindPublicKey(validatorAddr, false).Return(
			validatorAddr, nil,
		).AnyTimes()

		td.mockWallet.EXPECT().BondTransaction(gomock.Any(), validatorAddr,
			"voucher 12345678 claimed by Pagu", testVoucher.Amount).Return(
			"0x1", nil,
		).AnyTimes()

		args := map[string]string{
			"code":    voucherCode,
			"address": validatorAddr,
		}
		result := td.voucherCmd.claimHandler(caller, cmd, args)
		assert.True(t, result.Successful)
		assert.Equal(t, result.Message, "Voucher claimed successfully!\n\n https://pacviewer.com/transaction/0x1")
	})

	t.Run("Claim again", func(t *testing.T) {
		args := map[string]string{
			"code":    voucherCode,
			"address": "pc1p...",
		}
		result := td.voucherCmd.claimHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "An error occurred: voucher code claimed before")
	})
}
