package voucher

import (
	"strconv"

	"github.com/pagu-project/Pagu/pkg/utils"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (v *Voucher) createHandler(cmd command.Command, _ entity.AppID, _ string, args ...string) command.CommandResult {
	//! Admin only check

	code := utils.RandomString(8, utils.CapitalLetterNumbers)
	for _, err := v.db.GetVoucherByCode(code); err == nil; {
		code = utils.RandomString(8, utils.CapitalLetterNumbers)
	}

	recipient := args[0]
	amount := args[1]
	validMonths := args[2]
	description := args[3]

	expireMonths, err := strconv.Atoi(validMonths)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	intAmount, err := strconv.Atoi(amount)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	err = v.db.AddVoucher(&entity.Voucher{
		Creator:     cmd.User.ID,
		Code:        code,
		Desc:        description,
		Recipient:   recipient,
		ValidMonths: uint8(expireMonths),
		Amount:      uint(intAmount),
	})
	if err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Voucher crated successfully!")
}