package app

import (
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/labstack/echo"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

type Context struct {
	echo.Context
}

func NewContext(c echo.Context) *Context {
	return &Context{c}
}

func (c *Context) GetBlockNumberFromPath() (uint64, error) {
	blkNum, err := strconv.ParseUint(c.Param("blkNum"), 10, 64)
	if err != nil {
		return 0, ErrInvalidBlockNumber
	}

	return blkNum, nil
}

func (c *Context) GetTxIndexFromPath() (int, error) {
	txIndex, err := strconv.Atoi(c.Param("txIndex"))
	if err != nil {
		return 0, ErrInvalidTxIndex
	}

	return txIndex, nil
}

func (c *Context) GetOwnerFromForm() (common.Address, error) {
	if _, ok := c.Request().Form["owner"]; !ok {
		return common.Address{}, ErrOwnerRequired
	}

	ownerStr := c.FormValue("owner")
	if !common.IsHexAddress(ownerStr) {
		return common.Address{}, ErrInvalidOwnerAddressHex
	}

	return common.HexToAddress(ownerStr), nil
}

func (c *Context) GetAmountFromForm() (uint64, error) {
	if _, ok := c.Request().Form["amount"]; !ok {
		return 0, ErrAmountRequired
	}

	amountStr := c.FormValue("amount")
	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidAmount
	}

	return amount, nil
}

func (c *Context) GetTxFromForm() (*types.Tx, error) {
	if _, ok := c.Request().Form["tx"]; !ok {
		return nil, ErrTxRequired
	}

	txCoreStr := c.FormValue("tx")
	txCoreBytes, err := hexutil.Decode(txCoreStr)
	if err != nil {
		return nil, ErrInvalidTxHex
	}

	var txc types.TxCore
	if err := rlp.DecodeBytes(txCoreBytes, &txc); err != nil {
		return nil, ErrInvalidTxHex
	}

	tx := types.NewTx()
	tx.TxCore = &txc

	return tx, nil
}

func (c *Context) JSONSuccess(result interface{}) error {
	return c.JSON(http.StatusOK, NewSuccessResponse(result))
}

func (c *Context) JSONError(err error) error {
	var appErr *Error
	if e, ok := err.(*Error); ok {
		appErr = e
	} else {
		c.Logger().Error(err)
		appErr = ErrUnexpected
	}

	return c.JSON(http.StatusOK, NewErrorResponse(appErr))
}
