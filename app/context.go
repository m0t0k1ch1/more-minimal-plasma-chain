package app

import (
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/labstack/echo"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type Context struct {
	echo.Context
}

func NewContext(c echo.Context) *Context {
	return &Context{c}
}

func (c *Context) GetAddressFromPath() (common.Address, error) {
	return c.getAddressFromPath("address")
}

func (c *Context) GetBlockNumberFromPath() (uint64, error) {
	return c.getUint64FromPath("blkNum")
}

func (c *Context) GetTxPositionFromPath() (types.Position, error) {
	return c.getPositionFromPath("txPos")
}

func (c *Context) GetTxInPositionFromPath() (types.Position, error) {
	return c.getPositionFromPath("txInPos")
}

func (c *Context) getAddressFromPath(key string) (common.Address, error) {
	addrStr := c.getPathParam(key)
	if !utils.IsHexAddress(addrStr) {
		return types.NullAddress, NewInvalidPathParamError(key)
	}

	return utils.HexToAddress(addrStr), nil
}

func (c *Context) getUint64FromPath(key string) (uint64, error) {
	return strconv.ParseUint(c.getPathParam(key), 10, 64)
}

func (c *Context) getPositionFromPath(key string) (types.Position, error) {
	i, err := c.getUint64FromPath(key)
	if err != nil {
		return 0, err
	}

	return types.Position(i), nil
}

func (c *Context) getPathParam(key string) string {
	return c.Param(key)
}

func (c *Context) GetConfirmationSignatureFromForm() (types.Signature, error) {
	return c.getRequiredSignatureFromForm("confsig")
}

func (c *Context) GetTxFromForm() (*types.Tx, error) {
	return c.getRequiredTxFromForm("tx")
}

func (c *Context) GetOwnerAddressFromForm() (common.Address, error) {
	return c.getRequiredAddressFromForm("owner")
}

func (c *Context) GetAmountFromForm() (uint64, error) {
	return c.getRequiredUint64FromForm("amount")
}

func (c *Context) getRequiredSignatureFromForm(key string) (types.Signature, error) {
	sigStr, err := c.getRequiredFormParam(key)
	if err != nil {
		return types.NullSignature, err
	}

	sig, err := types.HexToSignature(sigStr)
	if err != nil {
		return types.NullSignature, NewInvalidFormParamError(key)
	}

	return sig, nil
}

func (c *Context) getRequiredTxFromForm(key string) (*types.Tx, error) {
	txStr, err := c.getRequiredFormParam(key)
	if err != nil {
		return nil, err
	}

	txBytes, err := utils.DecodeHex(txStr)
	if err != nil {
		return nil, NewInvalidFormParamError(key)
	}

	var tx types.Tx
	if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
		return nil, NewInvalidFormParamError(key)
	}

	return &tx, nil
}

func (c *Context) getRequiredAddressFromForm(key string) (common.Address, error) {
	addrStr, err := c.getRequiredFormParam(key)
	if err != nil {
		return types.NullAddress, err
	}

	if !utils.IsHexAddress(addrStr) {
		return types.NullAddress, NewInvalidFormParamError(key)
	}

	return utils.HexToAddress(addrStr), nil
}

func (c *Context) getRequiredUint64FromForm(key string) (uint64, error) {
	s, err := c.getRequiredFormParam(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(s, 64, 10)
}

func (c *Context) getRequiredFormParam(key string) (string, error) {
	if !c.isExistFormParam(key) {
		return "", NewRequiredFormParamError(key)
	}
	return c.getFormParam(key), nil
}

func (c *Context) getFormParam(key string) string {
	return c.FormValue(key)
}

func (c *Context) isExistFormParam(key string) bool {
	_, ok := c.Request().Form[key]
	return ok
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
