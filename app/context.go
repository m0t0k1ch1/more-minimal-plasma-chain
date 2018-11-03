package app

import (
	"math/big"
	"net/http"

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

func (c *Context) GetBlockNumberFromPath() (*big.Int, error) {
	return c.getBigIntFromPath("blkNum")
}

func (c *Context) GetBlockHashFromPath() (common.Hash, error) {
	return c.getHashFromPath("blkHash")
}

func (c *Context) GetTxHashFromPath() (common.Hash, error) {
	return c.getHashFromPath("txHash")
}

func (c *Context) getBigIntFromPath(key string) (*big.Int, error) {
	i, ok := new(big.Int).SetString(c.getPathParam(key), 10)
	if !ok {
		return big.NewInt(0), NewInvalidPathParamError(key)
	}

	return i, nil
}

func (c *Context) getHashFromPath(key string) (common.Hash, error) {
	b, err := utils.DecodeHex(c.getPathParam(key))
	if err != nil {
		return types.NullHash, NewInvalidPathParamError(key)
	}
	if len(b) != common.HashLength {
		return types.NullHash, NewInvalidFormParamError(key)
	}

	return common.BytesToHash(b), nil
}

func (c *Context) getPathParam(key string) string {
	return c.Param(key)
}

func (c *Context) GetInputIndexFromForm() (*big.Int, error) {
	return c.getRequiredBigIntFromForm("index")
}

func (c *Context) GetConfirmationSignatureFromForm() (types.Signature, error) {
	return c.getRequiredSignatureFromForm("confsig")
}

func (c *Context) GetTxFromForm() (*types.Tx, error) {
	return c.getRequiredTxFromForm("tx")
}

func (c *Context) getRequiredBigIntFromForm(key string) (*big.Int, error) {
	iStr, err := c.getRequiredFormParam(key)
	if err != nil {
		return big.NewInt(0), err
	}

	i, ok := new(big.Int).SetString(iStr, 10)
	if !ok {
		return big.NewInt(0), NewInvalidFormParamError(key)
	}

	return i, nil
}

func (c *Context) getRequiredSignatureFromForm(key string) (types.Signature, error) {
	sigStr, err := c.getRequiredFormParam(key)
	if err != nil {
		return types.NullSignature, err
	}

	sig, err := types.NewSignatureFromHex(sigStr)
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
