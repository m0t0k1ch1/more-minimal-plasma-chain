package app

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/contract"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

type HandlerFunc func(*Context) error

type ChildChain struct {
	e          *echo.Echo
	config     *Config
	rootChain  *contract.RootChain
	operator   *types.Account
	blockchain *core.Blockchain
}

func NewChildChain(conf *Config) (*ChildChain, error) {
	rc, err := newRootChain(conf.RootChain)
	if err != nil {
		return nil, err
	}

	op, err := newOperator(conf.Operator)
	if err != nil {
		return nil, err
	}

	cc := &ChildChain{
		e:          echo.New(),
		config:     conf,
		rootChain:  rc,
		operator:   op,
		blockchain: core.NewBlockchain(),
	}

	cc.e.Use(middleware.Logger())
	cc.e.Use(middleware.Recover())
	cc.e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{c})
		}
	})
	cc.e.HTTPErrorHandler = cc.httpErrorHandler

	cc.GET("/ping", cc.PingHandler)

	cc.GET("/chain/:blkNum", cc.GetChainHandler)

	cc.POST("/blocks", cc.PostBlockHandler)
	cc.GET("/blocks/:blkHash", cc.GetBlockHandler)

	cc.POST("/txes", cc.PostTxHandler)
	cc.GET("/txes/:txHash", cc.GetTxHandler)
	cc.GET("/txes/:txHash/proof", cc.GetTxProofHandler)
	cc.PUT("/txes/:txHash", cc.PutTxHandler)

	return cc, nil
}

func (cc *ChildChain) GET(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.Add(http.MethodGet, path, h, m...)
}

func (cc *ChildChain) POST(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.Add(http.MethodPost, path, h, m...)
}

func (cc *ChildChain) PUT(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.Add(http.MethodPut, path, h, m...)
}

func (cc *ChildChain) Add(method, path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.e.Add(method, path, func(c echo.Context) error {
		return h(NewContext(c))
	})
}

func (cc *ChildChain) Logger() echo.Logger {
	return cc.e.Logger
}

func (cc *ChildChain) Start() error {
	return cc.e.Start(fmt.Sprintf(":%d", cc.config.Port))
}

func (cc *ChildChain) httpErrorHandler(err error, c echo.Context) {
	cc.e.Logger.Error(err)

	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
		msg = fmt.Sprintf("%v", httpErr.Message)
	}

	appErr := NewError(code, msg)

	if err := c.JSON(appErr.Code, NewErrorResponse(appErr)); err != nil {
		cc.e.Logger.Error(err)
	}
}

func newOperator(conf *OperatorConfig) (*types.Account, error) {
	privKey, err := crypto.HexToECDSA(conf.PrivateKey)
	if err != nil {
		return nil, err
	}

	return types.NewAccount(privKey), nil
}

func newRootChain(conf *RootChainConfig) (*contract.RootChain, error) {
	if ok := common.IsHexAddress(conf.Address); !ok {
		return nil, fmt.Errorf("invalid root chain address")
	}
	rcAddr := common.HexToAddress(conf.Address)

	conn, err := ethclient.Dial(conf.RPC)
	if err != nil {
		return nil, err
	}

	return contract.NewRootChain(rcAddr, conn)
}
