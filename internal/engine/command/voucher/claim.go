package voucher

import (
	"errors"
	"fmt"
	"time"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
)

func (v *VoucherCmd) claimHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	code := args["code"]
	if len(code) != 8 {
		return cmd.ErrorResult(errors.New("voucher code is not valid, length must be 8"))
	}

	voucher, err := v.db.GetVoucherByCode(code)
	if err != nil {
		return cmd.ErrorResult(errors.New("voucher code is not valid, no voucher found"))
	}

	if voucher.CreatedAt.AddDate(0, int(voucher.ValidMonths), 0).Before(time.Now()) {
		return cmd.ErrorResult(errors.New("voucher is expired"))
	}

	if voucher.IsClaimed() {
		return cmd.ErrorResult(errors.New("voucher code claimed before"))
	}

	address := args["address"]
	valInfo, _ := v.clientManager.GetValidatorInfo(address)
	if valInfo != nil {
		err = errors.New("this address is already a staked validator")
		log.Warn(fmt.Sprintf("staked validator found. %s", address))

		return cmd.ErrorResult(err)
	}

	pubKey, err := v.clientManager.FindPublicKey(address, false)
	if err != nil {
		log.Warn(fmt.Sprintf("peer not found. %s", address))

		return cmd.ErrorResult(err)
	}

	memo := fmt.Sprintf("voucher %s claimed by Pagu", code)
	txHash, err := v.wallet.BondTransaction(pubKey, address, memo, voucher.Amount)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if txHash == "" {
		return cmd.ErrorResult(errors.New("can't send bond transaction"))
	}

	if err = v.db.ClaimVoucher(voucher.ID, txHash, caller.ID); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResultF("Voucher claimed successfully!\n\n https://pacviewer.com/transaction/%s", txHash)
}
